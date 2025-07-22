package main

import (
	"chat-server/internal/ws"
	"chat-server/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	key := "mysecretkey"
	maxAge := 86400 * 30
	isProd := false
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	gothic.Store = store
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"), os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"), "http://localhost:8080/auth/google/redirect", "email", "profile"))
	router := gin.Default()
	router.Use(gin.Logger())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization", "Content-Type", "Origin", "Accept", "X-Requested-With")
	config.AddAllowMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	router.Use(cors.New(config))
	h := ws.NewHub()
	go h.Run()
	routes.ChatRoutes(router, h)
	routes.UserRoutes(router)
	routes.SocialRRoutes(router)

	log.Fatal(router.Run(":" + "8080"))
}
