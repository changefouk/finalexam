package customer

import (
	"fmt"
	"net/http"

	"github.com/changefouk/finalexam/database"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func createCustomer(c *gin.Context) {
	cus := Customer{}
	if err := c.ShouldBindJSON(&cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := database.Connect().QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3)  RETURNING id", cus.Name, cus.Email, cus.Status)

	err := row.Scan(&cus.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, cus)
}

func getCustomerByID(c *gin.Context) {
	id := c.Param("id")

	stmt, err := database.Connect().Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	cus := &Customer{}

	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cus)
}

func getCustomer(c *gin.Context) {
	stmt, err := database.Connect().Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	customers := []*Customer{}

	for rows.Next() {
		cus := &Customer{}

		err = rows.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		customers = append(customers, cus)
	}

	c.JSON(http.StatusOK, customers)
}

func updateCustomer(c *gin.Context) {
	// get customer from database
	id := c.Param("id")

	stmt, err := database.Connect().Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	cus := &Customer{}

	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := c.ShouldBindJSON(cus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update to database
	stmt, err = database.Connect().Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	_, err = stmt.Exec(id, cus.Name, cus.Email, cus.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cus)
}

func deleteCustomer(c *gin.Context) {
	id := c.Param("id")
	stmt, err := database.Connect().Prepare("DELETE FROM customers WHERE id = $1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := stmt.Exec(id); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "customer deleted",
	})
}

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
