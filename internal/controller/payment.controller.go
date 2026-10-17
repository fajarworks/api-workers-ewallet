package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type PaymentController struct {
	service *service.PaymentService
}

func NewPaymentController(service *service.PaymentService) *PaymentController {
	return &PaymentController{
		service: service,
	}
}

func (p *PaymentController) CreateTopup(ctx *gin.Context) {
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

	var payload dto.CreateTopupRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, dto.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	resp, err := p.service.CreateTopupService(ctx.Request.Context(), userID, payload)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(201, dto.Response{
		Message: "Topup request generated successfully",
		Data:    resp,
	})
}

func (p *PaymentController) DokuNotification(c *gin.Context) {
	var payload dto.DokuNotificationPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, dto.Response{
			Message: "Invalid webhook payload",
			Data:    nil,
		})
		return
	}

	if err := p.service.ProcessDokuNotificationService(c.Request.Context(), payload); err != nil {
		c.JSON(500, dto.Response{
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(200, dto.Response{
		Message: "Notification processed successfully",
		Data:    nil,
	})
}
