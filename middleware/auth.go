package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	model "agatra/model"
	jwt "github.com/golang-jwt/jwt/v4"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		data, err := ctx.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				if ctx.GetHeader("Content-Type") == "application/json" {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				} else {
					ctx.Redirect(http.StatusSeeOther, "/user/login")
					ctx.Abort()
				}
				return
			}
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		claims := &model.Claims{}

		tkn, err := jwt.ParseWithClaims(data, claims, func(token *jwt.Token) (interface{}, error) {
            return model.JwtKey, nil
        })
        if err != nil || !tkn.Valid {
            ctx.JSON(400, model.ErrorResponse{Error: "ga valid bang"})
            return
        }
		ctx.Set("email", claims.Email)
		ctx.Next()
	})
}
