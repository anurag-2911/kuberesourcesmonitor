package controller

import (
	"context"
	"strconv"

	monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
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

	configMap := &corev1.ConfigMap{}
	err = r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: instance.Spec.ConfigMapName}, configMap)
	if err != nil {
		log.Error(err, "Failed to get ConfigMap", "name", instance.Spec.ConfigMapName)
		return reconcile.Result{}, err
	}

	isCollectMetrics, err := strconv.ParseBool(configMap.Data["collectMetrics"])
	if err != nil {
		isCollectMetrics = false
	}
	isRabbitMQAutoScale, err := strconv.ParseBool(configMap.Data["rabbitMQAutoScale"])
	if err != nil {
		isRabbitMQAutoScale = false
	}
	isTimeBasedAutoScale, err := strconv.ParseBool(configMap.Data["timeBasedAutoScale"])
	if err != nil {
		isTimeBasedAutoScale = false
	}

	log.Info("ConfigMap values", "collectMetrics", isCollectMetrics, "rabbitMQAutoScale", isRabbitMQAutoScale, "timeBasedAutoScale", isTimeBasedAutoScale)

	if isCollectMetrics {
		// Adds metrics for different resources, logs it, and provides it to Prometheus to scrape
		collectMetrics(ctx, r, req, instance, log)
	}

	if isRabbitMQAutoScale {
		rabbitMQAutoScale(ctx, instance, r, log)
	}

	if isTimeBasedAutoScale {
		timeBasedAutoScale(ctx, instance, r, log)
	}

	log.Info("Successfully reconciled this iteration")

	interval := getInterval(instance, log)

	return reconcile.Result{RequeueAfter: interval}, nil
}
