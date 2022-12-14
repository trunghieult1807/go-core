package opentracing

import (
	"strings"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// InitOpentracing initializes opentracing
func InitOpentracing(addr, name string) error {
	cfg := config.Configuration{
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: addr,
			LogSpans:           false,
		},
	}
	_, err := cfg.InitGlobalTracer(
		name,
		config.Logger(jaeger.StdLogger),
		config.Sampler(jaeger.NewConstSampler(true)),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		return err
	}
	return nil
}

// InitOpentracingWithProtocol   initializes opentracing with protocol option
func InitOpentracingWithProtocol(addr, name, protocol string) error {
	cfg := config.Configuration{}

	switch strings.ToLower(protocol) {
	case "http":
		cfg = config.Configuration{
			Reporter: &config.ReporterConfig{
				CollectorEndpoint: addr,
				LogSpans:          false,
			},
		}
	case "udp":
		cfg = config.Configuration{
			Reporter: &config.ReporterConfig{
				LocalAgentHostPort: addr,
				LogSpans:           false,
			},
		}
	}

	_, err := cfg.InitGlobalTracer(
		name,
		config.Logger(jaeger.StdLogger),
		config.Sampler(jaeger.NewConstSampler(true)),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		return err
	}
	return nil
}
