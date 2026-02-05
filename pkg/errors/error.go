package errors

import (
	goerrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code  int    `json:"code" example:"400"`
	Error string `json:"error" example:"Descrição do erro"`
}

func ErrorHandler(err error, c *gin.Context) (result ErrorResponse) {
	type statusCoder interface{ StatusCode() int }

	var sc statusCoder
	if goerrors.As(err, &sc) {
		result.Code = sc.StatusCode()
		result.Error = err.Error()
	} else {

		result.Code = http.StatusBadRequest
		result.Error = "bad request"
	}

	c.JSON(result.Code, result)
	c.Abort()
	return
}
