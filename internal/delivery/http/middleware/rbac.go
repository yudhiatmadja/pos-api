package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the user has one of the allowed roles.
// Assumes AuthMiddleware has already set "roles" (slice of strings) or "role" (string) in context.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Support both "role" (legacy/single) and "roles" (plural)
		var userRoles []string

		if rolesVal, exists := c.Get("roles"); exists {
			if rolesSlice, ok := rolesVal.([]string); ok {
				userRoles = rolesSlice
			}
		}

		// Fallback to single role if "roles" not present or strictly single
		if len(userRoles) == 0 {
			if roleVal, exists := c.Get("role"); exists {
				if roleStr, ok := roleVal.(string); ok {
					userRoles = append(userRoles, roleStr)
				}
			}
		}

		if len(userRoles) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found in context"})
			return
		}

		for _, allowed := range allowedRoles {
			for _, userRole := range userRoles {
				if userRole == allowed {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
	}
}
