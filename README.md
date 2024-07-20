#This project builds a Kubernetes Operator to monitor Kubernetes resources like pods, services, configMaps, cronjobs, secrets, etc, and on a configurable period pushes the metrics to the Prometheus push gateway from there Prometheus can be configured to monitor the metrics sent by the operator. The operator pod logs the same info like the number of pods in all the namespaces, the number of nodes, CPU usage, memory usage, cronJobs, etc.
## Prerequisite
   ##Linux OS, I have used Ubuntu (Ubuntu 22.04.4 LTS) for developing this project
   ##Go run time installed, for this project, I have used go version go1.21.0 linux/amd64
   ##Git
   ##Kubectl
   ##Operator-SDK: operator-sdk version
      (operator-sdk version: "v1.35.0", commit: "e95abdbd5ccb7ca0fd586e0c6f578e491b0a025b", kubernetes version: "v1.28.0", go version: "go1.21.0",          GOOS:    "linux", GOARCH: "amd64")
      Note: Make sure the Go version is the same as the one mentioned in the operator SDK version used.
      (https://sdk.operatorframework.io/docs/installation)
   ##Dockers


# Kubernetes Operator to Monitor Cluster Resources and Raise Alerts in Prometheus Alert Manager

This Kubernetes Operator monitors cluster resources like pods, services, config maps, secrets, and nodes. It raises alerts in Prometheus Alert Manager based on specified thresholds.

## Step 1: Initialize the Operator

1. **Create a new directory for the project:**
   
    mkdir kuberesourcesmonitor
    cd kuberesourcesmonitor
  

2. **Initialize a new Operator SDK project:**
   
    operator-sdk init --domain=example.com --repo=github.com/anurag-2911/kuberesourcesmonitor
   

3. **Create a new API:**
    bash
    operator-sdk create api --group=monitor --version=v1alpha1 --kind=KubeResourcesMonitor --resource --controller
    

## Step 2: Define the CRD

Edit the `api/v1alpha1/kuberesourcesmonitor_types.go` to define the CRD structure:


package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubeResourcesMonitorSpec defines the desired state of KubeResourcesMonitor
type KubeResourcesMonitorSpec struct {
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


## Step 3: Implement the Controller

Edit `controllers/kuberesourcesmonitor_controller.go` to implement the reconciliation logic.

## Step 4: Build and Push the Operator Image

1. Build the Docker image:   
    make docker-build IMG=anurag2911/kuberesourcesmonitor:latest
  

2. Push the Docker image:
    make docker-push IMG=anurag2911/kuberesourcesmonitor:latest
   

## Step 5: Deploy the Operator

1. Generate the manifests:
    make manifests   

2. Apply the CRD:    
    kubectl apply -f config/crd/bases   

3. Deploy the Operator:    
    make deploy IMG=anurag2911/kuberesourcesmonitor:latest
    

## Step 6: Create an Instance of the CRD

Create a YAML file for the custom resource:

apiVersion: monitor.example.com/v1alpha1
kind: KubeResourcesMonitor
metadata:
  name: kuberesourcesmonitor-sample
spec:
  interval: "5m"


Apply the custom resource:
kubectl apply -f kuberesourcesmonitor-sample.yaml

## Step 7: Verify the Operator

Check the logs of the Operator:

kubectl logs -f deployment/kuberesourcesmonitor-controller-manager -c manager

You should see logs indicating the counts of pods, services, config maps, and secrets being fetched and pushed to Prometheus.

## Step 8: Set Up Prometheus Alert Manager

Set up Prometheus Alert Manager to receive and manage alerts raised by the Kubernetes Operator.
