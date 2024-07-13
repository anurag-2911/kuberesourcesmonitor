package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// KubeResourcesMonitorReconciler reconciles a KubeResourcesMonitor object
type KubeResourcesMonitorReconciler struct {
    client.Client
    Scheme *runtime.Scheme
    Log    logr.Logger
}

// +kubebuilder:rbac:groups=monitor.example.com,resources=kuberesourcesmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitor.example.com,resources=kuberesourcesmonitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitor.example.com,resources=kuberesourcesmonitors/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods;services;configmaps;secrets,verbs=get;list

func (r *KubeResourcesMonitorReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
    log := r.Log.WithValues("kuberesourcesmonitor", req.NamespacedName)
    log.Info("Reconcile function called")
    log.Info("Reconciling KubeResourcesMonitor", "namespace", req.Namespace, "name", req.Name)
    
    // Fetch the KubeResourcesMonitor instance
    instance := &monitorv1alpha1.KubeResourcesMonitor{}
    err := r.Get(ctx, req.NamespacedName, instance)
    if err != nil {
        if errors.IsNotFound(err) {
            log.Info("KubeResourcesMonitor resource not found. Ignoring since object must be deleted")
            return reconcile.Result{}, nil
        }
        log.Error(err, "Failed to get KubeResourcesMonitor")
        return reconcile.Result{}, err
    }
    // Get the Prometheus endpoint from the spec
    prometheusEndpoint := instance.Spec.PrometheusEndpoint
    log.Info("Getting podList")
    // Define a new Pod object
    podList := &corev1.PodList{}
    err = r.Client.List(ctx, podList)
    if err != nil {
        return reconcile.Result{}, err
    }
    log.Info("getting podCount")
    // Get the counts of resources
    podCount := len(podList.Items)
    log.Info("getting services")
    serviceList := &corev1.ServiceList{}
    err = r.Client.List(ctx, serviceList)
    if err != nil {
        return reconcile.Result{}, err
    }
    serviceCount := len(serviceList.Items)
    log.Info("getting configMapList")
    configMapList := &corev1.ConfigMapList{}
    err = r.Client.List(ctx, configMapList)
    if err != nil {
        return reconcile.Result{}, err
    }
    configMapCount := len(configMapList.Items)
    
    secretList := &corev1.SecretList{}
    err = r.Client.List(ctx, secretList)
    if err != nil {
        return reconcile.Result{}, err
    }
    secretCount := len(secretList.Items)

    // Log the counts
    log.Info("Resource counts", "Pods", podCount, "Services", serviceCount, "ConfigMaps", configMapCount, "Secrets", secretCount)

    // Define the gauges
    podGauge := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "pod_count",
        Help: "Number of pods",
    })
    serviceGauge := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "service_count",
        Help: "Number of services",
    })
    configMapGauge := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "configmap_count",
        Help: "Number of configmaps",
    })
    secretGauge := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "secret_count",
        Help: "Number of secrets",
    })

    // Set the gauge values
    podGauge.Set(float64(podCount))
    serviceGauge.Set(float64(serviceCount))
    configMapGauge.Set(float64(configMapCount))
    secretGauge.Set(float64(secretCount))

    // Push metrics to Prometheus Pushgateway
    pusher := push.New(prometheusEndpoint, "kuberesourcesmonitor").
        Collector(podGauge).
        Collector(serviceGauge).
        Collector(configMapGauge).
        Collector(secretGauge)

    if err := pusher.Push(); err != nil {
        log.Error(err, "Could not push to Prometheus Pushgateway")
    }else {
        log.Info("Successfully pushed metrics to Prometheus Pushgateway")
    }
    log.Info("Successfully reconciled KubeResourcesMonitor")
    // Parse interval from the instance spec
    timeInterval := instance.Spec.Interval
    interval := 5 * time.Minute // default interval

    if timeInterval != "" {
        parsedInterval, err := time.ParseDuration(timeInterval)
        if err != nil {
            log.Error(err, "Failed to parse interval, using default", "interval", timeInterval)
        } else {
            interval = parsedInterval
        }
    }

    return reconcile.Result{RequeueAfter: interval}, nil
}

func (r *KubeResourcesMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
    r.Log.Info("SetupWithManager called")
    // Expose the metrics endpoint
    r.Log.Info("Exposing the metrics endpoint for Prometheus to scrape ")
    http.Handle("/metrics", promhttp.Handler())
    go func() {
        err := http.ListenAndServe(":2112", nil)
        if err != nil {
            r.Log.Error(err, "Error starting HTTP server for metrics")
        }
    }()
    return ctrl.NewControllerManagedBy(mgr).
        For(&monitorv1alpha1.KubeResourcesMonitor{}).
        Complete(r)
}
