package todo

import (
	"fmt"
	"net/http"
	"todoproject/api/users"
	"todoproject/api/util"
	"todoproject/apperror"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	Id = "id"
)

var (
	GetAllUrl       = "/todos"
	GetByIdUrl      = fmt.Sprintf("/todo/:%s", Id)
	RelativeTodoUrl = "/todo"
)

type Handler struct {
	Storage     *Storage
	userHandler *users.Handler
	Log         *logrus.Logger
}

func NewHandler(storage *Storage, userHandler *users.Handler, log *logrus.Logger) *Handler {
	return &Handler{Storage: storage, userHandler: userHandler, Log: log}
}

func (h *Handler) InitTodoHandler(e *gin.Engine) {
	api := e.Group(util.ApiV1, h.userHandler.IsLogin())
	{
		api.GET(GetAllUrl, h.GetAll)
		api.GET(GetByIdUrl, h.GetById)
		api.POST(RelativeTodoUrl, h.Create)
		api.PUT(RelativeTodoUrl, h.Update)
		api.DELETE(RelativeTodoUrl, h.Delete)
	}
}

func (h *Handler) GetAll(ctx *gin.Context) {
	user := GetCurrentUser(ctx)
	todos, err := h.Storage.GetAllTodoByUserId(ctx, user.Id)
	if err != nil {
		h.Log.Errorf("failed to get todos. due to error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, apperror.NewJsonMessage("success", todos))
}

func (h *Handler) GetById(ctx *gin.Context) {
	id := ctx.Param(Id)
	user := GetCurrentUser(ctx)
	todo, err := h.Storage.GetTodoById(ctx, id, user.Id)
	if err != nil {
		h.Log.Errorf("failed to get todo by id=(%s). due to error: %v", id, err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", "failed to get todo"))
		return
	}
	ctx.JSON(http.StatusOK, apperror.NewJsonMessage("success", todo))
}

func (h *Handler) Create(ctx *gin.Context) {

	var todoDto CreateTodoDto

	if err := ctx.ShouldBindJSON(&todoDto); err != nil {
		h.Log.Errorf("failed to bind todo. due to error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", err.Error()))
		return
	}

	user := GetCurrentUser(ctx)
	todo := Todo{
		Title: todoDto.Title, UserId: user.Id,
	}

	if err := h.Storage.Create(ctx, &todo); err != nil {
		h.Log.Errorf("failed to create todo. due to error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, apperror.NewJsonMessage("success", todo))
}

func (h *Handler) Update(ctx *gin.Context) {
	var todo Todo
	if err := ctx.BindJSON(&todo); err != nil {
		h.Log.Errorf("failed to bind todo. due to error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", err.Error()))
		return
	}

	user := GetCurrentUser(ctx)

	if user.Id != todo.Id {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", "failed to update"))
		return
	}

	if err := h.Storage.Update(ctx, todo); err != nil {
		h.Log.Errorf("failed to update todo. due to error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, apperror.NewJsonMessage("success", todo))
}

func (h *Handler) Delete(ctx *gin.Context) {
	var deleteDto DeleteTodoDto
	if err := ctx.ShouldBindJSON(&deleteDto); err != nil {
		h.Log.Errorf("failed to bind json. due to error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", err.Error()))
		return
	}

	user := GetCurrentUser(ctx)

	if err := h.Storage.Delete(ctx, deleteDto.TodoId, user.Id); err != nil {
		h.Log.Errorf("failed to delete todo. due to error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", err.Error()))
		return
	}
	ctx.JSON(http.StatusNoContent, apperror.NewJsonMessage("success", "deleted"))
}

func GetCurrentUser(ctx *gin.Context) users.User {
	return ctx.MustGet(util.CURRENT_USER).(users.User)
}
