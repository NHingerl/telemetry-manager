# Telemetry Module

Learn more about the Telemetry Module. Use it to enable observability for your application.

## What Is Telemetry?

Fundamentally, "Observability" is a measure of how well the application's external outputs can reflect the internal states of single components. The insights that an application and the surrounding infrastructure expose are displayed in the form of metrics, traces, and logs - collectively, that's called "telemetry" or ["signals"](https://opentelemetry.io/docs/concepts/signals/). These can be exposed by employing modern instrumentation.

![Stages of Observability](./assets/telemetry-stages.drawio.svg)

1. In order to implement Day-2 operations for a distributed application running in a container runtime, the single components of an application must expose these signals by employing modern instrumentation.
2. Furthermore, the signals must be collected and enriched with the infrastructural metadata in order to ship them to a target system.
3. Instead of providing a one-size-for-all backend solution, the Telemetry module supports you with instrumenting and shipping your telemetry data in a vendor-neutral way.
4. This way, you can conveniently enable observability for your application by integrating it into your existing or desired backends. Pick your favorite among many observability backends (available either as a service or as a self-manageable solution) that focus on different aspects and scenarios.

The Telemetry module focuses exactly on the aspects of instrumentation, collection, and shipment that happen in the runtime and explicitly defocuses on backends.

> [!TIP]
> An enterprise-grade setup demands a central solution outside the cluster, so we recommend in-cluster solutions only for testing purposes. If you want to install lightweight in-cluster backends for demo or development purposes, see [Integration Guides](#integration-guides).

## Features

To support telemetry for your applications, the Telemetry module provides the following features:

- Tooling for collection, filtering, and shipment: Based on the [Open Telemetry Collector](https://opentelemetry.io/docs/collector/) and [Fluent Bit](https://fluentbit.io/), you can configure basic pipelines to filter and ship telemetry data.
- Integration in a vendor-neutral way to a vendor-specific observability system (traces and metrics only): Based on the [OpenTelemetry protocol (OTLP)](https://opentelemetry.io/docs/reference/specification/protocol/), you can integrate backend systems.
- Guidance for the instrumentation (traces and metrics only): Based on [Open Telemetry](https://opentelemetry.io/), you get community samples on how to instrument your code using the [Open Telemetry SDKs](https://opentelemetry.io/docs/instrumentation/) in nearly every programming language.
- Enriching telemetry data by automatically adding common attributes. This is done in compliance with established semantic conventions, ensuring that the enriched data adheres to industry best practices and is more meaningful for analysis. For details, see [Data Enrichment](gateways.md#data-enrichment).
- Opt-out of features for advanced scenarios: At any time, you can opt out for each data type, and use custom tooling to collect and ship the telemetry data.
- SAP BTP as first-class integration: Integration into SAP BTP Observability services, such as SAP Cloud Logging, is prioritized. For more information, see [SAP Cloud Logging](integration/sap-cloud-logging/README.md). <!--- replace with Help Portal link once published? --->

## Scope

The Telemetry module focuses only on the signals of application logs, distributed traces, and metrics. Other kinds of signals are not considered. Also, audit logs are not in scope.

Supported integration scenarios are neutral to the vendor of the target system.

## Architecture

![Components](./assets/telemetry-arch.drawio.svg)

### Telemetry Manager

The Telemetry module ships Telemetry Manager as its core component. Telemetry Manager is a Kubernetes [operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) that implements the Kubernetes controller pattern and manages the whole lifecycle of all other components covered in the Telemetry module. Telemetry Manager watches for the user-created Kubernetes resources: LogPipeline, TracePipeline, and MetricPipeline. In these resources, you specify what data of a signal type to collect and where to ship it.
If Telemetry Manager detects a configuration, it deploys the related gateway and agent components accordingly and keeps them in sync with the requested pipeline definition.

For more information, see [Telemetry Manager](01-manager.md).

### Gateways

The Traces and Metrics features share the common approach of providing a gateway based on the [OTel Collector](https://opentelemetry.io/docs/collector/). It acts as a central point in the cluster to push data in the OTLP format. From here, the data is enriched and filtered and then dispatched as defined in the individual pipeline resources.

For more information, see [Telemetry Gateways](gateways.md).

### Log Agent

The log agent is based on a [Fluent Bit](https://fluentbit.io/) installation running as a [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/). It reads all containers' logs in the runtime and ships them according to a LogPipeline configuration.

For more information, see [Logs](02-logs.md).

For information on the feature supporting OTLP based on an [OTel Collector](https://opentelemetry.io/docs/collector/), see [Logs (OTLP)](logs.md).

### Trace Gateway

The trace gateway is based on an [OTel Collector](https://opentelemetry.io/docs/collector/) [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/). The gateway provides an [OTLP-based](https://opentelemetry.io/docs/reference/specification/protocol/) endpoint to which applications can push the trace signals. According to a TracePipeline configuration, the gateway processes and ships the trace data to a target system.

For more information, see [Traces](03-traces.md) and [Telemetry Gateways](gateways.md).

### Metric Gateway/Agent

The metric gateway and agent are based on an [OTel Collector](https://opentelemetry.io/docs/collector/) [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) and a [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/). The gateway provides an [OTLP](https://opentelemetry.io/docs/reference/specification/protocol/)-based endpoint to which applications can push the metric signals. The agent scrapes annotated Prometheus-based workloads. According to a MetricPipeline configuration, the gateway processes and ships the metric data to a target system.

For more information, see [Metrics](04-metrics.md) and [Telemetry Gateways](gateways.md).

## Integration Guides

To learn about integration with SAP Cloud Logging, read [SAP Cloud Logging](./integration/sap-cloud-logging/README.md). <!--- replace with Help Portal link once published? --->

For integration with other backends, such as Dynatrace, see:

- [Dynatrace](./integration/dynatrace/README.md)
- [Prometheus](./integration/prometheus/README.md)
- [Loki](./integration/loki/README.md)
- [Jaeger](./integration/jaeger/README.md)
- [Amazon CloudWatch and AWS X-Ray](./integration/aws-cloudwatch/README.md)

To learn how to collect data from applications based on the OpenTelemetry SDK, see:

- [OpenTelemetry Demo App](./integration/opentelemetry-demo/README.md)
- [Sample App](./integration/sample-app/)

## API / Custom Resource Definitions

The API of the Telemetry module is based on Kubernetes Custom Resource Definitions (CRD), which extend the Kubernetes API with custom additions. To inspect the specification of the Telemetry module API, see:

- [Telemetry CRD](./resources/01-telemetry.md)
- [LogPipeline CRD](./resources/02-logpipeline.md)
- [TracePipeline CRD](./resources/04-tracepipeline.md)
- [MetricPipeline CRD](./resources/05-metricpipeline.md)

## Resource Usage

To learn more about the resources used by the Telemetry module, see [Kyma Modules' Sizing](https://help.sap.com/docs/btp/sap-business-technology-platform/kyma-modules-sizing#telemetry).
