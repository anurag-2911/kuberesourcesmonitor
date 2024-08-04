package controller

import (
	"context"

	monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type KubeResourcesMonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

func (r *KubeResourcesMonitorReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("kuberesourcesmonitor", req.NamespacedName)
	log.Info("Reconcile function called")

	instance := &monitorv1alpha1.KubeResourcesMonitor{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "Failed to get KubeResourcesMonitor")
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.collectKubResourcesMetrics(ctx, req, instance, log); err != nil {
		log.Info("error in getting metrics ", "err", err)
	}

	// scale deployments based on time
	timeBasedAutoScale(ctx, instance, r, log)

	// scale deployments based on high volume of messages in Message Queue
	rabbitMQAutoScale(ctx, instance, r, log)

	log.Info("Successfully reconciled KubeResourcesMonitor")

	interval := getInterval(instance, log)

	return reconcile.Result{RequeueAfter: interval}, nil
}


