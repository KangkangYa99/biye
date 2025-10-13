package middleware

import (
	"biye/share/error_code"
	"biye/share/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(error_code.NotLogin)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(error_code.InvalidToken)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.Error(error_code.InvalidToken)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserId)
		c.Next()
	}

}
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last().Err
			if myerr, ok := lastErr.(*error_code.APIError); ok {
				httpCode := getHTTPStatus(myerr.Code)
				c.JSON(httpCode, gin.H{
					"success":    false,
					"error_code": myerr.Code,
					"message":    myerr.Message,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": error_code.ServerErrorCode,
				"message":    error_code.ServerError.Message,
			})
		}
	}
}
func getHTTPStatus(code int) int {
	const (
		ServerErrorCodeBase = 10000
		UserErrorCodeBase   = 20000
		DeviceErrorCodeBase = 30000
	)
	if code == error_code.NotLoginCode || code == error_code.InvalidTokenCode {
		return http.StatusUnauthorized
	} else if code >= UserErrorCodeBase && code < DeviceErrorCodeBase {
		return http.StatusUnauthorized
	} else if code >= DeviceErrorCodeBase && code < 40000 {
		return http.StatusForbidden
	} else if code >= ServerErrorCodeBase && code < UserErrorCodeBase {
		return http.StatusInternalServerError
	} else {
		return http.StatusBadRequest
	}
}
