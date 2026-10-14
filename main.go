package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rcontrisha/api-workers-ewallet/internal/config"
	"github.com/rcontrisha/api-workers-ewallet/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error oading env file")
	}

	pdb := config.NewPsqlDb(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	pool, err := pdb.Connect()
	if err != nil {
		log.Println("Cannot Connect to DB\nReason: ", err.Error())
		return
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Println("Database is not ready\nReason: ", err.Error())
		return
	}
	log.Println("Database is ready")

	r := gin.Default()
	router.MainRouter(r, pool)
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
