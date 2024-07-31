package controller

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
// +kubebuilder:rbac:groups="",resources=pods;services;configmaps;secrets;cronjobs;jobs;deployments;nodes;persistentvolumeclaims;events,verbs=get;list

var (
	// Define Prometheus gauges for monitoring various Kubernetes resources
	podGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pod_count",
		Help: "Number of pods",
	})
	serviceGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "service_count",
		Help: "Number of services",
	})
	configMapGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "configmap_count",
		Help: "Number of configmaps",
	})
	secretGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "secret_count",
		Help: "Number of secrets",
	})
	cronJobGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cronjob_count",
		Help: "Number of cronjobs",
	})
	deploymentGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "deployment_count",
		Help: "Number of deployments",
	})
	nodeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_count",
		Help: "Number of nodes",
	})
	nodeReadyGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_ready_count",
		Help: "Number of nodes in ready condition",
	})
	nodeMemoryPressureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_memory_pressure_count",
		Help: "Number of nodes with memory pressure",
	})
	nodeDiskPressureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_disk_pressure_count",
		Help: "Number of nodes with disk pressure",
	})
	pvcGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pvc_count",
		Help: "Number of PVCs",
	})
	eventGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "event_count",
		Help: "Number of events",
	})
	restartGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "restart_count",
		Help: "Number of container restarts",
	})
	crashGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "crash_count",
		Help: "Number of container crashes",
	})
	cpuUsageGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_cpu_usage",
		Help: "Total CPU usage across all nodes",
	})
	memoryUsageGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_memory_usage",
		Help: "Total memory usage across all nodes",
	})
)

func init() {
	// Register the gauges with Prometheus's default registry
	prometheus.MustRegister(podGauge, serviceGauge, configMapGauge, secretGauge, cronJobGauge, deploymentGauge, nodeGauge,
		nodeReadyGauge, nodeMemoryPressureGauge, nodeDiskPressureGauge, pvcGauge, eventGauge, restartGauge, crashGauge,
		cpuUsageGauge, memoryUsageGauge)
}

// Reconcile reads the state of the cluster for a KubeResourcesMonitor object and makes changes based on the state read
// and what is in the KubeResourcesMonitor.Spec
func (r *KubeResourcesMonitorReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("kuberesourcesmonitor", req.NamespacedName)
	log.Info("Reconcile function called")
	log.Info("Reconciling KubeResourcesMonitor", "namespace", req.Namespace, "name", req.Name)

	// Fetch the KubeResourcesMonitor instance
	instance := &monitorv1alpha1.KubeResourcesMonitor{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			log.Info("KubeResourcesMonitor resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get KubeResourcesMonitor")
		return reconcile.Result{}, err
	}

	// Collect Kubernetes resources metrics
	shouldReturn, result, err := r.collectKubResourcesMetrics(ctx, req, instance, log)
	if shouldReturn {
		return result, err
	}

	// Check if deployments are provided
	if len(instance.Spec.Deployments) == 0 {
		log.Info("No deployments specified for scaling, skipping autoscaling logic")
	} else {
		// Autoscaling logic based on time
		timeBasedAutoScale(ctx, instance, r, req, log)
	}

	log.Info("Successfully reconciled KubeResourcesMonitor")

	// Set the requeue interval based on the interval specified in the CR spec
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

// SetupWithManager sets up the controller with the Manager.
func (r *KubeResourcesMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Log.Info("SetupWithManager called")

	// Start a metrics server for Prometheus
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		r.Log.Info("Starting metrics server on port 2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			r.Log.Error(err, "Error starting HTTP server")
		}
	}()

	return ctrl.NewControllerManagedBy(mgr).
		For(&monitorv1alpha1.KubeResourcesMonitor{}).
		Complete(r)
}

// timeBasedAutoScale scales deployments based on specified time windows
func timeBasedAutoScale(ctx context.Context, instance *monitorv1alpha1.KubeResourcesMonitor, r *KubeResourcesMonitorReconciler, req reconcile.Request, log logr.Logger) {
	currentTime := time.Now().Format("15:04")
	for _, deployment := range instance.Spec.Deployments {
		for _, scaleTime := range deployment.ScaleTimes {
			if currentTime >= scaleTime.StartTime && currentTime <= scaleTime.EndTime {
				log.Info("Scaling deployment", "deployment", deployment.Name, "replicas", scaleTime.Replicas)
				r.scaleDeployment(ctx, deployment.Name, scaleTime.Replicas, deployment.Namespace, log)
			}
		}
	}
}

// collectKubResourcesMetrics collects various Kubernetes resource metrics and updates Prometheus gauges
func (r *KubeResourcesMonitorReconciler) collectKubResourcesMetrics(ctx context.Context, req reconcile.Request, instance *monitorv1alpha1.KubeResourcesMonitor, log logr.Logger) (bool, reconcile.Result, error) {
	// Fetch the ConfigMap
	configMap := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: instance.Spec.ConfigMapName}, configMap)
	if err != nil {
		log.Error(err, "Failed to get ConfigMap", "name", instance.Spec.ConfigMapName)
		return true, reconcile.Result{}, err
	}
	metricsToCollect := strings.Split(configMap.Data["metrics"], "\n")
	log.Info("Metrics to collect from ConfigMap", "metrics", metricsToCollect)

	// Collect various metrics from the cluster
	podList := &corev1.PodList{}
	err = r.Client.List(ctx, podList)
	if err != nil {
		log.Error(err, "Failed to list Pods")
		return true, reconcile.Result{}, err
	}
	podCount := len(podList.Items)
	log.Info("Collected pod metrics", "count", podCount)

	serviceList := &corev1.ServiceList{}
	err = r.Client.List(ctx, serviceList)
	if err != nil {
		log.Error(err, "Failed to list Services")
		return true, reconcile.Result{}, err
	}
	serviceCount := len(serviceList.Items)
	log.Info("Collected service metrics", "count", serviceCount)

	configMapList := &corev1.ConfigMapList{}
	err = r.Client.List(ctx, configMapList)
	if err != nil {
		log.Error(err, "Failed to list ConfigMaps")
		return true, reconcile.Result{}, err
	}
	configMapCount := len(configMapList.Items)
	log.Info("Collected config map metrics", "count", configMapCount)

	secretList := &corev1.SecretList{}
	err = r.Client.List(ctx, secretList)
	if err != nil {
		log.Error(err, "Failed to list Secrets")
		return true, reconcile.Result{}, err
	}
	secretCount := len(secretList.Items)
	log.Info("Collected secret metrics", "count", secretCount)

	cronJobList := &batchv1.CronJobList{}
	err = r.Client.List(ctx, cronJobList)
	if err != nil {
		log.Error(err, "Failed to list CronJobs")
		return true, reconcile.Result{}, err
	}
	cronJobCount := len(cronJobList.Items)
	log.Info("Collected cronjob metrics", "count", cronJobCount)

	deploymentList := &appsv1.DeploymentList{}
	err = r.Client.List(ctx, deploymentList)
	if err != nil {
		log.Error(err, "Failed to list Deployments")
		return true, reconcile.Result{}, err
	}
	deploymentCount := len(deploymentList.Items)
	log.Info("Collected deployment metrics", "count", deploymentCount)

	nodeList := &corev1.NodeList{}
	err = r.Client.List(ctx, nodeList)
	if err != nil {
		log.Error(err, "Failed to list Nodes")
		return true, reconcile.Result{}, err
	}
	nodeCount := len(nodeList.Items)
	log.Info("Collected node metrics", "count", nodeCount)

	pvcList := &corev1.PersistentVolumeClaimList{}
	err = r.Client.List(ctx, pvcList)
	if err != nil {
		log.Error(err, "Failed to list PVCs")
		return true, reconcile.Result{}, err
	}
	pvcCount := len(pvcList.Items)
	log.Info("Collected PVC metrics", "count", pvcCount)

	eventList := &corev1.EventList{}
	err = r.Client.List(ctx, eventList)
	if err != nil {
		log.Error(err, "Failed to list Events")
		return true, reconcile.Result{}, err
	}
	eventCount := len(eventList.Items)
	log.Info("Collected event metrics", "count", eventCount)

	// Calculate container restart and crash metrics
	var restartCount int
	var crashCount int
	for _, pod := range podList.Items {
		for _, status := range pod.Status.ContainerStatuses {
			restartCount += int(status.RestartCount)
			if status.State.Terminated != nil && status.State.Terminated.ExitCode != 0 {
				crashCount++
			}
		}
	}
	log.Info("Collected container metrics", "restartCount", restartCount, "crashCount", crashCount)

	// Calculate node metrics
	var totalCPUUsageMilli, totalMemoryUsageBytes float64
	var nodeReady, nodeMemoryPressure, nodeDiskPressure float64
	for _, node := range nodeList.Items {
		for _, condition := range node.Status.Conditions {
			switch condition.Type {
			case corev1.NodeReady:
				if condition.Status == corev1.ConditionTrue {
					nodeReady++
				}
			case corev1.NodeMemoryPressure:
				if condition.Status == corev1.ConditionTrue {
					nodeMemoryPressure++
				}
			case corev1.NodeDiskPressure:
				if condition.Status == corev1.ConditionTrue {
					nodeDiskPressure++
				}
			}
		}
		for resourceName, allocatable := range node.Status.Allocatable {
			if resourceName == corev1.ResourceCPU {
				totalCPUUsageMilli += float64(allocatable.MilliValue())
			} else if resourceName == corev1.ResourceMemory {
				totalMemoryUsageBytes += float64(allocatable.Value())
			}
		}
	}
	log.Info("Collected node resource metrics", "totalCPUUsageMilli", totalCPUUsageMilli, "totalMemoryUsageBytes", totalMemoryUsageBytes)

	// Convert CPU usage from milli to cores and memory usage to GB
	totalCPUUsage := totalCPUUsageMilli / 1000.0
	totalMemoryUsage := totalMemoryUsageBytes / (1024.0 * 1024.0 * 1024.0)

	// Log the collected metrics
	log.Info("Resource counts", "Pods", podCount, "Services", serviceCount, "ConfigMaps", configMapCount, "Secrets", secretCount, "CronJobs", cronJobCount, "Deployments", deploymentCount, "Nodes", nodeCount, "PVCs", pvcCount, "Events", eventCount, "Restarts", restartCount, "Crashes", crashCount, "TotalCPUUsage (cores)", totalCPUUsage, "TotalMemoryUsage (GB)", totalMemoryUsage, "NodeReady", nodeReady, "NodeMemoryPressure", nodeMemoryPressure, "NodeDiskPressure", nodeDiskPressure)

	// Update Prometheus gauge values with the collected metrics
	log.Info("Setting Prometheus gauge values")
	podGauge.Set(float64(podCount))
	serviceGauge.Set(float64(serviceCount))
	configMapGauge.Set(float64(configMapCount))
	secretGauge.Set(float64(secretCount))
	cronJobGauge.Set(float64(cronJobCount))
	deploymentGauge.Set(float64(deploymentCount))
	nodeGauge.Set(float64(nodeCount))
	nodeReadyGauge.Set(nodeReady)
	nodeMemoryPressureGauge.Set(nodeMemoryPressure)
	nodeDiskPressureGauge.Set(nodeDiskPressure)
	pvcGauge.Set(float64(pvcCount))
	eventGauge.Set(float64(eventCount))
	restartGauge.Set(float64(restartCount))
	crashGauge.Set(float64(crashCount))
	cpuUsageGauge.Set(totalCPUUsage)
	memoryUsageGauge.Set(totalMemoryUsage)

	return false, reconcile.Result{}, nil
}

// scaleDeployment scales the specified deployment to the given number of replicas
// scaleDeployment scales the specified deployment to the given number of replicas
func (r *KubeResourcesMonitorReconciler) scaleDeployment(ctx context.Context, deploymentName string, replicas int32, namespace string, log logr.Logger) {
	deployment := &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: namespace}, deployment)
	if err != nil {
		log.Error(err, "Failed to get Deployment for scaling", "name", deploymentName, "namespace", namespace)
		return
	}

	currentReplicas := *deployment.Spec.Replicas
	if currentReplicas != replicas {
		deployment.Spec.Replicas = &replicas
		err = r.Client.Update(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment replicas", "name", deploymentName, "namespace", namespace)
		} else {
			log.Info("Scaled deployment", "name", deploymentName, "namespace", namespace, "new replicas", replicas)
		}
	}
}


