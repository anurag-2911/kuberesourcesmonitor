package controller

import (
    "context"
    "errors"

    monitorv1alpha1 "github.com/anurag-2911/kuberesourcesmonitor/api/v1alpha1"
    "github.com/go-logr/logr"
    "github.com/streadway/amqp"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/types"
)

func (r *KubeResourcesMonitorReconciler) checkAndScaleDeployment(ctx context.Context, mq monitorv1alpha1.MessageQueueSpec, log logr.Logger) error {
    secret := &corev1.Secret{}
    log.Info("namespace used for secret ", "namespace", mq.QueueNamespace)
    if err := r.Get(ctx, types.NamespacedName{Name: mq.QueueSecretName, Namespace: mq.QueueNamespace}, secret); err != nil {
        log.Error(err, "Failed to get Secret")
        return err
    }

    queueUrl, exists := secret.Data[mq.QueueSecretKey]
    if !exists {
        err := errors.New("queue URL not found in secret")
        log.Error(err, "Queue URL key not found in Secret")
        return err
    }

    conn, err := amqp.Dial(string(queueUrl))
    if err != nil {
        log.Error(err, "Failed to connect to RabbitMQ")
        return err
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Error(err, "Failed to open a channel")
        return err
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "message.test.queue",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Error(err, "Failed to declare a queue")
        return err
    }

    queue, err := ch.QueueInspect(q.Name)
    if err != nil {
        log.Error(err, "Failed to inspect the queue")
        return err
    }

    log.Info("Queue length", "messages", queue.Messages)

    deployment := &appsv1.Deployment{}
    if err := r.Get(ctx, types.NamespacedName{Name: mq.DeploymentName, Namespace: mq.DeploymentNamespace}, deployment); err != nil {
        log.Error(err, "Failed to get Deployment")
        return err
    }

    currentReplicas := *deployment.Spec.Replicas

    if queue.Messages > int(mq.ThresholdMessages) && currentReplicas != mq.ScaleUpReplicas {
        log.Info("Scaling up deployment", "deployment", mq.DeploymentName, "replicas", mq.ScaleUpReplicas)
        deployment.Spec.Replicas = &mq.ScaleUpReplicas
    } else if queue.Messages <= int(mq.ThresholdMessages) && currentReplicas != mq.ScaleDownReplicas {
        log.Info("Scaling down deployment", "deployment", mq.DeploymentName, "replicas", mq.ScaleDownReplicas)
        deployment.Spec.Replicas = &mq.ScaleDownReplicas
    }

    if err := r.Update(ctx, deployment); err != nil {
        log.Error(err, "Failed to update Deployment replicas")
        return err
    }

    return nil
}

func rabbitMQAutoScale(ctx context.Context, instance *monitorv1alpha1.KubeResourcesMonitor, r *KubeResourcesMonitorReconciler, log logr.Logger) {
	for _, mq := range instance.Spec.MessageQueues {
		if err := r.checkAndScaleDeployment(ctx, mq, log); err != nil {
			log.Error(err, "Failed to check and scale deployment", "queue", mq.QueueName)
		}
	}
}

