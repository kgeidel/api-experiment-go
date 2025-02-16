package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Product defines a structure for an item in product catalog
type Product struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	AvailableFlag bool    `json:"available_flag"`
}

func get_db_conn_str() string {
	// fetch parameters from env
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")

	conn_str := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DB_USER,
		DB_PASSWORD,
		DB_HOST,
		DB_PORT,
		DB_NAME,
	)
	return conn_str
}

func get_db_connection() (*sql.DB, error) {
	conn_str := get_db_conn_str()
	// open database
	db, err := sql.Open("postgres", conn_str)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetProductHandler is used to get data inside the products defined on our product catalog
func GetProductHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Query the DB and return all results
		db, err := get_db_connection()
		if err != nil {
			fmt.Println(err)
		}
		products := []Product{}
		rows, err := db.Query("SELECT * FROM product;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var product Product
			err := rows.Scan(
				&product.ID,
				&product.Name,
				&product.Description,
				&product.Price,
				&product.AvailableFlag,
			)
			if err != nil {
				log.Fatal(err)
			}
			products = append(products, product)
		}
		b, err := json.MarshalIndent(products, "", "  ")
		if err != nil {
			fmt.Println("Error:", err)
		}
		rw.Header().Add("content-type", "application/json")
		rw.WriteHeader(http.StatusFound)
		rw.Write(b)
	}
}

// CreateProductHandler is used to create a new product and add to our product store.
func CreateProductHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Read incoming JSON from request body
		data, err := io.ReadAll(r.Body)
		// If no body is associated return with StatusBadRequest
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		// Check if data is proper JSON (data validation)
		var product Product
		err = json.Unmarshal(data, &product)
		if err != nil {
			rw.WriteHeader(http.StatusExpectationFailed)
			rw.Write([]byte("Invalid Data Format"))
			return
		}
		// Insert record into product table
		db, err := get_db_connection()
		if err != nil {
			fmt.Println(err)
		}
		qstr := `
			INSERT INTO product (name, description, price, available_flag)
			VALUES ($1, $2, $3, $4) returning "id";
		`
		statement, err := db.Prepare(qstr)
		if err != nil {
			panic(err)
		}
		defer statement.Close()
		var new_pk int
		err = statement.QueryRow(product.Name, product.Description, product.Price, product.AvailableFlag).Scan(&new_pk)
		if err != nil {
			panic(err)
		}
		success_message := fmt.Sprintf("Added new product [id: %d]", new_pk)
		// return after writing Body
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(success_message))
	}
}

// Create new Router

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Create new Router
	router := mux.NewRouter()

	// route properly to respective handlers
	router.Handle("/", GetProductHandler()).Methods("GET")
	router.Handle("/", CreateProductHandler()).Methods("POST")

	// Create new server and assign the router
	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	fmt.Println("Server listening on Port 8000...")
	// Start Server on defined port/host.
	server.ListenAndServe()
}
