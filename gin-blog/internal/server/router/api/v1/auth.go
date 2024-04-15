package v1

import (
	"net/http"
	"time"

	"blog/internal/server/code"
	"blog/internal/server/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rabbit-rm/rabbit/validateKit"
)

const AuthKey = "rabbit@rm99@gmail.com"

type AuthClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

func GetAuth(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	if err := validateKit.ValidateVar(username, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "username is required",
		})
		return
	}
	if err := validateKit.ValidateVar(password, "required"); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.Error,
			"msg":  "password is required",
		})
		return
	}
	auth0, err := models.GetAuthByUsername(username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": code.ErrorAuth,
			"msg":  "username is invalid",
		})
		return
	}
	if auth0.Password != password {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": code.ErrorAuth,
			"msg":  "password is invalid",
		})
		return
	}
	token, err := generateToken(username, password)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code.ErrorAuth,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.Success,
		"msg":  code.GetMsg(code.Success),
		"data": gin.H{
			"token": token,
		},
	})
}

func generateToken(username string, password string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		Username: username,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}).SignedString([]byte(AuthKey))
}
