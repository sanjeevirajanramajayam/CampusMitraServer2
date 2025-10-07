package middleware
import (
	"bitresume/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)
func AuthorizeRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("BITRESUME")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing auth token"})
			return
		}
		claims, err := utils.ParseJWT(cookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found in token"})
			return
		}
		// Check if the role is allowed
		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Set("email", claims["email"])
				c.Set("rollNo", claims["rollNo"])
				c.Set("role", role)
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied for this role"})
	}
}
