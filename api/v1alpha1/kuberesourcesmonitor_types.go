package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubeResourcesMonitorSpec defines the desired state of KubeResourcesMonitor
type KubeResourcesMonitorSpec struct {
    Interval            string `json:"interval,omitempty"`
    PrometheusEndpoint  string `json:"prometheusEndpoint,omitempty"`
    ConfigMapName       string `json:"configMapName,omitempty"` // Add this line
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
