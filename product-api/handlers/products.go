// Package handlers
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"product-api/data"
	"strconv"

	"github.com/gorilla/mux"
)

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct{}

type Products struct {
	l *log.Logger
	v *data.Validation
}

// NewProducts returns a new products handler with the given logger
func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Error("Invalid Path, path should be /products/[id]")

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		//should never happen
		panic(err)
	}

	return id
}

//func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
//	lp := data.GetProducts()
//
//	err := lp.ToJSON(rw)
//	if err != nil {
//		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
//	}
//}
//
//func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
//	p.l.Println("Handle POST Product")
//
//	prod := r.Context().Value(KeyProduct{}).(data.Product)
//
//	p.l.Printf("Prod: %#v", prod)
//	data.AddProduct(&prod)
//}
//
//func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	id, err := strconv.Atoi(vars["id"])
//	if err != nil {
//		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
//	}
//
//	p.l.Println("Handle PUT Product", id)
//
//	prod := r.Context().Value(KeyProduct{}).(data.Product)
//
//	err = data.UpdateProducts(id, &prod)
//	if err == data.ErrProductNotFound {
//		http.Error(rw, "Product not found", http.StatusNotFound)
//		return
//	}
//
//	if err != nil {
//		http.Error(rw, "Product not found", http.StatusInternalServerError)
//		return
//	}
//}
//
//type KeyProduct struct{}
//
//func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
//		prod := data.Product{}
//
//		err := prod.FromJSON(r.Body)
//		if err != nil {
//			p.l.Println("[ERROR] deserializing product", err)
//			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
//			return
//		}
//
//		// validate the product
//		err = prod.Validate()
//		if err != nil {
//			p.l.Println("[ERROR] validating product", err)
//			http.Error(
//				rw,
//				fmt.Sprintf("Error validating product: %s", err),
//				http.StatusBadRequest,
//			)
//			return
//		}
//
//		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
//		r = r.WithContext(ctx)
//		next.ServeHTTP(rw, r)
//	})
//}
