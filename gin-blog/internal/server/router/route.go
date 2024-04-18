package router

import (
	"net/http"
	"strings"

	"blog/docs"
	"blog/internal/server/code"
	"blog/internal/server/router/api"
	v1 "blog/internal/server/router/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rabbit-rm/rabbit/errorKit"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(engin *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	engin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	engin.GET("/auth", api.GetAuth)

}

func parseToken(str string) (*api.AuthClaims, error) {
	token, err := jwt.ParseWithClaims(str, &api.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(api.AuthKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errorKit.New("token is invalid")
	}
	if claims, ok := token.Claims.(*api.AuthClaims); ok {
		return claims, nil
	}
	return nil, errorKit.New("token type(%t) is invalid", token.Claims)
}
