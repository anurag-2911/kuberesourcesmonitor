package controller

import (
	"context"
	monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

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
	statefulSetGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "statefulset_count",
		Help: "Number of stateful sets",
	})
	daemonSetGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "daemonset_count",
		Help: "Number of daemon sets",
	})
	replicaSetGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "replicaset_count",
		Help: "Number of replica sets",
	})	
	endpointGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "endpoint_count",
		Help: "Number of endpoints",
	})
	persistentVolumeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "persistentvolume_count",
		Help: "Number of persistent volumes",
	})
	namespaceGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "namespace_count",
		Help: "Number of namespaces",
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
	prometheus.MustRegister(podGauge, serviceGauge, configMapGauge, secretGauge, cronJobGauge, deploymentGauge, statefulSetGauge,
		daemonSetGauge, replicaSetGauge, endpointGauge, persistentVolumeGauge, namespaceGauge, nodeGauge,
		nodeReadyGauge, nodeMemoryPressureGauge, nodeDiskPressureGauge, pvcGauge, eventGauge, restartGauge, crashGauge,
		cpuUsageGauge, memoryUsageGauge)
}

func (r *KubeResourcesMonitorReconciler) collectKubResourcesMetrics(ctx context.Context, req reconcile.Request, instance *monitorv1alpha1.KubeResourcesMonitor, log logr.Logger) error {
	configMap := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: instance.Spec.ConfigMapName}, configMap)
	if err != nil {
		log.Error(err, "Failed to get ConfigMap", "name", instance.Spec.ConfigMapName)
		return err
	}
	metricsToCollect := strings.Split(configMap.Data["metrics"], "\n")
	log.Info("Metrics to collect from ConfigMap", "metrics", metricsToCollect)

	// Collect and log various metrics...
	// Collect various metrics from the cluster
	podList := &corev1.PodList{}
	err = r.Client.List(ctx, podList)
	if err != nil {
		log.Error(err, "Failed to list Pods")
		return err
	}
	podCount := len(podList.Items)
	log.Info("Collected pod metrics", "count", podCount)

	serviceList := &corev1.ServiceList{}
	err = r.Client.List(ctx, serviceList)
	if err != nil {
		log.Error(err, "Failed to list Services")
		return err
	}
	serviceCount := len(serviceList.Items)
	log.Info("Collected service metrics", "count", serviceCount)

	configMapList := &corev1.ConfigMapList{}
	err = r.Client.List(ctx, configMapList)
	if err != nil {
		log.Error(err, "Failed to list ConfigMaps")
		return err
	}
	configMapCount := len(configMapList.Items)
	log.Info("Collected config map metrics", "count", configMapCount)

	secretList := &corev1.SecretList{}
	err = r.Client.List(ctx, secretList)
	if err != nil {
		log.Error(err, "Failed to list Secrets")
		return err
	}
	secretCount := len(secretList.Items)
	log.Info("Collected secret metrics", "count", secretCount)

	cronJobList := &batchv1.CronJobList{}
	err = r.Client.List(ctx, cronJobList)
	if err != nil {
		log.Error(err, "Failed to list CronJobs")
		return err
	}
	cronJobCount := len(cronJobList.Items)
	log.Info("Collected cronjob metrics", "count", cronJobCount)

	deploymentList := &appsv1.DeploymentList{}
	err = r.Client.List(ctx, deploymentList)
	if err != nil {
		log.Error(err, "Failed to list Deployments")
		return err
	}
	deploymentCount := len(deploymentList.Items)
	log.Info("Collected deployment metrics", "count", deploymentCount)

	statefulSetList := &appsv1.StatefulSetList{}
	err = r.Client.List(ctx, statefulSetList)
	if err != nil {
		log.Error(err, "Failed to list StatefulSets")
		return err
	}
	statefulSetCount := len(statefulSetList.Items)
	log.Info("Collected statefulset metrics", "count", statefulSetCount)

	daemonSetList := &appsv1.DaemonSetList{}
	err = r.Client.List(ctx, daemonSetList)
	if err != nil {
		log.Error(err, "Failed to list DaemonSets")
		return err
	}
	daemonSetCount := len(daemonSetList.Items)
	log.Info("Collected daemonset metrics", "count", daemonSetCount)

	replicaSetList := &appsv1.ReplicaSetList{}
	err = r.Client.List(ctx, replicaSetList)
	if err != nil {
		log.Error(err, "Failed to list ReplicaSets")
		return err
	}
	replicaSetCount := len(replicaSetList.Items)
	log.Info("Collected replicaset metrics", "count", replicaSetCount)

	endpointList := &corev1.EndpointsList{}
	err = r.Client.List(ctx, endpointList)
	if err != nil {
		log.Error(err, "Failed to list Endpoints")
		return err
	}
	endpointCount := len(endpointList.Items)
	log.Info("Collected endpoint metrics", "count", endpointCount)

	persistentVolumeList := &corev1.PersistentVolumeList{}
	err = r.Client.List(ctx, persistentVolumeList)
	if err != nil {
		log.Error(err, "Failed to list PersistentVolumes")
		return err
	}
	persistentVolumeCount := len(persistentVolumeList.Items)
	log.Info("Collected persistentvolume metrics", "count", persistentVolumeCount)

	namespaceList := &corev1.NamespaceList{}
	err = r.Client.List(ctx, namespaceList)
	if err != nil {
		log.Error(err, "Failed to list Namespaces")
		return err
	}
	namespaceCount := len(namespaceList.Items)
	log.Info("Collected namespace metrics", "count", namespaceCount)

	nodeList := &corev1.NodeList{}
	err = r.Client.List(ctx, nodeList)
	if err != nil {
		log.Error(err, "Failed to list Nodes")
		return err
	}
	nodeCount := len(nodeList.Items)
	log.Info("Collected node metrics", "count", nodeCount)

	pvcList := &corev1.PersistentVolumeClaimList{}
	err = r.Client.List(ctx, pvcList)
	if err != nil {
		log.Error(err, "Failed to list PVCs")
		return err
	}
	pvcCount := len(pvcList.Items)
	log.Info("Collected PVC metrics", "count", pvcCount)

	eventList := &corev1.EventList{}
	err = r.Client.List(ctx, eventList)
	if err != nil {
		log.Error(err, "Failed to list Events")
		return err
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
	log.Info("Resource counts", "Pods", podCount, "Services", serviceCount, "ConfigMaps", configMapCount, "Secrets", secretCount, "CronJobs", cronJobCount, "Deployments", deploymentCount, "StatefulSets", statefulSetCount, "DaemonSets", daemonSetCount, "ReplicaSets", replicaSetCount, "Endpoints", endpointCount, "PersistentVolumes", persistentVolumeCount, "Namespaces", namespaceCount, "Nodes", nodeCount, "PVCs", pvcCount, "Events", eventCount, "Restarts", restartCount, "Crashes", crashCount, "TotalCPUUsage (cores)", totalCPUUsage, "TotalMemoryUsage (GB)", totalMemoryUsage, "NodeReady", nodeReady, "NodeMemoryPressure", nodeMemoryPressure, "NodeDiskPressure", nodeDiskPressure)

	// Update Prometheus gauge values with the collected metrics
	log.Info("Setting Prometheus gauge values")
	podGauge.Set(float64(podCount))
	serviceGauge.Set(float64(serviceCount))
	configMapGauge.Set(float64(configMapCount))
	secretGauge.Set(float64(secretCount))
	cronJobGauge.Set(float64(cronJobCount))
	deploymentGauge.Set(float64(deploymentCount))
	statefulSetGauge.Set(float64(statefulSetCount))
	daemonSetGauge.Set(float64(daemonSetCount))
	replicaSetGauge.Set(float64(replicaSetCount))	
	endpointGauge.Set(float64(endpointCount))
	persistentVolumeGauge.Set(float64(persistentVolumeCount))
	namespaceGauge.Set(float64(namespaceCount))
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
	log.Info("Set Prometheus gauge values")

	return nil
}
