package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

func main() {
	// connect to database
	db, err := connectDB()
	if err != nil {
		fmt.Println("Database connection error: ", err)
		return
	}

	// migrate tables
	err = migrateDatabase(db)
	if err != nil {
		fmt.Println("Database migration error: ", err)
		return
	}

	fmt.Println("Database migrated successfully!")

	// load and seed Customers table
	customers, err := loadCustomers()
	if err != nil {
		fmt.Println("Error loading customers: ", err)
		return
	}

	if err := seedCustomers(db, customers); err != nil {
		fmt.Println("Error injecting customers: ", err)
		return
	}

	fmt.Println("Customers injected into table successfully")

	// load and seed Equipment table
	equipment, err := loadEquipment()
	if err != nil {
		fmt.Println("Error loading equipment: ", err)
		return
	}

	if err := seedEquipment(db, equipment); err != nil {
		fmt.Println("Error injecting equipment: ", err)
		return
	}

	fmt.Println("Equipment injected into table successfully")

	// load and seed LaborRates table
	laborRates, err := loadLaborRates()
	if err != nil {
		fmt.Println("Error loading labor rates: ", err)
		return
	}

	if err := seedLaborRates(db, laborRates); err != nil {
		fmt.Println("Error injecting labor rates: ", err)
		return
	}

	fmt.Println("Labor rates injected into table successfully")

	// test API endpoint to check backend is running
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Backend is running!")
	})

	// API endpoint to get customer data using customerID
	http.HandleFunc("/api/customers/", func(w http.ResponseWriter, r *http.Request) {
		customerID := strings.TrimPrefix(r.URL.Path, "/api/customers/")

		var customer Customer

		if err := db.First(&customer, "id = ?", customerID).Error; err != nil {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customer)
	})

	// API endpoint to get equipment data
	http.HandleFunc("/api/equipment", func(w http.ResponseWriter, r *http.Request) {
		var equipment []Equipment

		if err := db.Find(&equipment).Error; err != nil {
			http.Error(w, "Failed to retrieve equipment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(equipment); err != nil {
			http.Error(w, "Failed to encode equipment", http.StatusInternalServerError)
		}
	})

	// API endpoint to get labor rates data
	http.HandleFunc("/api/laborRates", func(w http.ResponseWriter, r *http.Request) {
		var laborRates []LaborRate

		if err := db.Find(&laborRates).Error; err != nil {
			http.Error(w, "Failed to retrieve labor rates", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(laborRates); err != nil {
			http.Error(w, "Failed to encode labor rates", http.StatusInternalServerError)
		}
	})

	// CORS to connect backend to frontend
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
	})

	handler := c.Handler(http.DefaultServeMux)

	fmt.Println("Starting server on port 8080")

	// listen and serve on port 8080
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println("Server error: ", err)
	}
}
