package common

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type ValErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(c *fiber.Ctx, args QueryArgs) []*ValErrorResponse {
	var errors []*ValErrorResponse
	validate := validator.New()
	err := validate.Struct(args)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateStructBody(body BodyArgs) []*ValErrorResponse {
	var errors []*ValErrorResponse
	validate := validator.New()
	err := validate.Struct(body)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
