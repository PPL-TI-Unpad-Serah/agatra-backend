package middleware

import (
	model "agatra/model"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
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

func AuthAdmin(db *gorm.DB) gin.HandlerFunc {
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
		
		var compare model.User
		err = db.Where("email = ?", claims.Email).First(&compare).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Trouble finding user"})
		}

		if compare.Role == "admin"{
			ctx.Set("email", claims.Email)
			ctx.Next()
		}else{
			ctx.AbortWithStatusJSON(403, gin.H{"error": "Insufficient Permission"})
		}
	})
}

func AuthMaintainer(db *gorm.DB) gin.HandlerFunc {
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
		
		var compare model.User
		err = db.Where("email = ?", claims.Email).First(&compare).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Trouble finding user"})
		}

		if compare.Role == "member"{
			ctx.AbortWithStatusJSON(403, gin.H{"error": "Insufficient Permission"})
		}else{
			ctx.Set("email", claims.Email)
			ctx.Next()
		}
	})
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}