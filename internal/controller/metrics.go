package controller

import (
	"github.com/prometheus/client_golang/prometheus"
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
