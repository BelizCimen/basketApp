package main

import (
	"modules/services"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/products", services.GetProducts).Methods("GET")
	r.HandleFunc("/customers", services.GetCustomers).Methods("GET")
	r.HandleFunc("/addProduct", services.AddProduct).Methods("POST")

	r.HandleFunc("/addCustomer", services.AddCustomer).Methods("POST")

	r.HandleFunc("/addToCart/{id}", services.AddToCart).Methods("POST")
	r.HandleFunc("/deleteOneItemFromCart/{id}", services.DeleteFromCart).Methods("DELETE")
	r.HandleFunc("/cart", services.GetCart).Methods("GET")

	r.HandleFunc("/order", services.Order).Methods("POST")
	r.HandleFunc("/order/{id}", services.GetOrderById).Methods("GET")

	handler := cors.AllowAll().Handler(r)
	fmt.Printf("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

}
