package controller

import (
	"context"
	"time"

	monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *KubeResourcesMonitorReconciler) scaleDeployment(ctx context.Context, deploymentName string, replicas int32, namespace string, log logr.Logger) {
	deployment := &appsv1.Deployment{}
	err := r.Client.Get(ctx, client.ObjectKey{Name: deploymentName, Namespace: namespace}, deployment)
	if err != nil {
		log.Error(err, "Failed to get Deployment for scaling", "name", deploymentName, "namespace", namespace)
		return
	}

	if *deployment.Spec.Replicas != replicas {
		deployment.Spec.Replicas = &replicas
		err = r.Client.Update(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment replicas", "name", deploymentName, "namespace", namespace)
		} else {
			log.Info("Scaled deployment", "name", deploymentName, "namespace", namespace, "new replicas", replicas)
		}
	}
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
