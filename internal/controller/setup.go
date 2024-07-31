package controller

import (
	"net/http"

	monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KubeResourcesMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Log.Info("SetupWithManager called")

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
