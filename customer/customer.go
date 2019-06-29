package customer

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/korrawit/finalexam/repository"
)

type CustomerContext struct {
	Repo interface {
		CreateNewCustomer(c *repository.Customer) error
		GetCustomers() ([]repository.Customer, error)
		GetCustomerById(id string) (*repository.Customer, error)
		UpdateCustomer(id string, c *repository.Customer) error
		DeleteCustomerById(id string) error
	}
}

func (cc CustomerContext) CreateCustomerHandler(c *gin.Context) {
	var cus repository.Customer
	err := c.BindJSON(&cus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cc.Repo.CreateNewCustomer(&cus)
	if err != nil {
		log.Print("Error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusCreated, cus)
}

func (cc CustomerContext) GetCustomerByIdHandler(c *gin.Context) {
	id := c.Param("id")
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid id format"})
		return
	}

	cus, err := cc.Repo.GetCustomerById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := fmt.Sprintf("Customer id %s not found", id)
			c.JSON(http.StatusNotFound, gin.H{"message": msg})
			return
		}
		log.Print("Error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, cus)
}

func (cc CustomerContext) GetListOfCustomerHandler(c *gin.Context) {
	cus, err := cc.Repo.GetCustomers()
	if err != nil {
		log.Print("Error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, cus)
}

func (cc CustomerContext) UpdateCustomerIdHandler(c *gin.Context) {
	id := c.Param("id")
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid id format"})
		return
	}

	var cus repository.Customer
	err := c.BindJSON(&cus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cc.Repo.UpdateCustomer(id, &cus)
	if err != nil {
		log.Print("Error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, cus)
}

func (cc CustomerContext) DeleteCustomerByIdHandler(c *gin.Context) {
	id := c.Param("id")
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid id format"})
		return
	}

	err := cc.Repo.DeleteCustomerById(id)
	if err != nil {
		log.Print("Error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}
