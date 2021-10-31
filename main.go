package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mizhexiaoxiao/front-tracing-go/common"
	"github.com/mizhexiaoxiao/front-tracing-go/config"
	"github.com/mizhexiaoxiao/front-tracing-go/logger"
	"github.com/mizhexiaoxiao/front-tracing-go/tracing"
)

func main() {
	//configuration initialization
	err := config.Parse()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// start
	logger.InfoLogger().Println("Starting application")
	app := fiber.New()
	tracer, closer, err := tracing.NewJaegerTracer()
	defer closer.Close()

	if err != nil {
		logger.ErrorLogger().Panicf("cannot init Jaeger: %v", err)
	}
	app.Get("/", func(c *fiber.Ctx) error {
		args := new(common.QueryArgs)
		if err := c.QueryParser(args); err != nil {
			return err
		}
		// vaildate data
		errors := common.ValidateStruct(c, *args)
		if errors != nil {
			return c.JSON(common.Response{
				Code: 401,
				Msg:  "validate error",
				Data: errors,
			})
		}

		startTime, _ := StringToTime(args.StartTime)
		finishTime, _ := StringToTime(args.FinishTime)
		if err := tracing.HandleSpan(tracer, c, args.Api, startTime, finishTime); err != nil {
			logger.InfoLogger().Println(err)
			return c.JSON(common.Response{
				Code: 500,
				Msg:  err.Error(),
			})
		}

		return c.JSON(common.Response{
			Code: 200,
			Msg:  "success",
		})
	})
	logger.ErrorLogger().Fatalln(app.Listen(":" + config.FiberPort()))

}

func StringToTime(s string) (t time.Time, err error) {
	data, err := strconv.ParseInt(s, 10, 64)
	// 19位时间戳(ns)转为time类型
	t = time.Unix(0, data)
	return
}
