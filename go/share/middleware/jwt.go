package middleware

import (
	"biye/share/error_code"
	"biye/share/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    error_code.NotLogin.Code,
				"message": error_code.NotLogin.Message,
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    error_code.InvalidToken.Code,
				"message": error_code.InvalidToken.Message,
			})

			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    error_code.InvalidToken.Code,
				"message": "Token无效或已过期",
			})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserId)
		c.Next()
	}
}
