package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

// Struktur untuk data customer
type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Middleware untuk Basic Auth
func basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username != "lumoshive" || password != "academy" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=postgres host=localhost")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var lastInsertId int
	query := `INSERT INTO customers (name, email) VALUES ($1, $2) RETURNING id`
	err = db.QueryRow(query, customer.Name, customer.Email).Scan(&lastInsertId)
	if err != nil {
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
		return
	}

	customer.ID = lastInsertId

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func listCustomersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=postgres host=localhost")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, email FROM customers")
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Email); err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func getCustomerByIDHandler(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	id, err := strconv.Atoi(customerID)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=postgres host=localhost")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var customer Customer
	query := `SELECT id, name, email FROM customers WHERE id = $1`
	err = db.QueryRow(query, id).Scan(&customer.ID, &customer.Name, &customer.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching data", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func updateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	id, err := strconv.Atoi(customerID)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=postgres host=localhost")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `UPDATE customers SET name = $1, email = $2 WHERE id = $3`
	_, err = db.Exec(query, customer.Name, customer.Email, id)
	if err != nil {
		http.Error(w, "Error updating data", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": fmt.Sprintf("Customer with ID %d updated successfully", id),
		"id":      id,
		"status":  "success",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func deleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	id, err := strconv.Atoi(customerID)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=postgres host=localhost")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	query := `DELETE FROM customers WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting data", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error getting rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message": fmt.Sprintf("Customer with ID %d deleted successfully", id),
		"id":      id,
		"status":  "success",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()

	// Tambahkan handler untuk route "/customers"
	mux.HandleFunc("POST /customers", createCustomerHandler)

	// Menambahkan route untuk endpoint list customers
	mux.Handle("GET /customers", basicAuthMiddleware(http.HandlerFunc(listCustomersHandler)))

	// Menambahkan route untuk endpoint get customer by ID
	mux.Handle("GET /customers/{id}", basicAuthMiddleware(http.HandlerFunc(getCustomerByIDHandler)))

	// Menambahkan route untuk endpoint update customer
	mux.Handle("PUT /customers/{id}", basicAuthMiddleware(http.HandlerFunc(updateCustomerHandler)))

	// Menambahkan route untuk endpoint delete customer
	mux.Handle("DELETE /customers/{id}", basicAuthMiddleware(http.HandlerFunc(deleteCustomerHandler)))

	fmt.Println("Server started at port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
