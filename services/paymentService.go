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
)

func totalPrice(quantity int, productPrice float64, tax float64) float64 {
	price := productPrice + (productPrice * tax)
	if quantity > 3 {
		total_price := (price * 3) + (price*(float64(quantity)-3) - totalDiscount(quantity, price))
		return total_price
	} else {
		totalPrice := price * (float64(quantity))
		return totalPrice
	}
}

func totalDiscount(quantity int, price float64) float64 {
	discount := price * 0.08

	if quantity > 3 {
		total_discount := discount * (float64(quantity) - 3)
		return total_discount
	} else {
		return 0.00
	}

}

func Order(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		var payment models.Payment
		_ = json.NewDecoder(r.Body).Decode(&payment)

		rows, err := db.DB.Query("SELECT * FROM payment WHERE customer_id = $1", payment.CustomerID)

		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No rows found")
				return
			}
			log.Fatal(err)
		}

		defer rows.Close()

		var custOrder []*models.Payment
		for rows.Next() {
			ord := &models.Payment{}
			err := rows.Scan(&ord.ID, &ord.CustomerID, &ord.TotalPrice)

			if err != nil {
				log.Fatal(err)
			}

			custOrder = append(custOrder, ord)
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		var count []float32
		for _, element := range custOrder {
			if element.TotalPrice > 500 { // given amount
				count = append(count, element.TotalPrice)
			}
		}

		if len(count) < 1 {
			rows, err := db.DB.Query("SELECT * FROM cart")

			if err != nil {
				if err == sql.ErrNoRows {
					fmt.Println("No rows found")
					return
				}
				log.Fatal(err)
			}

			defer rows.Close()

			var currCart []*models.Cart
			for rows.Next() {
				cart := &models.Cart{}
				err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)

				if err != nil {
					log.Fatal(err)
				}

				currCart = append(currCart, cart)
			}

			if err = rows.Err(); err != nil {
				log.Fatal(err)
			}

			for index, element := range currCart {

				var product float64
				var tax float64

				err := db.DB.QueryRow("SELECT price FROM products WHERE id = $1", element.ProductID).Scan(&product)
				err2 := db.DB.QueryRow("SELECT tax FROM products WHERE id = $1", element.ProductID).Scan(&tax)
				switch {
				case err == sql.ErrNoRows:
					fmt.Println("No rows found")
					return
				case err != nil:
					log.Fatal(err)
				}
				if err2 == sql.ErrNoRows {
					fmt.Println("No rows found")
					return
				}

				currCart[index].Price = totalPrice(currCart[index].Quantity, product, tax)
			}
			var total_price float64
			for _, element := range currCart {
				total_price += element.Price
			}
			result, err := db.DB.Exec("INSERT INTO payment(customer_id,total_price) VALUES($1,$2)", payment.CustomerID, total_price)

			if err != nil {
				log.Fatal(err)
			}
			count, err := result.RowsAffected()
			if err != nil {
				log.Fatal(count)
			}

			table, err := db.DB.Exec("DELETE FROM cart")

			if err != nil {
				log.Fatal(err)
			}
			effect, err := table.RowsAffected()
			if err != nil {
				log.Fatal(effect)
			}
			s := fmt.Sprintf("%f", total_price)
			w.Write([]byte("You just can use the Quantity promotion.Total price with discount:" + s))

		} else if len(count) >= 1 {
			lastOrder := len(count)
			if (lastOrder%4+1)%4 == 0 {
				rows, err := db.DB.Query("SELECT * FROM cart")

				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Println("No rows found")
						return
					}
					log.Fatal(err)
				}

				defer rows.Close()

				var currCart []*models.Cart
				for rows.Next() {
					cart := &models.Cart{}
					err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
					if err != nil {
						log.Fatal(err)
					}
					currCart = append(currCart, cart)
				}

				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}

				for index, element := range currCart {

					var product float64
					var tax float64
					var vatProduct = product + (product * tax)

					err := db.DB.QueryRow("SELECT price,tax FROM products WHERE id = $1", element.ProductID).Scan(&product, &tax)

					switch {
					case err == sql.ErrNoRows:
						fmt.Println("No rows found")
						return
					case err != nil:
						log.Fatal(err)
					}

					if tax == 1.00 {
						currCart[index].Price = vatProduct * float64(currCart[index].Quantity)
					} else if tax == 8.00 {
						currCart[index].Price = (vatProduct - (vatProduct * 0.1)) * float64(currCart[index].Quantity)
					} else if tax == 18.00 {
						currCart[index].Price = (vatProduct - (vatProduct * 0.15)) * float64(currCart[index].Quantity)
					}
				}
				var total_price float64
				for _, element := range currCart {

					total_price += element.Price

				}
				result, err := db.DB.Exec("INSERT INTO payment(customer_id,total_price) VALUES($1,$2)", payment.CustomerID, total_price)

				if err != nil {
					log.Fatal(err)
				}
				count, err := result.RowsAffected()
				if err != nil {
					log.Fatal(count)
				}

				table, err := db.DB.Exec("DELETE FROM cart")

				if err != nil {
					log.Fatal(err)
				}
				effect, err := table.RowsAffected()
				if err != nil {
					log.Fatal(effect)
				}
				s := fmt.Sprintf("%f", total_price)
				w.Write([]byte("You made your 4th order of higher than given amount, Discount applied. Total price with discount:" + s))

			} else {
				rows, err := db.DB.Query("SELECT * FROM cart WHERE total_discount > 0")

				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Println("No rows found")
						return
					}
					log.Fatal(err)
				}

				defer rows.Close() //

				var currCart []*models.Cart
				for rows.Next() {
					cart := &models.Cart{}
					err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
					if err != nil {
						log.Fatal(err)
					}
					currCart = append(currCart, cart)
				}

				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}

				var discounts []float64

				for _, element := range currCart {
					if element.Discount > 0 {
						discounts = append(discounts, element.Discount)
					}
				}
				if len(discounts) < 1 {
					var currCart []*models.Cart
					for rows.Next() {
						cart := &models.Cart{}
						err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
						if err != nil {
							log.Fatal(err)
						}
						currCart = append(currCart, cart)
					}

					if err = rows.Err(); err != nil {
						log.Fatal(err)
					}

					for index, element := range currCart {

						var product float64
						var tax float64
						var vatProduct = product + product*tax
						err := db.DB.QueryRow("SELECT price, tax FROM products WHERE id = $1", element.ProductID).Scan(&product, &tax)
						switch {
						case err == sql.ErrNoRows:
							fmt.Println("No rows found")
							return
						case err != nil:
							log.Fatal(err)
						}

						currCart[index].Price = vatProduct * float64(currCart[index].Quantity)
					}
					var total_price float64
					for _, element := range currCart {
						total_price += element.Price
					}
					discount := total_price * 0.1
					result, err := db.DB.Exec("INSERT INTO payment(customer_id,total_price) VALUES($1,$2)", payment.CustomerID, total_price-discount)

					if err != nil {
						log.Fatal(err)
					}
					count, err := result.RowsAffected()
					if err != nil {
						log.Fatal(count)
					}

					table, err := db.DB.Exec("DELETE FROM cart")

					if err != nil {
						log.Fatal(err)
					}
					effect, err := table.RowsAffected()
					if err != nil {
						log.Fatal(effect)
					}
					s := fmt.Sprintf("%f", total_price-discount)
					w.Write([]byte("You made your order higher than given amount, Thanks!  %10 discount is applied to total price of your cart.Total Price:" + s))

				} else if len(discounts) >= 1 {

					rows, err := db.DB.Query("SELECT * FROM cart")

					if err != nil {
						if err == sql.ErrNoRows {
							fmt.Println("No rows found")
							return
						}
						log.Fatal(err)
					}

					defer rows.Close()

					var currCart []*models.Cart
					for rows.Next() {
						cart := &models.Cart{}
						err := rows.Scan(&cart.ID, &cart.ProductID, &cart.Quantity, &cart.Price, &cart.Discount)
						if err != nil {
							log.Fatal(err)
						}
						currCart = append(currCart, cart)
					}

					if err = rows.Err(); err != nil {
						log.Fatal(err)
					}

					for index, element := range currCart {

						var product float64
						var tax float64
						var vatProduct = product + (product * tax)

						err := db.DB.QueryRow("SELECT price, tax FROM products WHERE id = $1", element.ProductID).Scan(&product, &tax)
						switch {
						case err == sql.ErrNoRows:
							fmt.Println("No rows found")
							return
						case err != nil:
							log.Fatal(err)
						}

						currCart[index].Price = vatProduct * float64(currCart[index].Quantity)
					}
					var total_price float64
					for _, element := range currCart {
						total_price += element.Price
					}
					discount := total_price * 0.1
					result, err := db.DB.Exec("INSERT INTO payment(customer_id,total_price) VALUES($1,$2)", payment.CustomerID, total_price-discount)

					if err != nil {
						log.Fatal(err)
					}
					count, err := result.RowsAffected()
					if err != nil {
						log.Fatal(count)
					}

					table, err := db.DB.Exec("DELETE FROM cart")

					if err != nil {
						log.Fatal(err)
					}
					effect, err := table.RowsAffected()
					if err != nil {
						log.Fatal(effect)
					}
				}
			}
		}
	}
}

func GetOrderById(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	row := db.DB.QueryRow("SELECT * FROM payment WHERE customer_id = $1", id)
	pymnt := &models.Payment{}
	err := row.Scan(&pymnt.ID, &pymnt.CustomerID, &pymnt.TotalPrice)

	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No rows found")
		return
	case err != nil:
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(pymnt)

}
