package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"wc22/repo"
	"wc22/tools"
	"wc22/types"

	"github.com/gin-gonic/gin"
)

func PostAuth(c *gin.Context) {
	var body types.PostAuthRequestBody
	ctx := c.Request.Context()

	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	client, err := tools.GetFirebaseClient(ctx)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusForbidden)
		return
	}

	decoded, err := client.VerifyIDToken(ctx, body.IdToken)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusUnauthorized)
		return
	}

	expiresIn := time.Hour * 24 * 5
	cookie, err := client.SessionCookie(ctx, body.IdToken, expiresIn)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusUnauthorized)
		return
	}

	userRecord, err := client.GetUser(ctx, decoded.UID)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusNotFound)
		return
	}

	fmt.Println(time.Now().UTC())

	user := &types.User{
		Uid:            userRecord.UID,
		Name:           userRecord.DisplayName,
		Email:          userRecord.Email,
		ProfilePicture: userRecord.PhotoURL,
		LastLogin:      time.Now().UTC(),
	}

	err = repo.CreateOrUpdateUser(user)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusNotModified)
		return
	}

	c.SetCookie(
		"authSession",
		cookie,
		int(expiresIn.Seconds()),
		"/",
		"localhost",
		true,
		true,
	)

	accessTokenCookie, err := tools.CreateAccessToken(user)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusUnauthorized)
		return
	}

	c.SetCookie(
		"accessToken",
		accessTokenCookie.Cookie,
		accessTokenCookie.MaxAge,
		"/",
		"localhost",
		true,
		true,
	)

	refreshTokenCookie, err := tools.CreateRefreshToken(user)
	if err == nil {
		c.SetCookie(
			"refreshToken",
			refreshTokenCookie.Cookie,
			refreshTokenCookie.MaxAge,
			"/",
			"localhost",
			true,
			true,
		)
	}

	c.Status(http.StatusOK)
}

func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	user, verificationRes := tools.VerifyJWTSignature(&refreshToken, os.Getenv("JWT_REFRESH_SECRET"))
	if verificationRes != types.JWT_VALID {
		c.SetCookie("accessToken", "", -1, "/", "localhost", true, true)
		c.SetCookie("refreshToken", "", -1, "/", "localhost", true, true)
		c.Status(http.StatusUnauthorized)
		return
	}

	accessTokenCookie, err := tools.CreateAccessToken(user)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusUnauthorized)
		return
	}

	c.SetCookie(
		"accessToken",
		accessTokenCookie.Cookie,
		accessTokenCookie.MaxAge,
		"/",
		"localhost",
		true,
		true,
	)

	refreshTokenCookie, err := tools.CreateRefreshToken(user)
	if err == nil {
		c.SetCookie(
			"refreshToken",
			refreshTokenCookie.Cookie,
			refreshTokenCookie.MaxAge,
			"/",
			"localhost",
			true,
			true,
		)
	}

	c.Status(http.StatusOK)
}

func IsAuthenticated(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, &user)
}
