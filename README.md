# KubeResourcesMonitor

KubeResourcesMonitor is a Kubernetes controller designed to collect various counts of Kubernetes resources, provide these metrics to Prometheus, and offer autoscaling capabilities based on time intervals and RabbitMQ queue lengths.

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
   - [1. Collection of Kubernetes Resource Metrics](#1-collection-of-kubernetes-resource-metrics)
   - [2. Time-based Autoscaling](#2-time-based-autoscaling)
   - [3. Message Queue-based Autoscaling](#3-message-queue-based-autoscaling)
3. [Prerequisites](#prerequisites)
4. [Installation](#installation)
   - [Install Helm from Public Repo](#install-helm-from-public-repo)
   - [Port Forward the Requests and Verify the Metrics](#port-forward-the-requests-and-verify-the-metrics)
5. [Configuration](#configuration)
   - [Changing ConfigMap Values](#changing-configmap-values)
6. [Usage](#usage)
7. [Contributing](#contributing)
8. [License](#license)

## Introduction

KubeResourcesMonitor is a Kubernetes controller that provides valuable insights into your cluster's health and performance by collecting metrics and enabling autoscaling based on specific criteria.

## Features

### 1. Collection of Kubernetes Resource Metrics

KubeResourcesMonitor collects various counts of Kubernetes cluster resources such as pods, services, configmaps, secrets, deployments, nodes, etc., and makes these metrics available for Prometheus to scrape.

**Metrics Collected:**
- `pod_count`: Number of pods
- `service_count`: Number of services
- `configmap_count`: Number of configmaps
- `secret_count`: Number of secrets
- `cronjob_count`: Number of cronjobs
- `deployment_count`: Number of deployments
- `statefulset_count`: Number of stateful sets
- `daemonset_count`: Number of daemon sets
- `replicaset_count`: Number of replica sets
- `endpoint_count`: Number of endpoints
- `persistentvolume_count`: Number of persistent volumes
- `namespace_count`: Number of namespaces
- `node_count`: Number of nodes
- `node_ready_count`: Number of nodes in ready condition
- `node_memory_pressure_count`: Number of nodes with memory pressure
- `node_disk_pressure_count`: Number of nodes with disk pressure
- `pvc_count`: Number of PersistentVolumeClaims (PVCs)
- `event_count`: Number of events
- `restart_count`: Number of container restarts
- `crash_count`: Number of container crashes
- `total_cpu_usage`: Total CPU usage across all nodes
- `total_memory_usage`: Total memory usage across all nodes

**Using Prometheus Alertmanager:**
Setting up Alertmanager in Prometheus allows you to track cluster-wide health based on these metrics. You can configure alerts for various metrics to get notified when certain thresholds are crossed. For example:
- **Pod Count**: Alert if the number of pods exceeds a certain limit, indicating potential resource overuse.
- **Service Count**: Monitor the number of services to ensure no unexpected spikes.
- **Node Conditions**: Alerts on `node_memory_pressure_count` and `node_disk_pressure_count` to track node health.
- **Container Restarts and Crashes**: Alerts on `restart_count` and `crash_count` to detect unstable deployments.
- **Resource Usage**: Alerts on `total_cpu_usage` and `total_memory_usage` to prevent resource exhaustion.

### 2. Time-based Autoscaling

KubeResourcesMonitor can scale specified deployments based on configured time intervals. This is useful for scaling applications during peak and off-peak hours.

**Configuration:**
In the CustomResource definition (CRD), you can specify the deployments and their scaling times. Below is an example configuration:

```yaml
apiVersion: monitor.example.com/v1alpha1
kind: KubeResourcesMonitor
metadata:
  name: example-kuberesourcesmonitor
  namespace: kuberesourcesmonitor-system
spec:
  interval: "45s"  # Interval to requeue the reconciler
  prometheusEndpoint: "2112"  # Endpoint for Prometheus metrics
  configMapName: "kuberesourcesmonitor-config"  # Name of the ConfigMap containing metrics configuration
  deployments:
    - name: "busybox-deployment"
      namespace: "default"
      minReplicas: 2
      maxReplicas: 5
      scaleTimes:
        - startTime: "00:00"  # Time to start scaling
          endTime: "06:00"    # Time to end scaling
          replicas: 5         # Number of replicas during this time window
        - startTime: "06:00"
          endTime: "23:59"
          replicas: 2
```

This configuration will scale the `busybox-deployment` to 5 replicas between midnight and 6 AM, and scale it down to 2 replicas for the rest of the day.

### 3. Message Queue-based Autoscaling

KubeResourcesMonitor can scale consumer microservice deployments based on the number of messages in a RabbitMQ queue. This ensures that your application can handle varying loads efficiently.

**Configuration:**
You need to create a secret containing the RabbitMQ URL and specify the message queue details in the CustomResource definition.

**Example Secret:**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: rabbitmq-credentials
  namespace: kuberesourcesmonitor-system
type: Opaque
data:
  # Base64 encode the RabbitMQ URL
  queueUrl: YW1xcDovL25vdmVsbDpub3ZlbGxAMTIzQDE3Mi4xMDUuNTEuMjE2OjU2NzIv
```

**Example Configuration:**

```yaml
apiVersion: monitor.example.com/v1alpha1
kind: KubeResourcesMonitor
metadata:
  name: example-kuberesourcesmonitor
  namespace: kuberesourcesmonitor-system
spec:
  interval: "45s"
  prometheusEndpoint: "2112"
  configMapName: "kuberesourcesmonitor-config"
  messageQueues:
    - queueSecretName: "rabbitmq-credentials"
      queueSecretKey: "queueUrl"
      queueNamespace: "kuberesourcesmonitor-system"
      queueName: "message.test.queue"
      deploymentName: "rabbitmq-consumer"
      deploymentNamespace: "default"
      thresholdMessages: 10
      scaleUpReplicas: 5
      scaleDownReplicas: 1
```

This configuration will scale the `rabbitmq-consumer` deployment based on the number of messages in the `message.test.queue`. If the number of messages exceeds 10, the deployment will scale up to 5 replicas. If the number of messages falls below or equal to 10, the deployment will scale down to 1 replica.

## Prerequisites

- Kubernetes cluster
- Helm installed
- Prometheus installed and configured

## Installation

### Install Helm from Public Repo

1. Add the Helm repository:

    ```sh
    helm repo add krm https://anurag-2911.github.io/helm-charts
    helm repo update
    ```

2. Install the Helm chart in a specific namespace:
    
    kubectl create namespace kuberesourcesmonitor
    helm install krm krm/kuberesourcesmonitor --namespace kuberesourcesmonitor --set global.prometheusReleaseLabel=kube-prometheus-stack 

3. Verify the installation:
 
    kubectl get all -n kuberesourcesmonitor  

4. Uninstall the Helm chart:

    helm uninstall krm --namespace kuberesourcesmonitor
    kubectl delete namespace kuberesourcesmonitor


5. List the Helm chart:

    helm list --namespace kuberesourcesmonitor  

### Port Forward the Requests and Verify the Metrics

1. Get the pods in the namespace:

    kubectl get pods -n kuberesourcesmonitor 

2. Port forward the requests:

    kubectl port-forward pod/kuberesourcesmonitor-kubemntr-f6847f6fd-kmj5q 8080:2112 -n kuberesourcesmonitor

3. Verify the metrics:

    Open your browser and go to `http://localhost:8080/metrics`.

## Configuration

### Changing ConfigMap Values

To change the values in the `ConfigMap`, use the `kubectl patch` command:

kubectl patch configmap kuberesourcesmonitor-config -n kuberesourcesmonitor --type merge -p '{"data":{"collectMetrics":"true","rabbitMQAutoScale":"true","timeBasedAutoScale":"true"}}'


## Usage

### Monitor Kubernetes Resources

KubeResourcesMonitor provides various Kubernetes resource metrics to Prometheus, which can be visualized in Grafana dashboards.

### Autoscaling

- **Time-based Autoscaling:** Configure time intervals for scaling deployments to handle peak and off-peak loads.
- **Message Queue-based Autoscaling:** Scale consumer deployments based on the number of messages in RabbitMQ queues.

## Contributing

Contributions are welcome! Please read the [contributing guidelines](CONTRIBUTING.md) first.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Feel free to copy and paste this content into your README.md file.
