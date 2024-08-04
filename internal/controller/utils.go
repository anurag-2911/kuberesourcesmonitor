package controller

import (
	"time"

	monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
	"github.com/go-logr/logr"
)

func getInterval(instance *monitorv1alpha1.KubeResourcesMonitor, log logr.Logger) time.Duration {
	timeInterval := instance.Spec.Interval
	interval := 5 * time.Minute

	if timeInterval != "" {
		if parsedInterval, err := time.ParseDuration(timeInterval); err == nil {
			interval = parsedInterval
			log.Info("interval for requeue", "interval", interval)
		} else {
			log.Error(err, "Failed to parse interval, using default")
		}
	}
	return interval
}
