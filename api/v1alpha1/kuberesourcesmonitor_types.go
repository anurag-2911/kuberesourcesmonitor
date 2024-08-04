package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentScalingSpec defines the scaling schedule for a specific deployment
type DeploymentScalingSpec struct {
	Name        string      `json:"name"`        // Name of the deployment to scale
	Namespace   string      `json:"namespace"`   // Namespace of the deployment
	MinReplicas int32       `json:"minReplicas"` // Minimum number of replicas
	MaxReplicas int32       `json:"maxReplicas"` // Maximum number of replicas
	ScaleTimes  []ScaleTime `json:"scaleTimes"`  // List of times to scale the deployment
}

// ScaleTime defines a time window for scaling
type ScaleTime struct {
	StartTime string `json:"startTime"` // Start time in HH:MM format (24-hour)
	EndTime   string `json:"endTime"`   // End time in HH:MM format (24-hour)
	Replicas  int32  `json:"replicas"`  // Number of replicas during this time window
}

// MessageQueueSpec defines the message queue monitoring configuration
type MessageQueueSpec struct {
	QueueSecretName     string `json:"queueSecretName"`     // Name of the secret containing the queue URL
	QueueSecretKey      string `json:"queueSecretKey"`      // Key in the secret that holds the queue URL
	QueueNamespace      string `json:"queueNamespace"`      // Namespace of the secret
	QueueName           string `json:"queueName"`           // Name of the message queue (e.g., RabbitMQ, AWS SQS)
	DeploymentName      string `json:"deploymentName"`      // Name of the deployment to scale
	DeploymentNamespace string `json:"deploymentNamespace"` // Namespace of the deployment to scale
	ThresholdMessages   int32  `json:"thresholdMessages"`   // Threshold of messages to trigger scaling
	ScaleUpReplicas     int32  `json:"scaleUpReplicas"`     // Number of replicas to scale up
	ScaleDownReplicas   int32  `json:"scaleDownReplicas"`   // Number of replicas to scale down
}

// KubeResourcesMonitorSpec defines the desired state of KubeResourcesMonitor
type KubeResourcesMonitorSpec struct {
	Interval           string                  `json:"interval,omitempty"`
	PrometheusEndpoint string                  `json:"prometheusEndpoint,omitempty"`
	ConfigMapName      string                  `json:"configMapName,omitempty"`
	Deployments        []DeploymentScalingSpec `json:"deployments"`   // List of deployments to scale
	MessageQueues      []MessageQueueSpec      `json:"messageQueues"` // List of message queues to monitor
}

// KubeResourcesMonitorStatus defines the observed state of KubeResourcesMonitor
type KubeResourcesMonitorStatus struct {
	// Add custom fields here
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KubeResourcesMonitor is the Schema for the kuberesourcesmonitors API
type KubeResourcesMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeResourcesMonitorSpec   `json:"spec,omitempty"`
	Status KubeResourcesMonitorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubeResourcesMonitorList contains a list of KubeResourcesMonitor
type KubeResourcesMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeResourcesMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeResourcesMonitor{}, &KubeResourcesMonitorList{})
}
