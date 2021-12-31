package main

import (
	"encoding/json"
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

	//time calibration
	app.Head("/timestamp", func(c *fiber.Ctx) error {
		f := c.Response()
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		f.Header.Add("Access-Control-Expose-Headers", "timestamp")
		f.Header.Add("timestamp", timestamp)
		return nil
	})

	app.Post("/v2", func(c *fiber.Ctx) error {
		body := c.Body()
		serviceName := c.Query("servicename", "front")
		serviceName = serviceName + ".front"
		tracer, closer, err := tracing.NewJaegerTracer(serviceName)
		defer closer.Close()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(common.Response{
				Code: fiber.StatusInternalServerError,
				Msg:  err.Error(),
				Data: nil,
			})
		}

		ResMapSlice := make([]common.BodyArgs, 20)
		if err := json.Unmarshal(body, &ResMapSlice); err != nil {
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(common.Response{
					Code: fiber.StatusBadRequest,
					Msg:  err.Error(),
					Data: nil,
				})
			}
		}
		for _, value := range ResMapSlice {
			// vaildate data
			errors := common.ValidateStructBody(value)
			if errors != nil {
				return c.Status(fiber.StatusBadRequest).JSON(common.Response{
					Code: fiber.StatusBadRequest,
					Msg:  "validate error",
					Data: errors,
				})
			}
			if err := tracing.V2HandleSpan(tracer, value); err != nil {
				logger.InfoLogger().Println(err)
				return c.Status(fiber.StatusInternalServerError).JSON(common.Response{
					Code: fiber.StatusInternalServerError,
					Msg:  err.Error(),
				})
			}
		}

		return c.Status(fiber.StatusOK).JSON(common.Response{
			Code: fiber.StatusOK,
			Msg:  "success",
		})
	})

	logger.ErrorLogger().Fatalln(app.Listen(":" + config.FiberPort()))
}
