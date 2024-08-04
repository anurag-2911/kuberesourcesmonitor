
# KubeResourcesMonitor

KubeResourcesMonitor is a Kubernetes Operator that provides three key features:

1. **Collection of Kubernetes Resource Metrics**: Collects various counts of Kubernetes resources and makes them available for Prometheus.
2. **Time-based Autoscaling**: Automatically scales specified deployments based on defined time intervals.
3. **Message Queue-based Autoscaling**: Scales consumer microservice deployments based on the number of messages in a RabbitMQ queue.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
  - [Using Helm](#using-helm)
- [Configuration](#configuration)
  - [ConfigMap](#configmap)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

### 1. Collection of Kubernetes Resource Metrics

KubeResourcesMonitor collects various counts of Kubernetes cluster resources such as pods, services, configmaps, secrets, deployments, nodes, etc., and makes these metrics available for Prometheus to scrape.

Metrics Collected
pod_count: Number of pods
service_count: Number of services
configmap_count: Number of configmaps
secret_count: Number of secrets
cronjob_count: Number of cronjobs
deployment_count: Number of deployments
statefulset_count: Number of stateful sets
daemonset_count: Number of daemon sets
replicaset_count: Number of replica sets
endpoint_count: Number of endpoints
persistentvolume_count: Number of persistent volumes
namespace_count: Number of namespaces
node_count: Number of nodes
node_ready_count: Number of nodes in ready condition
node_memory_pressure_count: Number of nodes with memory pressure
node_disk_pressure_count: Number of nodes with disk pressure
pvc_count: Number of PersistentVolumeClaims (PVCs)
event_count: Number of events
restart_count: Number of container restarts
crash_count: Number of container crashes
total_cpu_usage: Total CPU usage across all nodes
total_memory_usage: Total memory usage across all nodes
Using Prometheus Alertmanager
Setting up Alertmanager in Prometheus allows you to track cluster-wide health based on these metrics. You can configure alerts for various metrics to get notified when certain thresholds are crossed. For example:

Pod Count: Alert if the number of pods exceeds a certain limit, indicating potential resource overuse.
Service Count: Monitor the number of services to ensure there are no unexpected spikes.
Node Conditions: Alerts on node_memory_pressure_count and node_disk_pressure_count to track node health.
Container Restarts and Crashes: Alerts on restart_count and crash_count to detect unstable deployments.
Resource Usage: Alerts on total_cpu_usage and total_memory_usage to prevent resource exhaustion.

### 2. Time-based Autoscaling

KubeResourcesMonitor can scale specified deployments based on configured time intervals. This is useful for scaling applications during peak and off-peak hours.

### 3. Message Queue-based Autoscaling

KubeResourcesMonitor can scale consumer microservice deployments based on the number of messages in a RabbitMQ queue. This ensures that your application can handle varying loads efficiently.

## Installation

### Using Helm

To install KubeResourcesMonitor using Helm, follow these steps:

1. Clone the repository:

   
   git clone https://github.com/anurag-2911/kuberesourcesmonitor.git
   
   

2. Install the Helm chart:

   cd kuberesourcesmonitor/charts
   helm upgrade --install krm ./kuberesourcesmonitor --namespace kuberesourcesmonitor --create-namespace --set global.prometheusReleaseLabel=kube-prometheus-stack
   

   This command will install the KubeResourcesMonitor in the `kuberesourcesmonitor` namespace.

## Configuration

KubeResourcesMonitor is configured using a ConfigMap. Below is an example configuration:

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: kuberesourcesmonitor-config
  namespace: kuberesourcesmonitor
data:
  collectMetrics: "true"
  rabbitMQAutoScale: "true"
  timeBasedAutoScale: "true"
```

The configuration options are as follows:

- `collectMetrics`: Enables or disables the collection of Kubernetes resource metrics.
- `rabbitMQAutoScale`: Enables or disables RabbitMQ-based autoscaling.
- `timeBasedAutoScale`: Enables or disables time-based autoscaling.

## Usage

After installing KubeResourcesMonitor and configuring the ConfigMap, the operator will start monitoring the specified resources and perform autoscaling based on the configuration.

You can view the collected metrics by accessing the Prometheus server configured to scrape metrics from KubeResourcesMonitor.

## Contributing

We welcome contributions to KubeResourcesMonitor! If you have any improvements or bug fixes, please submit a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.



This `README.md` provides a professional overview of your project, detailing its features, installation process, configuration options, usage, and contribution guidelines.