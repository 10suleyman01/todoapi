package handler

import "github.com/gin-gonic/gin"

type ControllerHandler interface {
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}
