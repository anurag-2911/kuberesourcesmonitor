package controller

import (
    "context"
    "time"
    "github.com/go-logr/logr"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/push"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/runtime"
    corev1 "k8s.io/api/core/v1"
    "sigs.k8s.io/controller-runtime/pkg/client"
    
    "sigs.k8s.io/controller-runtime/pkg/reconcile"
    ctrl "sigs.k8s.io/controller-runtime"
    monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
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

    // Fetch the KubeResourcesMonitor instance
    instance := &monitorv1alpha1.KubeResourcesMonitor{}
    err := r.Get(ctx, req.NamespacedName, instance)
    if err != nil {
        if errors.IsNotFound(err) {
            return reconcile.Result{}, nil
        }
        return reconcile.Result{}, err
    }

    // Define a new Pod object
    podList := &corev1.PodList{}
    err = r.Client.List(ctx, podList)
    if err != nil {
        return reconcile.Result{}, err
    }

    // Get the counts of resources
    podCount := len(podList.Items)
    
    serviceList := &corev1.ServiceList{}
    err = r.Client.List(ctx, serviceList)
    if err != nil {
        return reconcile.Result{}, err
    }
    serviceCount := len(serviceList.Items)
    
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

    // Push metrics to Prometheus
    pusher := push.New("http://prometheus:9091", "kuberesourcesmonitor").
        Collector(podGauge).
        Collector(serviceGauge).
        Collector(configMapGauge).
        Collector(secretGauge)

    if err := pusher.Push(); err != nil {
        log.Error(err, "Could not push to Prometheus")
    }

    return reconcile.Result{RequeueAfter: time.Minute * 5}, nil
}

func (r *KubeResourcesMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&monitorv1alpha1.KubeResourcesMonitor{}).
        Complete(r)
}

