package mw

import (
	"net/http"
	"os"
	"wc22/tools"
	"wc22/types"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("accessToken")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, validationRes := tools.VerifyJWTSignature(&accessToken, os.Getenv("JWT_ACCESS_SECRET"))
		if validationRes == types.JWT_EXPIRED {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", &user)
		c.Next()
	}
}

// func PostAuthGuardMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authSessionCookie, err := c.Cookie("authSession")
// 		if err != nil {
// 			c.Next()
// 			return
// 		}

// 		client, err := tools.GetFirebaseClient(c.Request.Context())
// 		if err != nil {
// 			c.Next()
// 			return
// 		}

// 		decoded, err := client.VerifySessionCookieAndCheckRevoked(c.Request.Context(), authSessionCookie)
// 		if err != nil {
// 			c.Next()
// 			return
// 		}

// 		_ = decoded

// 		c.AbortWithStatus(http.StatusOK)
// 	}
// }
