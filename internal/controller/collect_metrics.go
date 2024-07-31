package controller

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
)

func (r *KubeResourcesMonitorReconciler) collectKubResourcesMetrics(ctx context.Context, req reconcile.Request, instance *monitorv1alpha1.KubeResourcesMonitor, log logr.Logger) (bool, reconcile.Result, error) {
	configMap := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: instance.Spec.ConfigMapName}, configMap)
	if err != nil {
		log.Error(err, "Failed to get ConfigMap", "name", instance.Spec.ConfigMapName)
		return true, reconcile.Result{}, err
	}
	metricsToCollect := strings.Split(configMap.Data["metrics"], "\n")
	log.Info("Metrics to collect from ConfigMap", "metrics", metricsToCollect)

	// Collect and log various metrics...
	return false, reconcile.Result{}, nil
}
