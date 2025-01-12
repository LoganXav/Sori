package validators

import (
	"LoganXav/sori/app/structs"
	"LoganXav/sori/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Validate payload for Quality Control job
//
//	param c *fiber.Ctx
//	return error
func JobsQualityControl(c *fiber.Ctx) error {
	var errors []*structs.IError
	body := new(structs.JobsQualityControlRequest) 
	errb := c.BodyParser(&body)
	if errb != nil {
		return helpers.UnprocessableResponse(c, nil, "Invalid request format")
	}

	// Validate the request body
	err := Validator.Struct(body)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el structs.IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		return helpers.UnprocessableResponse(c, errors, "Unprocessable entity")
	}

	return c.Next()
}
