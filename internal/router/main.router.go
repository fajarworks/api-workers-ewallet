package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/rcontrisha/api-workers-ewallet/docs"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MainRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client, mailer *pkg.Mailer) {
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRouter(router, db, rdb, mailer)
	UserRouter(router, db, rdb)
}
