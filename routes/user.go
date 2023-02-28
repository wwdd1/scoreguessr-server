package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.SetCookie("accessToken", "", -1, "/", "localhost", true, true)
	c.SetCookie("refreshToken", "", -1, "/", "localhost", true, true)
	c.SetCookie("authSession", "", -1, "/", "localhost", true, true)
	c.Status(http.StatusOK)
}
