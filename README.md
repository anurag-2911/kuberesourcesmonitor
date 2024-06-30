## Kubernetes Operator to monitor Kubernetes cluster resources like pods, services, config maps, secrets, nodes and raise alerts in Prometheus Alert Manager
## based on the specified threshold defined.

Step 1: Initialize the Operator

1. Create a new directory for the project:
    
    mkdir kuberesourcesmonitor
    cd kuberesourcesmonitor
    

2. Initialize a new Operator SDK project:
    
    operator-sdk init --domain=example.com --repo=github.com/anurag-2911/kuberesourcesmonitor
    

3. Create a new API:
    
    operator-sdk create api --group=monitor --version=v1alpha1 --kind=KubeResourcesMonitor --resource --controller
    

Step 2: Define the CRD

Edit the `api/v1alpha1/kuberesourcesmonitor_types.go` to define the CRD structure:

go
package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubeResourcesMonitorSpec defines the desired state of KubeResourcesMonitor
type KubeResourcesMonitorSpec struct {
    // Add custom fields here
    Interval string `json:"interval,omitempty"`
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

Step 3: Implement the Controller

Edit `controllers/kuberesourcesmonitor_controller.go` to implement the reconciliation logic:


Step 4: Build and Push the Operator Image

1. Build the Docker image:
    
    make docker-build IMG=anurag2911/kuberesourcesmonitor:latest
    

2. Push the Docker image:
    
    make docker-pu IMG=anurag2911/kuberesourcesmonitor:latest
    

 Step 5: Deploy the Operator

1. Generate the manifests:
    
    make manifests
    

2. Apply the CRD:
    
    kubectl apply -f config/crd/bases
    

3. Deploy the Operator:
    
    make deploy IMG=anurag2911/kuberesourcesmonitor:latest
    

Step 6: Create an Instance of the CRD

Create a YAML file for the custom resource:

yaml
apiVersion: monitor.example.com/v1alpha1
kind: KubeResourcesMonitor
metadata:
  name: kuberesourcesmonitor-sample
spec:
  interval: "5m"


Apply the custom resource:


kubectl apply -f kuberesourcesmonitor-sample.yaml


 Step 7: Verify the Operator

Check the logs of the Operator:


kubectl logs -f deployment/kuberesourcesmonitor-controller-manager -c manager


should see logs indicating the counts of pods, services, config maps, and secrets being fetched and pushed to Prometheus.

Step 8: Set Up Prometheus Alert Manager


