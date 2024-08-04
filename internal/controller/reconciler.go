package controller

import (
	"context"
	"time"

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

	shouldReturn, result, err := r.collectKubResourcesMetrics(ctx, req, instance, log)
	if shouldReturn {
		return result, err
	}

	if len(instance.Spec.Deployments) > 0 {
		timeBasedAutoScale(ctx, instance, r, log)
	}
	// Check RabbitMQ queues and scale deployments if necessary
	for _, mq := range instance.Spec.MessageQueues {
		if err := r.checkAndScaleDeployment(ctx, mq, log); err != nil {
			log.Error(err, "Failed to check and scale deployment", "queue", mq.QueueName)
		}
	}
	log.Info("Successfully reconciled KubeResourcesMonitor")

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

	return reconcile.Result{RequeueAfter: interval}, nil
}
