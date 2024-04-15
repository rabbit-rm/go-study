package router

import (
	"net/http"
	"strings"

	"blog/internal/server/code"
	v1 "blog/internal/server/router/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rabbit-rm/rabbit/errorKit"
)

func InitRouter(engin *gin.Engine) {

	apiV1Group := engin.Group("/api/v1", func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if strings.Compare("", auth) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": code.ErrorAuth,
				"msg":  code.GetMsg(code.ErrorAuth),
			})
			return
		}
		str := strings.TrimPrefix(auth, "Bearer ")
		claims, err := parseToken(str)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": code.ErrorAuthCheckTokenFail,
				"msg":  err.Error(),
			})
			return
		}
		ctx.Set("username", claims.Username)
		ctx.Set("password", claims.Password)
	})
	{
		apiV1Group.GET("/articles", v1.GetArticles)
		apiV1Group.GET("/articles/:id", v1.GetArticle)
		apiV1Group.POST("/articles", v1.AddArticle)
		apiV1Group.PUT("/article/:id", v1.EditArticle)
		apiV1Group.DELETE("/article/:id", v1.DeleteArticle)

		apiV1Group.GET("/tags/:id", v1.GetTag)
		apiV1Group.POST("/tags", v1.AddTag)
		apiV1Group.PUT("/tags/:id", v1.EditTag)
		apiV1Group.DELETE("/tags/:id", v1.DeleteTag)

	}
	engin.GET("/auth", v1.GetAuth)

}

func parseToken(str string) (*v1.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(str, &v1.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(v1.AuthKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errorKit.New("token is invalid")
	}
	if claims, ok := token.Claims.(*v1.AuthClaims); ok {
		return claims, nil
	}
	return nil, errorKit.New("token type(%t) is invalid", token.Claims)
}
