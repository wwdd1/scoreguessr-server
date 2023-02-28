package tools

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"wc22/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func CreateUUID() string {
	return uuid.NewString()
}

func CreateAccessToken(user *types.User) (*types.JWTCookieInfo, error) {
	expiry, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRE"))
	if err != nil {
		expiry = 60
	}

	expireDate := time.Now().Add(time.Minute * time.Duration(expiry))
	// expireDate := time.Now().Add(time.Minute)
	issuedDate := time.Now()

	claims := jwt.MapClaims{}
	claims["sig"] = CreateUUID()
	claims["aud"] = os.Getenv("JWT_AUD")
	claims["iss"] = os.Getenv("JWT_ISS")
	claims["iat"] = issuedDate.UTC().Unix()
	claims["exp"] = expireDate.UTC().Unix()
	claims["sub"] = user.Uid
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["profilePicture"] = user.ProfilePicture

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return &types.JWTCookieInfo{
		Cookie: token,
		MaxAge: int(expireDate.UTC().Sub(time.Now().UTC()).Milliseconds() / 1000),
		// MaxAge: 3600,
	}, nil
}

func CreateRefreshToken(user *types.User) (*types.JWTCookieInfo, error) {
	expiry, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRE"))
	if err != nil {
		expiry = 60 * 24 * 7
	}

	expireDate := time.Now().Add(time.Minute * time.Duration(expiry))
	issuedDate := time.Now()

	claims := jwt.MapClaims{}
	claims["sig"] = CreateUUID()
	claims["aud"] = os.Getenv("JWT_AUD")
	claims["iss"] = os.Getenv("JWT_ISS")
	claims["iat"] = issuedDate.UTC().Unix()
	claims["exp"] = expireDate.UTC().Unix()
	claims["sub"] = user.Uid
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["profilePicture"] = user.ProfilePicture

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return &types.JWTCookieInfo{
		Cookie: token,
		MaxAge: int(expireDate.UTC().Sub(time.Now().UTC()).Milliseconds() / 1000),
	}, nil
}

func VerifyJWTSignature(token *string, secret string) (*types.User, int) {
	result, err := jwt.ParseWithClaims(*token, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	var claims *types.JWTClaims
	var ok bool
	if claims, ok = result.Claims.(*types.JWTClaims); !ok || err != nil || !result.Valid {
		log.Println(err.Error())
		isExpired := strings.Contains(err.Error(), "token is expired")
		if isExpired {
			return nil, types.JWT_EXPIRED
		}
		return nil, types.JWT_INVALID
	}

	if !claims.VerifyAudience(os.Getenv("JWT_AUD"), true) || claims.Issuer != os.Getenv("JWT_ISS") {
		return nil, types.JWT_INVALID
	}

	return &types.User{
		Name:           claims.Name,
		Email:          claims.Email,
		Uid:            claims.Subject,
		ProfilePicture: claims.ProfilePicture,
	}, types.JWT_VALID
}
