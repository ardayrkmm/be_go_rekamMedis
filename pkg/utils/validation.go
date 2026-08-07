package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HandleValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorsMap := make(map[string][]string)
		for _, e := range validationErrors {
			field := e.Field() // Or use custom logic to convert to snake_case if needed
			errorsMap[field] = []string{"The " + field + " field is invalid."}
		}
		ErrorResponse(c, 422, "The given data was invalid.", errorsMap)
		return
	}
	ErrorResponse(c, 422, err.Error(), nil)
}
