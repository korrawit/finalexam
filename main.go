package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korrawit/finalexam/customer"
	"github.com/korrawit/finalexam/database"
	"github.com/korrawit/finalexam/repository"

	_ "github.com/lib/pq"
)

func main() {
	// Initializing
	db := database.PostgresDB{}
	repo := repository.Repository{
		DB: db,
	}
	customerContext := customer.CustomerContext{
		Repo: repo,
	}

	createTableIfExist(repo)

	r := setupRouter(customerContext)

	r.Run(":2019")
}

func setupRouter(cc customer.CustomerContext) *gin.Engine {
	r := gin.Default()
	r.Use(authMiddleware)
	r.POST("/customers", cc.CreateCustomerHandler)
	r.GET("/customers/:id", cc.GetCustomerByIdHandler)
	r.GET("/customers", cc.GetListOfCustomerHandler)
	r.PUT("/customers/:id", cc.UpdateCustomerIdHandler)
	r.DELETE("/customers/:id", cc.DeleteCustomerByIdHandler)
	return r
}

func createTableIfExist(repo interface {
	CreateTableIfNotExist() error
}) {
	err := repo.CreateTableIfNotExist()
	if err != nil {
		s := fmt.Sprintf("Unable to create table: %v", err)
		log.Fatal(s)
	}
}

func authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		c.Abort()
		return
	}
	c.Next()
}
