package tracing

import (
	"errors"
	"io"

	"github.com/mizhexiaoxiao/front-tracing-go/common"
	"github.com/mizhexiaoxiao/front-tracing-go/config"
	"github.com/mizhexiaoxiao/front-tracing-go/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
)

func NewJaegerTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	collectorEndpoint := config.JaegerCollectorEndpoint()
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Gen128Bit: true,
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          false,
			CollectorEndpoint: collectorEndpoint,
		},
	}

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
	)

	return tracer, closer, err
}

func spanCtxFromReq(tracer opentracing.Tracer, value common.BodyArgs) (opentracing.SpanContext, error) {
	traceHeader := map[string][]string{
		"x-b3-traceid":      {value.TraceID},
		"x-b3-spanid":       {value.SpanID},
		"x-b3-parentspanid": {value.ParentSpanID},
		"x-b3-sampled":      {value.Sampled},
	}
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(traceHeader),
	)
}

func ExtractData(tracer opentracing.Tracer, value common.BodyArgs) error {
	spanCtx, err := spanCtxFromReq(tracer, value)
	if err != nil {
		return errors.New("cannot extract spancontext from request headers")
	}

	api := value.Api
	domain := value.Domain
	timing := value.Timing

	span := HandleSpan(tracer, api, spanCtx, value.StartTime, value.FinishTime)

	for k, v := range timing {
		HandleSpan(tracer, k, span.Context(), v.StartTime, v.FinishTime)
	}

	span.SetTag("domian", domain)

	logger.InfoLogger().Println(span, api)
	return nil
}

func HandleSpan(
	tracer opentracing.Tracer,
	operationName string,
	spanCtx opentracing.SpanContext,
	startTime string,
	finishTime string,
) opentracing.Span {
	st := common.StringToTime(startTime)
	ft := common.StringToTime(finishTime)
	span := tracer.StartSpan(operationName, opentracing.ChildOf(spanCtx), opentracing.StartTime(st))
	span.FinishWithOptions(opentracing.FinishOptions{FinishTime: ft})
	return span
}
