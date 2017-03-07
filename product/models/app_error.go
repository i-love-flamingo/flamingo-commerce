package models

import "fmt"

type (
	AppError struct {
		Code    int
		Message string
	}
)

func (e *AppError) Error() string {
	return fmt.Sprintf("%d : %s", e.Code, e.Message)
}

func (e *AppError) SetCode(code int) {
	e.Code = code
}

func (e *AppError) SetMessage(message string) {
	e.Message = message
}

func (e *AppError) HasError() bool {
	if &e.Code == nil || e.Code == 0 {
		return false
	}

	return true
}
