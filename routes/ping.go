package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func Ping(c *gin.Context) {
	c.SetCookie(
		"XSRF-TOKEN",
		csrf.GetToken(c),
		60*60*1000,
		"/",
		"localhost",
		false,
		false,
	)
	c.Status(http.StatusOK)
}
