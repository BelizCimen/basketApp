package services

import (
	"encoding/json"
	"log"
	"modules/db"
	"modules/models"
	"net/http"
)

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT * FROM customer")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Surname)

		if err != nil {
			log.Fatal(err)
		}

		customers = append(customers, &customer)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(customers)
}

func AddCustomer(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var customer models.Customer
		_ = json.NewDecoder(r.Body).Decode(&customer)

		result, err := db.DB.Exec("INSERT INTO customer(name,surname) VALUES($1, $2)", customer.Name, customer.Surname)

		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}

		w.Write([]byte("Customer added successfully"))

	}
}
