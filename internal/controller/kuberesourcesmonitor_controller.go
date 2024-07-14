package controller

import (
	"context"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
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
// +kubebuilder:rbac:groups="",resources=pods;services;configmaps;secrets;cronjobs;jobs;deployments;nodes;persistentvolumeclaims;events;daemonsets;replicasets;replicationcontrollers;statefulsets,verbs=get;list

func (r *KubeResourcesMonitorReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("kuberesourcesmonitor", req.NamespacedName)
	log.Info("Reconcile function called")
	log.Info("Reconciling KubeResourcesMonitor", "namespace", req.Namespace, "name", req.Name)

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

	prometheusEndpoint := instance.Spec.PrometheusEndpoint

	configMap := &corev1.ConfigMap{}
	err = r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: "kuberesourcesmonitor-config"}, configMap)
	if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return reconcile.Result{}, err
	}

	metricsToCollect := strings.Split(configMap.Data["metrics"], "\n")

	var (
		podCount, serviceCount, configMapCount, secretCount            int
		cronJobCount, jobCount, daemonSetCount, replicaSetCount        int
		replicationControllerCount, statefulSetCount, deploymentCount  int
		nodeCount, pvcCount, eventCount, restartCount, crashCount      int
		totalCPUUsage, totalMemoryUsage, nodeReady, nodeMemoryPressure float64
		nodeDiskPressure                                               float64
	)

	podList := &corev1.PodList{}
	err = r.Client.List(ctx, podList)
	if err != nil {
		return reconcile.Result{}, err
	}

	for _, metric := range metricsToCollect {
		switch metric {
		case "pods":
			podCount = len(podList.Items)
		case "services":
			serviceList := &corev1.ServiceList{}
			err = r.Client.List(ctx, serviceList)
			if err != nil {
				return reconcile.Result{}, err
			}
			serviceCount = len(serviceList.Items)
		case "configmaps":
			configMapList := &corev1.ConfigMapList{}
			err = r.Client.List(ctx, configMapList)
			if err != nil {
				return reconcile.Result{}, err
			}
			configMapCount = len(configMapList.Items)
		case "secrets":
			secretList := &corev1.SecretList{}
			err = r.Client.List(ctx, secretList)
			if err != nil {
				return reconcile.Result{}, err
			}
			secretCount = len(secretList.Items)
		case "cronjobs":
			cronJobList := &batchv1.CronJobList{}
			err = r.Client.List(ctx, cronJobList)
			if err != nil {
				return reconcile.Result{}, err
			}
			cronJobCount = len(cronJobList.Items)
		case "jobs":
			jobList := &batchv1.JobList{}
			err = r.Client.List(ctx, jobList)
			if err != nil {
				return reconcile.Result{}, err
			}
			jobCount = len(jobList.Items)
		case "daemonsets":
			daemonSetList := &appsv1.DaemonSetList{}
			err = r.Client.List(ctx, daemonSetList)
			if err != nil {
				return reconcile.Result{}, err
			}
			daemonSetCount = len(daemonSetList.Items)
		case "replicasets":
			replicaSetList := &appsv1.ReplicaSetList{}
			err = r.Client.List(ctx, replicaSetList)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicaSetCount = len(replicaSetList.Items)
		case "replicationcontrollers":
			replicationControllerList := &corev1.ReplicationControllerList{}
			err = r.Client.List(ctx, replicationControllerList)
			if err != nil {
				return reconcile.Result{}, err
			}
			replicationControllerCount = len(replicationControllerList.Items)
		case "statefulsets":
			statefulSetList := &appsv1.StatefulSetList{}
			err = r.Client.List(ctx, statefulSetList)
			if err != nil {
				return reconcile.Result{}, err
			}
			statefulSetCount = len(statefulSetList.Items)
		case "deployments":
			deploymentList := &appsv1.DeploymentList{}
			err = r.Client.List(ctx, deploymentList)
			if err != nil {
				return reconcile.Result{}, err
			}
			deploymentCount = len(deploymentList.Items)
		case "nodes":
			nodeList := &corev1.NodeList{}
			err = r.Client.List(ctx, nodeList)
			if err != nil {
				return reconcile.Result{}, err
			}
			nodeCount = len(nodeList.Items)
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
						totalCPUUsage += float64(allocatable.MilliValue()) / 1000
					} else if resourceName == corev1.ResourceMemory {
						totalMemoryUsage += float64(allocatable.Value()) / (1024 * 1024 * 1024)
					}
				}
			}
		case "pvcs":
			pvcList := &corev1.PersistentVolumeClaimList{}
			err = r.Client.List(ctx, pvcList)
			if err != nil {
				return reconcile.Result{}, err
			}
			pvcCount = len(pvcList.Items)
		case "events":
			eventList := &corev1.EventList{}
			err = r.Client.List(ctx, eventList)
			if err != nil {
				return reconcile.Result{}, err
			}
			eventCount = len(eventList.Items)
		case "restarts":
			for _, pod := range podList.Items {
				for _, status := range pod.Status.ContainerStatuses {
					restartCount += int(status.RestartCount)
				}
			}
		case "crashes":
			for _, pod := range podList.Items {
				for _, status := range pod.Status.ContainerStatuses {
					if status.State.Terminated != nil && status.State.Terminated.ExitCode != 0 {
						crashCount++
					}
				}
			}
		}
	}

	log.Info("Resource counts", "Pods", podCount, "Services", serviceCount, "ConfigMaps", configMapCount, "Secrets", secretCount, "CronJobs", cronJobCount, "Jobs", jobCount, "DaemonSets", daemonSetCount, "ReplicaSets", replicaSetCount, "ReplicationControllers", replicationControllerCount, "StatefulSets", statefulSetCount, "Deployments", deploymentCount, "Nodes", nodeCount, "PVCs", pvcCount, "Events", eventCount, "Restarts", restartCount, "Crashes", crashCount, "TotalCPUUsage (cores)", totalCPUUsage, "TotalMemoryUsage (GB)", totalMemoryUsage, "NodeReady", nodeReady, "NodeMemoryPressure", nodeMemoryPressure, "NodeDiskPressure", nodeDiskPressure)

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
	cronJobGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cronjob_count",
		Help: "Number of cronjobs",
	})
	jobGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "job_count",
		Help: "Number of jobs",
	})
	daemonSetGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "daemonset_count",
		Help: "Number of daemonsets",
	})
	replicaSetGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "replicaset_count",
		Help: "Number of replicasets",
	})
	replicationControllerGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "replicationcontroller_count",
		Help: "Number of replication controllers",
	})
	statefulSetGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "statefulset_count",
		Help: "Number of stateful sets",
	})
	deploymentGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "deployment_count",
		Help: "Number of deployments",
	})
	nodeGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_count",
		Help: "Number of nodes",
	})
	pvcGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pvc_count",
		Help: "Number of PVCs",
	})
	eventGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "event_count",
		Help: "Number of events",
	})
	restartGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "restart_count",
		Help: "Number of container restarts",
	})
	crashGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "crash_count",
		Help: "Number of container crashes",
	})
	cpuUsageGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_cpu_usage",
		Help: "Total CPU usage across all nodes",
	})
	memoryUsageGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_memory_usage",
		Help: "Total memory usage across all nodes",
	})
	nodeReadyGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_ready_count",
		Help: "Number of nodes in ready condition",
	})
	nodeMemoryPressureGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_memory_pressure_count",
		Help: "Number of nodes with memory pressure",
	})
	nodeDiskPressureGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "node_disk_pressure_count",
		Help: "Number of nodes with disk pressure",
	})

	podGauge.Set(float64(podCount))
	serviceGauge.Set(float64(serviceCount))
	configMapGauge.Set(float64(configMapCount))
	secretGauge.Set(float64(secretCount))
	cronJobGauge.Set(float64(cronJobCount))
	jobGauge.Set(float64(jobCount))
	daemonSetGauge.Set(float64(daemonSetCount))
	replicaSetGauge.Set(float64(replicaSetCount))
	replicationControllerGauge.Set(float64(replicationControllerCount))
	statefulSetGauge.Set(float64(statefulSetCount))
	deploymentGauge.Set(float64(deploymentCount))
	nodeGauge.Set(float64(nodeCount))
	pvcGauge.Set(float64(pvcCount))
	eventGauge.Set(float64(eventCount))
	restartGauge.Set(float64(restartCount))
	crashGauge.Set(float64(crashCount))
	cpuUsageGauge.Set(totalCPUUsage)
	memoryUsageGauge.Set(totalMemoryUsage)
	nodeReadyGauge.Set(nodeReady)
	nodeMemoryPressureGauge.Set(nodeMemoryPressure)
	nodeDiskPressureGauge.Set(nodeDiskPressure)

	pusher := push.New(prometheusEndpoint, "kuberesourcesmonitor").
		Collector(podGauge).
		Collector(serviceGauge).
		Collector(configMapGauge).
		Collector(secretGauge).
		Collector(cronJobGauge).
		Collector(jobGauge).
		Collector(daemonSetGauge).
		Collector(replicaSetGauge).
		Collector(replicationControllerGauge).
		Collector(statefulSetGauge).
		Collector(deploymentGauge).
		Collector(nodeGauge).
		Collector(pvcGauge).
		Collector(eventGauge).
		Collector(restartGauge).
		Collector(crashGauge).
		Collector(cpuUsageGauge).
		Collector(memoryUsageGauge).
		Collector(nodeReadyGauge).
		Collector(nodeMemoryPressureGauge).
		Collector(nodeDiskPressureGauge)

	if err := pusher.Push(); err != nil {
		log.Error(err, "Could not push to Prometheus Pushgateway")
	} else {
		log.Info("Successfully pushed metrics to Prometheus Pushgateway")
	}

	log.Info("Successfully reconciled KubeResourcesMonitor")

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
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitorv1alpha1.KubeResourcesMonitor{}).
		Complete(r)
}
