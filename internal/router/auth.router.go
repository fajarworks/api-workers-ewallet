package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-workers-ewallet/internal/controller"
	"github.com/rcontrisha/api-workers-ewallet/internal/repository"
	"github.com/rcontrisha/api-workers-ewallet/internal/service"
)

func AuthRouter(r *gin.Engine, db *pgxpool.Pool) {
	authRouter := r.Group("auth")

	repo := repository.NewAuthRepo(db)
	service := service.NewAuthService(repo)
	controller := controller.NewAuthController(service)

	authRouter.POST("login", controller.LoginController)
}
