package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/rcontrisha/api-workers-ewallet/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MainRouter(router *gin.Engine, db *pgxpool.Pool) {
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRouter(router, db)
}
