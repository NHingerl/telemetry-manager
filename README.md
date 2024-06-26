# Telemetry Manager

## Status

[![REUSE status](https://api.reuse.software/badge/github.com/kyma-project/telemetry-manager)](https://api.reuse.software/info/github.com/kyma-project/telemetry-manager)

![GitHub tag checks state](https://img.shields.io/github/checks-status/kyma-project/telemetry-manager/main?label=telemetry-manager&link=https%3A%2F%2Fgithub.com%2Fkyma-project%2Ftelemetry-manager%2Fcommits%2Fmain)

## Overview

[Telemetry Manager](docs/user/01-manager.md) is a Kubernetes [operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) that fulfils the [Kyma module interface](https://kyma-project.io/#/06-modules/README.md). It provides APIs for a managed agent/gateway setup for log, trace, and metric ingestion and dispatching into 3rd-party backend systems, in order to reduce the pain of orchestrating such setup on your own. Read more on the [manager](./docs/user/01-manager.md) itself or the general [usage](docs/user/README.md) of the module.

### Logs

The logging controllers generate a Fluent Bit DaemonSet and configuration from one or more LogPipeline and LogParser custom resources. The controllers ensure that all Fluent Bit Pods run the current configuration by restarting Pods after the configuration has changed. See all [CRD attributes](apis/telemetry/v1alpha1/logpipeline_types.go) and some [examples](config/samples).

For more information, see [Logs](./docs/user/02-logs.md).

### Traces

The trace controller creates an [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) deployment and related Kubernetes objects from a `TracePipeline` custom resource. The collector is configured to receive traces using the [OpenTelemetry Protocol (OTLP)](https://opentelemetry.io/docs/specs/otel/protocol/), and forwards the received traces to a configurable OTLP backend.

For more information, see [Traces](./docs/user/03-traces.md).

### Metrics

The metric controller creates an [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) and related Kubernetes objects from a `MetricPipeline` custom resource. The collector is deployed as a [Gateway](https://opentelemetry.io/docs/collector/deployment/#gateway). The controller is configured to receive metrics in the OTLP protocol and forward them to a configurable OTLP backend.

For more information, see [Metrics](./docs/user/04-metrics.md).

## Installation

See the [installation instruction](docs/contributor/installation.md).

## Usage

See the [user documentation](docs/user/README.md).

## Development

For details, see:

- [Available commands for building/linting/installation](docs/contributor/development.md)
- [Testing strategy](docs/contributor/testing.md)
- [Troubleshooting and debugging](docs/contributor/troubleshooting.md)
- [Release process](docs/contributor/releasing.md)
- [Governance checks like linting](docs/contributor/governance.md)

## Contributing

<!--- mandatory section - do not change this! --->

See the [CONTRIBUTING](CONTRIBUTING.md) file.

## Code of Conduct

<!--- mandatory section - do not change this! --->

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## Licensing

<!--- mandatory section - do not change this! --->

See the [LICENSE](LICENSE) file.
