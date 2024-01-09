package middleware

import (
	"net/http"

	"rest-apishka/internal/http/repository"
	"rest-apishka/internal/model"

	"github.com/gin-gonic/gin"
)

func ModeratorOnly(r *repository.Repository, c *gin.Context) bool {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		c.Abort()
	}

	userID := ctxUserID.(uint)

	role, err := r.GetUserRoleByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}

	if role == model.ModeratorRole {
		return true
	}
	return false
}
