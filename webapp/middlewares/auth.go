package middlewares

import (
	"net/http"
	"strings"

	"webapp-go/webapp/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const USER_ID_KEY = "user_id"

func AuthRequired(bearerService services.BearerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userId := session.Get(USER_ID_KEY)
		if userId != nil {
			userId, err := uuid.Parse(userId.(string))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
				return
			}

			c.Set(USER_ID_KEY, userId)
			c.Next()
			return
		}

		authorization := c.GetHeader("Authorization")
		if authorization != "" {
			tokens := strings.Split(authorization, " ")

			if len(tokens) != 2 || tokens[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
				return
			}

			accessToken := tokens[1]

			userId, err := bearerService.ValidateToken(accessToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
				return
			}

			c.Set(USER_ID_KEY, userId)

			c.Next()
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		c.Abort()
		return
	}
}
