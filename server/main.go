package main

import (
	"chat-server/internal/ws"
	"chat-server/routes"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	h := ws.NewHub()
	go h.Run()
	routes.ChatRoutes(router, h)
	routes.UserRoutes(router)
	log.Fatal(router.Run(":" + "8080"))
}
