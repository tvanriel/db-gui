package httpinterface

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func getStatusError(status int) func(err string) (int, any) {

	return func(err string) (int, any) {
		return status, ErrorResponse{
			Message: fmt.Sprintf("%s, %v", http.StatusText(status), err),
		}

	}
}

func ValidationError(err error) (int, any) {
	return getStatusError(http.StatusUnprocessableEntity)(err.Error())
}

func InternalServerError(err error) (int, any) {
	return getStatusError(http.StatusInternalServerError)(err.Error())
}

func Unauthorized(str string) (int, any) {
	return getStatusError(http.StatusUnauthorized)(str)
}

func NotImplemented() (int, any) {
	return getStatusError(http.StatusNotImplemented)("Not implemented")
}
