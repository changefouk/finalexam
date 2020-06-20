package customer

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func authMiddleware(c *gin.Context) {
	fmt.Println("start #middleware")
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you don't have Access"})
		c.Abort()
		return
	}

	c.Next()

	fmt.Println("end #middleware")
}
