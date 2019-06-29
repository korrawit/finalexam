package repository

import (
	"strconv"

	"github.com/korrawit/finalexam/database"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type Repository struct {
	DB database.Interface
}

type CustomerRepository interface {
	CreateTableIfNotExist() error
	CreateNewCustomer(c *Customer) error
	GetCustomers() ([]Customer, error)
	GetCustomerById(id string) (*Customer, error)
	UpdateCustomer(id string, c *Customer) error
	DeleteCustomerById(id string) error
}

func (r Repository) CreateTableIfNotExist() error {
	db, err := r.DB.Connect()
	defer db.Close()

	if err != nil {
		return err
	}

	ddl := `
		CREATE TABLE IF NOT EXISTS customer (
			id SERIAL PRIMARY KEY,
			name TEXT,  
			email TEXT,  
			status 	TEXT
	  )
	`

	_, err = db.Exec(ddl)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetCustomers() ([]Customer, error) {
	db, err := r.DB.Connect()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customer ORDER BY id")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	customers := []Customer{}
	for rows.Next() {
		c := Customer{}
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Status)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r Repository) CreateNewCustomer(c *Customer) error {
	db, err := r.DB.Connect()
	defer db.Close()

	if err != nil {
		return err
	}

	query := `
			INSERT INTO customer (name, email ,status) VALUES ($1,$2,$3) RETURNING id
		`

	row := db.QueryRow(query, c.Name, c.Email, c.Status)
	err = row.Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetCustomerById(id string) (*Customer, error) {
	db, err := r.DB.Connect()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customer WHERE id = $1")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	var c Customer
	err = row.Scan(&c.ID, &c.Name, &c.Email, &c.Status)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r Repository) UpdateCustomer(id string, c *Customer) error {
	db, err := r.DB.Connect()
	defer db.Close()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE customer SET name=$2, email=$3, status=$4 WHERE id=$1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id, c.Name, c.Email, c.Status)
	if err != nil {
		return err
	}

	c.ID, err = strconv.Atoi(id)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) DeleteCustomerById(id string) error {
	db, err := r.DB.Connect()
	defer db.Close()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("DELETE FROM customer WHERE id =$1")
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
