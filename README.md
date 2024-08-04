
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

#### Metrics Collected

- **pod_count**: Number of pods
- **service_count**: Number of services
- **configmap_count**: Number of configmaps
- **secret_count**: Number of secrets
- **cronjob_count**: Number of cronjobs
- **deployment_count**: Number of deployments
- **statefulset_count**: Number of stateful sets
- **daemonset_count**: Number of daemon sets
- **replicaset_count**: Number of replica sets
- **endpoint_count**: Number of endpoints
- **persistentvolume_count**: Number of persistent volumes
- **namespace_count**: Number of namespaces
- **node_count**: Number of nodes
- **node_ready_count**: Number of nodes in ready condition
- **node_memory_pressure_count**: Number of nodes with memory pressure
- **node_disk_pressure_count**: Number of nodes with disk pressure
- **pvc_count**: Number of PersistentVolumeClaims (PVCs)
- **event_count**: Number of events
- **restart_count**: Number of container restarts
- **crash_count**: Number of container crashes
- **total_cpu_usage**: Total CPU usage across all nodes
- **total_memory_usage**: Total memory usage across all nodes

#### Using Prometheus Alertmanager

Setting up Alertmanager in Prometheus allows you to track cluster-wide health based on these metrics. You can configure alerts for various metrics to get notified when certain thresholds are crossed. For example:

- **Pod Count**: Alert if the number of pods exceeds a certain limit, indicating potential resource overuse.
- **Service Count**: Monitor the number of services to ensure no unexpected spikes.
- **Node Conditions**: Alerts on `node_memory_pressure_count` and `node_disk_pressure_count` to track node health.
- **Container Restarts and Crashes**: Alerts on `restart_count` and `crash_count` to detect unstable deployments.
- **Resource Usage**: Alerts on `total_cpu_usage` and `total_memory_usage` to prevent resource exhaustion.

Here is an example Alertmanager configuration to alert on high CPU usage:

```yaml
groups:
- name: kubernetes-resources
  rules:
  - alert: HighCPUUsage
    expr: total_cpu_usage > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High CPU usage detected"
      description: "Total CPU usage across all nodes has exceeded 80 cores for more than 5 minutes."
```





### 2. Time-based Autoscaling

KubeResourcesMonitor can scale specified deployments based on configured time intervals. This feature is particularly useful for scaling applications during peak and off-peak hours, ensuring optimal resource utilization and cost efficiency.

#### Benefits

- **Resource Optimization**: Automatically scale up resources during peak hours to handle increased load and scale down during off-peak hours to save on costs.
- **Performance Management**: Maintain application performance by ensuring that the appropriate number of instances are running based on anticipated load.
- **Cost Efficiency**: Reduce costs by minimizing the number of running instances during periods of low demand.

#### Configuration

To configure time-based autoscaling, specify the `deployments` section in your Custom Resource Definition (CRD). Here's an example configuration:

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

#### How to Configure Time-based Autoscaling

1. **Define the CRD**: Ensure the CRD includes fields for specifying deployment names, namespaces, minimum and maximum replicas, and scale times.
2. **Specify Deployment Details**: In the `deployments` section, provide details of the deployments you wish to scale, including:
   - `name`: The name of the deployment.
   - `namespace`: The namespace where the deployment resides.
   - `minReplicas`: The minimum number of replicas.
   - `maxReplicas`: The maximum number of replicas.
   - `scaleTimes`: An array of time windows specifying the start and end times, and the number of replicas to scale to during those times.
3. **Apply the Configuration**: Use `kubectl` to apply the CRD configuration to your Kubernetes cluster.


kubectl apply customresource.yaml


#### Example Use Cases

- **E-commerce Applications**: Scale up during peak hours and scale down during the night when traffic is low.
- **Data Processing**: Increase the number of workers during data ingestion periods and reduce them during processing off-hours.
- **Seasonal Traffic**: Automatically adjust resources during seasonal peaks, such as holiday sales or promotional events.

By leveraging time-based autoscaling, you can ensure that your applications have the resources to handle varying workloads efficiently, maintaining performance while optimizing costs.

### Message Queue-based Autoscaling

KubeResourcesMonitor can scale consumer microservice deployments based on the number of messages in a RabbitMQ queue. This dynamic scaling ensures that your application can efficiently handle varying loads by automatically adjusting the number of consumer instances based on the queue length.

#### Benefits

- **Efficient Load Handling**: Automatically scale up the number of consumer instances when there are many messages in the queue to process them faster and scale down when the load decreases.
- **Resource Optimization**: Utilize resources more effectively by scaling down consumers during periods of low traffic, reducing operational costs.
- **Improved Performance**: Ensure timely processing of messages, which can be critical for applications that rely on real-time data processing.

#### Configuration

To configure RabbitMQ-based autoscaling, you need to define the `messageQueues` section in your Custom Resource Definition (CRD) and use a Kubernetes Secret to store the RabbitMQ URL securely.

##### Example Secret Configuration

Create a Kubernetes Secret to store the RabbitMQ connection details. This ensures that sensitive information is not exposed in your configuration files.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: rabbitmq-credentials
  namespace: kuberesourcesmonitor-system
type: Opaque
data:
  # Base64 encode the RabbitMQ URL
  queueUrl: YW1zcDovLc5vdmVsbDpuc3ZlbGxAMTIzQDE3Mi4xMDUcNTEuMjE2OjU2NzIv
```

The `queueUrl` is the base64 encoded RabbitMQ URL. To encode your RabbitMQ URL, you can use the following command:

echo -n 'amqp://user:password@hostname:port/vhost' | base64


##### Example CRD Configuration

Define the `messageQueues` section in your CRD to specify how the consumer deployment should scale based on the message queue length.

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
  messageQueues:
    - queueSecretName: "rabbitmq-credentials"
      queueSecretKey: "queueUrl"
      queueNamespace: "kuberesourcesmonitor-system"  # Namespace where the Secret is located
      queueName: "message.test.queue"
      deploymentName: "rabbitmq-consumer"
      deploymentNamespace: "default"
      thresholdMessages: 10
      scaleUpReplicas: 5
      scaleDownReplicas: 1
    # Add more message queues as needed
```

##### How to Configure Message Queue-based Autoscaling

1. **Create the Secret**: Store the RabbitMQ URL in a Kubernetes Secret.

kubectl apply -f config/samples/secret.yaml

2. **Define the CRD**: Include the `messageQueues` section with the necessary details for the RabbitMQ queue and the deployment to be scaled.

kubectl apply -f config/samples/customresource.yaml

3. **Apply the Configuration**: Use `kubectl` to apply the CRD configuration to your Kubernetes cluster.

kubectl apply -f config/samples/customresource.yaml

#### Example Use Cases

- **Real-time Data Processing**: Scale up the number of consumers during periods of high message traffic to process data faster.
- **Event-driven Applications**: Ensure that your event-driven architecture can handle spikes in event generation by dynamically scaling consumers.
- **Batch Processing**: Efficiently handle batch processing by scaling consumers based on the number of messages queued for processing.


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
