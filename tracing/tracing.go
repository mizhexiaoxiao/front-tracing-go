package tracing

import (
	"errors"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mizhexiaoxiao/front-tracing-go/config"
	"github.com/mizhexiaoxiao/front-tracing-go/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
)

func NewJaegerTracer() (opentracing.Tracer, io.Closer, error) {
	collectorEndpoint := config.JaegerCollectorEndpoint()
	servicename := config.JaegerServiceName()

	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Gen128Bit: true,
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			CollectorEndpoint:   collectorEndpoint,
		},
	}

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	tracer, closer, err := cfg.NewTracer(
		//jaegercfg.Reporter(reporter),
		jaegercfg.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
	)

	return tracer, closer, err
}

func getTraceHeaderFromReq(c *fiber.Ctx) opentracing.HTTPHeadersCarrier {
	traceHeader := map[string][]string{
		"x-b3-traceid":      {c.Get("x-b3-traceid")},
		"x-b3-spanid":       {c.Get("x-b3-spanid")},
		"x-b3-parentspanid": {c.Get("x-b3-parentspanid")},
		"x-b3-sampled":      {c.Get("x-b3-sampled")},
	}
	return traceHeader
}

func spanCtxFromReq(tracer opentracing.Tracer, c *fiber.Ctx) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(getTraceHeaderFromReq(c)),
	)
}

func HandleSpan(tracer opentracing.Tracer, c *fiber.Ctx, operationName string, startTime time.Time, finishTime time.Time) error {
	spanCtx, err := spanCtxFromReq(tracer, c)
	if err != nil {
		return errors.New("cannot extract spancontext from request headers")
	}
	span := tracer.StartSpan(operationName, opentracing.ChildOf(spanCtx), opentracing.StartTime(startTime))
	logger.InfoLogger().Println(span)
	span.FinishWithOptions(opentracing.FinishOptions{FinishTime: finishTime})
	return nil
}
