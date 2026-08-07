package utils

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponse mimicking Laravel's ApiResponse trait successResponse
func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse mimicking Laravel's ApiResponse trait errorResponse
func ErrorResponse(c *gin.Context, code int, message string, errors interface{}) {
	response := gin.H{
		"success": false,
		"message": message,
	}

	if errors != nil {
		response["errors"] = errors
	}

	c.JSON(code, response)
}

// SuccessResponsePaginated mimicking Laravel's Resource Collection pagination
func SuccessResponsePaginated(c *gin.Context, code int, message string, data interface{}, page, perPage int, total int64) {
	// Calculate total pages
	lastPage := int(total) / perPage
	if int(total)%perPage > 0 {
		lastPage++
	}
	if lastPage == 0 {
		lastPage = 1
	}

	c.JSON(code, gin.H{
		"success": true,
		"message": message,
		"data": gin.H{
			"data":         data,
			"current_page": page,
			"per_page":     perPage,
			"total":        total,
			"last_page":    lastPage,
		},
	})
}
