package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"modules/db"
	"modules/models"
	"net/http"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT * FROM products")

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found")
			return
		}
		log.Fatal(err)
	}

	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		prd := &models.Product{}
		err := rows.Scan(&prd.ID, &prd.Name, &prd.Price, &prd.Tax)

		if err != nil {
			log.Fatal(err)
		}

		products = append(products, prd)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(products)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var prdct models.Product
		_ = json.NewDecoder(r.Body).Decode(&prdct)

		result, err := db.DB.Exec("INSERT INTO products(name, price, tax) VALUES($1, $2, $3)", prdct.Name, prdct.Price, prdct.Tax)
		if err != nil {
			log.Fatal(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			log.Fatal(count)
		}

		w.Write([]byte("Product added successfully"))

	}
}
