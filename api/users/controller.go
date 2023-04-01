package users

import (
	"fmt"
	"net/http"
	"strings"
	"todoproject/api/util"
	"todoproject/apperror"
	"todoproject/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	Id = "id"
)

var (
	GetOneUser           = fmt.Sprintf("/users/:%s", Id)
	GetByNameAndPassword = "/user"
	LoginUser            = "/users/login"
	LogoutUser           = "/users/logout"
	RelativeUser         = "/users"
)

type Handler struct {
	storage *Storage
	config  *db.Config
	log     *logrus.Logger
}

func NewHandler(storage *Storage, config *db.Config, log *logrus.Logger) *Handler {
	return &Handler{storage: storage, config: config, log: log}
}

func (h *Handler) InitUserHandler(e *gin.Engine) {
	api := e.Group(util.ApiV1)
	{
		api.GET(RelativeUser, h.GetAll)
		api.GET(GetOneUser, h.GetById)
		api.GET(GetByNameAndPassword, h.GetByNameAndPassword)
		api.POST(RelativeUser, h.Create)
		api.POST(LoginUser, h.SignIn)
		api.PUT(RelativeUser, h.IsLogin(), h.Update)
		api.GET(LogoutUser, h.IsLogin(), h.Logout)
		api.DELETE(GetOneUser, h.IsLogin(), h.Delete)
	}
}

func (h *Handler) GetAll(ctx *gin.Context) {
	users, errAll := h.storage.GetAll(ctx)
	if errAll != nil {
		h.log.Error(errAll)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h *Handler) GetById(ctx *gin.Context) {
	user, err := h.storage.GetById(ctx, ctx.Param(Id))
	if err != nil {
		h.log.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) GetByNameAndPassword(ctx *gin.Context) {

	type User struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var u User

	err := ctx.ShouldBind(&u)
	if err != nil {
		return
	}

	fmt.Printf("U = %v", u)

	user, err := h.storage.GetByNameAndPassword(ctx, u.Name, u.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", "failed to get user by name and password"))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) Create(ctx *gin.Context) {
	var userDto CreateDtoUser
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		h.log.Errorf("failed to bind json users. due to error: %v", err)
		return
	}

	user := User{
		Name:     userDto.Name,
		Password: util.GeneratePasswordHash(userDto.Password),
	}
	if err := h.storage.Create(ctx, user); err != nil {
		h.log.Errorf("failed to create users. due to error: %v. method = (Create)", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apperror.NewJsonMessage("fail", "failed to create user"))
		return
	}
	user.CleanPassword()
	ctx.JSON(http.StatusCreated, user)
}

func (h *Handler) Update(ctx *gin.Context) {
	var user User
	errBind := ctx.ShouldBindJSON(&user)
	if errBind != nil {
		h.log.Errorf("failed to bind json users. due to error: %v. method = (Update)", errBind)
		return
	}
	errUpdate := h.storage.Update(ctx, user)
	if errUpdate != nil {
		h.log.Errorf("failed to update users. due to error: %v", errUpdate)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(ctx *gin.Context) {
	errDelete := h.storage.Delete(ctx, ctx.Param(Id))
	if errDelete != nil {
		h.log.Errorf("failed to delete users. due to error: %v", errDelete)
		return
	}
	ctx.JSON(http.StatusNoContent, errDelete)
}

func (h *Handler) SignIn(ctx *gin.Context) {

	var payload SignInDtoUser

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	user, err := h.storage.GetByNameAndPassword(ctx, payload.Name, util.GeneratePasswordHash(payload.Password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid login or password"})
		return
	}

	token, err := util.GenerateToken(h.config.TokenExpires, user.Id, h.config.TokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	h.SetTokenCookie(ctx, token, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

func (h *Handler) IsLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		cookie, err := ctx.Cookie("token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in!"})
			return
		}

		secret := viper.GetString(util.ConfigPath(util.Token, "secret_key"))
		sub, err := util.ValidateToken(token, secret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		user, errFind := h.storage.GetById(ctx, sub.(string))
		if errFind != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		ctx.Set(util.CURRENT_USER, user)
		ctx.Next()

	}
}

func (h *Handler) Logout(ctx *gin.Context) {
	h.SetTokenCookie(ctx, util.EmptyStr, false)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) SetTokenCookie(ctx *gin.Context, token string, isSignIn bool) {
	if isSignIn {
		ctx.SetCookie(util.Token, token, h.config.TokenMaxAge*60, util.HomePath, util.DOMAIN, false, true)
	} else {
		ctx.SetCookie(util.Token, util.EmptyStr, -1, util.HomePath, util.DOMAIN, false, true)
	}
}
