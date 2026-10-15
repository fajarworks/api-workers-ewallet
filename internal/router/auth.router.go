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

	repo := repository.NewAuthRepo()
	service := service.NewAuthService(repo, db)
	controller := controller.NewAuthController(service)

	authRouter.POST("login", controller.LoginController)
	authRouter.POST("register", controller.RegisterController)
}
