package helpers

import (

	"github.com/labstack/echo/v4"
)

// JSONResponse standardizes the JSON response format and supports server messages
func JSONResponse(c echo.Context, statusCode int, status bool, message string, data interface{}) error {
	response := map[string]interface{}{
		"data":    data,
		"status":  status,
		"message": message,
	}
	return c.JSON(statusCode, response)
}

