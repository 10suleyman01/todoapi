package apperror

import "github.com/gin-gonic/gin"

type AppError struct {
	Message string
}

func NewAppError(message string) *AppError {
	return &AppError{Message: message}
}

func (err *AppError) Error() string {
	return err.Message
}

func NewJsonMessage(status string, obj interface{}) gin.H {
	return gin.H{
		"status":  status,
		"message": obj,
	}
}
