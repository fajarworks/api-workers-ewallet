package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-workers-ewallet/internal/dto"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

type TransactionHistoryController struct {
	service *service.TransactionHistoryService
}

func NewTransactionHistoryController(service *service.TransactionHistoryService) *TransactionHistoryController {
	return &TransactionHistoryController{
		service: service,
	}

}

func (c *TransactionHistoryController) GetTransactionHistory(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "7"))
	if err != nil || limit < 1 {
		limit = 7
	}
	search := ctx.Query("search")

	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "token doesn't exist",
		})
		return
	}
	claims, ok := token.(pkg.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "invalid token",
		})
		return
	}
	userId := claims.Id
	qp := dto.TransactionHistoryQuery{
		Search: search,
		Limit:  limit,
		Page:   page,
	}

	trxHistories, metaData, err := c.service.TransactionHistory(ctx.Request.Context(), userId, ctx.Request.URL.Path, qp)
	if err != nil {
		log.Println("[transaction history controller] err: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "an error server occured",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "success get transaction histories",
		Data:    trxHistories,
		Meta:    metaData,
	})
}
