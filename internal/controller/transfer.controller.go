package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type TransferController struct {
	service *service.TransferService
}

func NewTransferController(service *service.TransferService) *TransferController {
	return &TransferController{
		service: service,
	}
}

func (t *TransferController) SearchUsers(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "token doesn't exist",
		})
		return
	}
	claims, ok := token.(pkg.JWTClaims)
	if !ok {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "invalid token",
		})
		return
	}
	userID := claims.Id

	query := ctx.Query("q")
	result, err := t.service.SearchUsersService(ctx.Request.Context(), query, userID)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "Users fetched successfully",
		Data:    result,
	})
}

func (t *TransferController) TransferMoney(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "token doesn't exist",
		})
		return
	}
	claims, ok := token.(pkg.JWTClaims)
	if !ok {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "invalid token",
		})
		return
	}
	userID := claims.Id

	var payload dto.TransferRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	resp, err := t.service.TransferMoneyService(ctx.Request.Context(), userID, payload)
	if err != nil {
		ctx.JSON(400, dto.Response{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "Transfer completed successfully",
		Data:    resp,
	})
}
