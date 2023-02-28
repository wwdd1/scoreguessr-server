package main

import (
	"log"
	"os"
	"time"
	"wc22/routes"
	"wc22/routes/mw"
	"wc22/tools"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	csrf "github.com/utrack/gin-csrf"
)

type AuthRequestBody struct {
	IdToken string `json:"idToken"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	gin.SetMode(gin.DebugMode)
	tools.Init()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://localhost:3000", "http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Cookie",
			"X-XSRF-TOKEN",
			"X-Version",
			"Access-Control-Allow-Credentials",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	store := cookie.NewStore([]byte(os.Getenv("CSRF_SESSION_SECRET")))
	router.Use(sessions.Sessions("userSession", store))
	router.Use(csrf.Middleware(csrf.Options{
		Secret: os.Getenv("CSRF_SESSION_SECRET"),
		ErrorFunc: func(c *gin.Context) {
			c.Status(404)
			c.Abort()
		},
	}))

	auth := router.Group("/auth")
	{
		auth.POST("", routes.PostAuth)
		auth.GET("/user", mw.AuthMiddleware(), routes.IsAuthenticated)
		auth.GET("/refresh", routes.RefreshToken)
	}

	user := router.Group("/user")
	{
		user.GET("", mw.AuthMiddleware(), func(c *gin.Context) {
			c.Status(200)
		})

		user.POST("", mw.AuthMiddleware(), func(c *gin.Context) {
			c.Status(200)
		})

		user.POST("/logout", mw.AuthMiddleware(), routes.Logout)
	}

	router.GET("/ping", routes.Ping)

	router.Run()
}
