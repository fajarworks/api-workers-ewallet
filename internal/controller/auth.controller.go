package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{
		service: *service,
	}
}

func (a *AuthController) LoginController(ctx *gin.Context) {
	var payload dto.LoginRequest
	if e := ctx.ShouldBindWith(&payload, binding.JSON); e != nil {
		log.Println("error", e.Error())
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}

	user, token, err := a.service.LoginService(ctx.Request.Context(), payload)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "Login Failed. Invalid email or password.",
			Data:    gin.H{},
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "Login Success.",
		Data: gin.H{
			"token": token,
			"user":  user,
		},
	})
}

func (a *AuthController) RegisterController(ctx *gin.Context) {
	var payload dto.LoginRequest

	if e := ctx.ShouldBindWith(&payload, binding.JSON); e != nil {
		log.Println("error", e.Error())
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}

	err := a.service.RegisterService(ctx.Request.Context(), payload)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "Register failed.",
			Data:    gin.H{},
		})
		return
	}

	ctx.JSON(201, dto.Response{
		Success: true,
		Message: "Register Success.",
	})
}
