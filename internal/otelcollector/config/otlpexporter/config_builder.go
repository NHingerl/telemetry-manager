package otlpexporter

import (
	"context"
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"

	telemetryv1alpha1 "github.com/kyma-project/telemetry-manager/apis/telemetry/v1alpha1"
	"github.com/kyma-project/telemetry-manager/internal/otelcollector/config"
	sharedtypesutils "github.com/kyma-project/telemetry-manager/internal/utils/sharedtypes"
)

type EnvVars map[string][]byte

const (
	SignalTypeMetric = "metric"
	SignalTypeTrace  = "trace"
	SignalTypeLog    = "log"
)

type ConfigBuilder struct {
	reader       client.Reader
	otlpOutput   *telemetryv1alpha1.OTLPOutput
	pipelineName string
	queueSize    int
	signalType   string
}

func NewConfigBuilder(reader client.Reader, otlpOutput *telemetryv1alpha1.OTLPOutput, pipelineName string, queueSize int, signalType string) *ConfigBuilder {
	return &ConfigBuilder{
		reader:       reader,
		otlpOutput:   otlpOutput,
		pipelineName: pipelineName,
		queueSize:    queueSize,
		signalType:   signalType,
	}
}

func (cb *ConfigBuilder) MakeConfig(ctx context.Context) (*config.OTLPExporter, EnvVars, error) {
	envVars, err := makeEnvVars(ctx, cb.reader, cb.otlpOutput, cb.pipelineName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make env vars: %w", err)
	}

	exportersConfig := makeExportersConfig(cb.otlpOutput, cb.pipelineName, envVars, cb.queueSize, cb.signalType)

	return exportersConfig, envVars, nil
}

func makeExportersConfig(otlpOutput *telemetryv1alpha1.OTLPOutput, pipelineName string, envVars map[string][]byte, queueSize int, signalType string) *config.OTLPExporter {
	headers := makeHeaders(otlpOutput, pipelineName)
	otlpEndpointVariable := makeOTLPEndpointVariable(pipelineName)
	otlpEndpointValue := string(envVars[otlpEndpointVariable])
	tlsConfig := makeTLSConfig(otlpOutput, otlpEndpointValue, pipelineName)

	sendingQueue := config.SendingQueue{
		Enabled: false,
	}
	if queueSize != 0 {
		sendingQueue.QueueSize = queueSize
		sendingQueue.Enabled = true
	}

	otlpExporterConfig := config.OTLPExporter{
		Endpoint:     fmt.Sprintf("${%s}", otlpEndpointVariable),
		Headers:      headers,
		TLS:          tlsConfig,
		SendingQueue: sendingQueue,
		RetryOnFailure: config.RetryOnFailure{
			Enabled:         true,
			InitialInterval: "5s",
			MaxInterval:     "30s",
			MaxElapsedTime:  "300s",
		},
	}

	if len(otlpOutput.Path) > 0 && SignalTypeMetric == signalType {
		otlpExporterConfig.Endpoint = ""
		otlpExporterConfig.MetricsEndpoint = fmt.Sprintf("${%s}", otlpEndpointVariable)
	}

	if len(otlpOutput.Path) > 0 && SignalTypeTrace == signalType {
		otlpExporterConfig.Endpoint = ""
		otlpExporterConfig.TracesEndpoint = fmt.Sprintf("${%s}", otlpEndpointVariable)
	}

	return &otlpExporterConfig
}

func ExporterID(protocol string, pipelineName string) string {
	var outputType string
	if protocol == telemetryv1alpha1.OTLPProtocolHTTP {
		outputType = "otlphttp"
	} else {
		outputType = "otlp"
	}

	return fmt.Sprintf("%s/%s", outputType, pipelineName)
}

func makeTLSConfig(output *telemetryv1alpha1.OTLPOutput, otlpEndpointValue, pipelineName string) config.TLS {
	var cfg config.TLS

	cfg.Insecure = isInsecureOutput(otlpEndpointValue)

	if output.TLS == nil {
		return cfg
	}

	if !cfg.Insecure {
		cfg.Insecure = output.TLS.Insecure
	}

	cfg.InsecureSkipVerify = output.TLS.InsecureSkipVerify
	if sharedtypesutils.IsValid(output.TLS.CA) {
		cfg.CAPem = fmt.Sprintf("${%s}", makeTLSCaVariable(pipelineName))
	}

	if sharedtypesutils.IsValid(output.TLS.Cert) {
		cfg.CertPem = fmt.Sprintf("${%s}", makeTLSCertVariable(pipelineName))
	}

	if sharedtypesutils.IsValid(output.TLS.Key) {
		cfg.KeyPem = fmt.Sprintf("${%s}", makeTLSKeyVariable(pipelineName))
	}

	return cfg
}

func makeHeaders(output *telemetryv1alpha1.OTLPOutput, pipelineName string) map[string]string {
	headers := make(map[string]string)

	if output.Authentication != nil && sharedtypesutils.IsValid(&output.Authentication.Basic.User) && sharedtypesutils.IsValid(&output.Authentication.Basic.Password) {
		basicAuthHeaderVariable := makeBasicAuthHeaderVariable(pipelineName)
		headers["Authorization"] = fmt.Sprintf("${%s}", basicAuthHeaderVariable)
	}

	for _, header := range output.Headers {
		headers[header.Name] = fmt.Sprintf("${%s}", makeHeaderVariable(header, pipelineName))
	}

	return headers
}

func isInsecureOutput(endpoint string) bool {
	return len(strings.TrimSpace(endpoint)) > 0 && strings.HasPrefix(endpoint, "http://")
}
