package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"modules/db"
	"modules/models"
	"net/http"
	"strconv"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var cart models.Cart
		_ = json.NewDecoder(r.Body).Decode(&cart)
		id := mux.Vars(r)["id"]

		var product string
		err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", id).Scan(&product)

		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}

		var tax float64
		err = db.DB.QueryRow("SELECT tax FROM products WHERE id = $1", id).Scan(&tax)

		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}

		if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) == sql.ErrNoRows && cart.Quantity > 3 {
			floatVar, _ := strconv.ParseFloat(product, 64)
			result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2, $3, $4)", id, cart.Quantity, totalPrice(cart.Quantity, floatVar, tax), totalDiscount(cart.Quantity, floatVar))

			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}

			w.Write([]byte("products were added with discount")) //

		} else if db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ProductID) == sql.ErrNoRows && cart.Quantity <= 3 { // Check if the product is already in the cart and the user wants to add less than 3 products

			floatVar, _ := strconv.ParseFloat(product, 64) // Convert the string to a float64

			result, err := db.DB.Exec("INSERT INTO cart(product_id, quantity, total_price, total_discount) VALUES($1, $2, $3, $4)", id, cart.Quantity, totalPrice(cart.Quantity, floatVar, tax), totalDiscount(cart.Quantity, floatVar)) // Execute the SQL Query to insert the product into the cart with the quantity and the total price and discount

			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}

		} else {

			var oldQuantity int
			err := db.DB.QueryRow("SELECT quantity FROM cart WHERE product_id = $1", id).Scan(&oldQuantity)

			switch {
			case err == sql.ErrNoRows:
				fmt.Println("No rows found")
				return
			case err != nil:
				log.Fatal(err)
			}

			floatVar, _ := strconv.ParseFloat(product, 64)
			newQuantity := cart.Quantity + oldQuantity
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, totalPrice(newQuantity, floatVar, tax), totalDiscount(newQuantity, floatVar), id)
			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}

			w.Write([]byte("cart is updated "))
		}
	}
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT * FROM cart")

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found")
			return
		}
		log.Fatal(err)
	}

	defer rows.Close()

	var cartItems []*models.Cart
	for rows.Next() {
		crt := &models.Cart{}
		err := rows.Scan(&crt.ID, &crt.ProductID, &crt.Quantity, &crt.Price, &crt.Discount)

		if err != nil {
			log.Fatal(err)
		}

		cartItems = append(cartItems, crt)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(cartItems)
}

func DeleteFromCart(w http.ResponseWriter, r *http.Request) {

	if r.Method == "DELETE" {
		w.Header().Set("Content-Type", "application/json")
		var cart models.Cart
		_ = json.NewDecoder(r.Body).Decode(&cart)
		id := mux.Vars(r)["id"]
		var product string
		var tax float64
		err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", id).Scan(&product)
		err2 := db.DB.QueryRow("SELECT tax FROM products WHERE id = $1", id).Scan(&product)

		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err != nil:
			log.Fatal(err)
		}

		switch {
		case err2 == sql.ErrNoRows:
			fmt.Println("No rows found")
			return
		case err2 != nil:
			log.Fatal(err)
		}

		row := db.DB.QueryRow("SELECT * FROM cart WHERE product_id = $1", id).Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)

		switch {
		case row == sql.ErrNoRows:
			fmt.Println("No rows found")
			w.Write([]byte("No Rows Found"))
			return
		case row != nil:
			log.Fatal(row)
		}

		floatVar, _ := strconv.ParseFloat(product, 64)
		newQuantity := cart.Quantity - 1

		if newQuantity <= 0 {
			result, err := db.DB.Exec("DELETE FROM cart WHERE product_id = $1", id)

			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}

			w.Write([]byte("cart empty"))
			return
		} else if newQuantity > 3 {
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, totalPrice(newQuantity, floatVar, tax), totalDiscount(newQuantity, floatVar), id)

			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}

			w.Write([]byte("One item deleted from the cart and discount updated.."))
			return
		} else {
			result, err := db.DB.Exec("UPDATE cart SET quantity = $1, total_price = $2, total_discount = $3 WHERE product_id = $4", newQuantity, totalPrice(newQuantity, floatVar, tax), totalDiscount(newQuantity, floatVar), id)

			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}

			w.Write([]byte("deleted from cart successfully"))
			return
		}
	}
}
