package traceSvc

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/afret0/wheel/log"
	"github.com/afret0/wheel/tool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
)

func Init(cfg *viper.Viper, svc string) {

	ctx, cancel := context.WithTimeout(tool.NewCtxBK(), time.Second*5)
	defer cancel()

	lg := log.CtxLogger(ctx).WithFields(logrus.Fields{"svc": svc})
	if !tool.EnvEnabled("TRACE") {
		lg.Infof("Trace disabled")
		return
	}

	endpoint := cfg.GetString("trace.endpoint")
	token := cfg.GetString("trace.token")
	lg.Infof("token: %s, endpoint: %s", token, endpoint)

	headers := map[string]string{
		"Authentication": token,
	}

	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithHeaders(headers),
	)
	if err != nil {
		lg.Errorf("otlptracegrpc.New failed, err: %s", err)
		return
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(svc),
		),
	)
	if err != nil {
		lg.Errorf("resource.New failed, err: %s", err)
		return
	}

	fraction := 1.00
	if ef := os.Getenv("TRACE-FRACTION"); ef != "" {
		eff, err := strconv.ParseFloat(ef, 64)
		if err != nil {
			lg.Errorf("strconv.ParseFloat failed, err: %s", err)
		} else {
			fraction = eff
		}
	}

	sampler := sdktrace.TraceIDRatioBased(fraction)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
		trace.WithSampler(sampler),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

}
