package customer

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	engine := gin.Default()

	engine.Use(authMiddleware)

	engine.POST("/customers", createCustomer)
	engine.GET("/customers/:id", getCustomerByID)
	engine.GET("/customers", getCustomer)
	engine.PUT("/customers/:id", updateCustomer)
	engine.DELETE("/customers/:id")

	return engine
}
