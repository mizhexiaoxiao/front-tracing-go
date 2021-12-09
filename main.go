package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mizhexiaoxiao/front-tracing-go/common"
	"github.com/mizhexiaoxiao/front-tracing-go/config"
	"github.com/mizhexiaoxiao/front-tracing-go/logger"
	"github.com/mizhexiaoxiao/front-tracing-go/middleware"
	"github.com/mizhexiaoxiao/front-tracing-go/tracing"
)

func main() {
	//configuration initialization
	err := config.Parse()
	if err != nil {
		logger.InfoLogger().Panic("fatal error config file: ", err)
	}
	// start
	logger.InfoLogger().Println("Starting application")
	app := fiber.New()
	app.Use(middleware.CORSConfig())
	tracer, closer, err := tracing.NewJaegerTracer()
	defer closer.Close()

	if err != nil {
		logger.ErrorLogger().Panicf("cannot init Jaeger: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		args := new(common.QueryArgs)
		if err := c.QueryParser(args); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(common.Response{
				Code: fiber.StatusBadRequest,
				Msg:  err.Error(),
			})
		}
		// vaildate data
		errors := common.ValidateStruct(c, *args)
		if errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(common.Response{
				Code: fiber.StatusBadRequest,
				Msg:  "validate error",
				Data: errors,
			})
		}

		if err := tracing.HandleSpan(tracer, c); err != nil {
			logger.InfoLogger().Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.Response{
				Code: fiber.StatusInternalServerError,
				Msg:  err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(common.Response{
			Code: fiber.StatusOK,
			Msg:  "success",
		})
	})

	//time calibration
	app.Head("/timestamp", func(c *fiber.Ctx) error {
		f := c.Response()
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		f.Header.Add("timestamp", timestamp)
		return nil
	})

	// app.Post("/v2", func(c *fiber.Ctx) error {
	// 	body := c.Body()
	// 	ResMapSlice := make([]common.BodyArgs, 1)
	// 	if err := json.Unmarshal(body, &ResMapSlice); err != nil {
	// 		if err != nil {
	// 			return c.Status(fiber.StatusBadRequest).JSON(common.Response{
	// 				Code: fiber.StatusBadRequest,
	// 				Msg:  err.Error(),
	// 				Data: nil,
	// 			})
	// 		}
	// 	}
	// 	for _, value := range ResMapSlice {
	// 		// vaildate data
	// 		errors := common.ValidateStructBody(value)
	// 		if errors != nil {
	// 			return c.Status(fiber.StatusBadRequest).JSON(common.Response{
	// 				Code: fiber.StatusBadRequest,
	// 				Msg:  "validate error",
	// 				Data: errors,
	// 			})
	// 		}
	// 		if err := HandleSpan(tracer, value); err != nil {
	// 			logger.InfoLogger().Println(err)
	// 			return c.Status(fiber.StatusInternalServerError).JSON(common.Response{
	// 				Code: fiber.StatusInternalServerError,
	// 				Msg:  err.Error(),
	// 			})
	// 		}
	// 	}

	// 	return c.Status(fiber.StatusOK).JSON(common.Response{
	// 		Code: fiber.StatusOK,
	// 		Msg:  "success",
	// 	})
	// })

	logger.ErrorLogger().Fatalln(app.Listen(":" + config.FiberPort()))
}

// func spanCtxFromReq(tracer opentracing.Tracer, value common.BodyArgs) (opentracing.SpanContext, error) {
// 	traceHeader := map[string][]string{
// 		"x-b3-traceid":      {value.TraceID},
// 		"x-b3-spanid":       {value.SpanID},
// 		"x-b3-parentspanid": {value.ParentSpanID},
// 		"x-b3-sampled":      {value.Sampled},
// 	}
// 	return tracer.Extract(
// 		opentracing.HTTPHeaders,
// 		opentracing.HTTPHeadersCarrier(traceHeader),
// 	)
// }

// func HandleSpan(tracer opentracing.Tracer, value common.BodyArgs) error {
// 	spanCtx, err := spanCtxFromReq(tracer, value)
// 	if err != nil {
// 		return errors.New("cannot extract spancontext from request headers")
// 	}
// 	api := value.Api
// 	startTime := common.StringToTime(value.StartTime)
// 	finishTime := common.StringToTime(value.FinishTime)
// 	span := tracer.StartSpan(api, opentracing.ChildOf(spanCtx), opentracing.StartTime(startTime))
// 	span.FinishWithOptions(opentracing.FinishOptions{FinishTime: finishTime})
// 	logger.InfoLogger().Println(span, api)
// 	return nil
// }
