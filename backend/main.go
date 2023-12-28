package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"slices"
	"space-travel/database"
	"space-travel/structs"
	"time"
)

const (
	tables          = "./database/sql/tables.sql"
	dbPath          = "./database/pricelists.db"
	travelPricesURL = "https://cosmos-odyssey.azurewebsites.net/api/v1.0/TravelPrices"
)

var validPlanets = []string{"Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"}

func checkLastPricelistValidity(db *sql.DB) (bool, time.Duration) {
	// Query to get the last entered pricelist
	query := "SELECT ID, ValidUntil FROM Pricelists ORDER BY ValidUntil DESC LIMIT 1"

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return false, 0
	}
	defer rows.Close()

	if rows.Next() {
		var pricelistID string
		var validUntil time.Time

		if err := rows.Scan(&pricelistID, &validUntil); err != nil {
			log.Println(err)
			return false, 0
		}

		if validUntil.Before(time.Now()) {
			return false, 0
		} else {
			durationUntilInvalid := validUntil.Sub(time.Now())
			return true, durationUntilInvalid
		}
	}
	return false, 0
}

// Fetch travel prices and store in the database
func fetchAndStoreTravelPrices(db *sql.DB) (error, time.Duration) {
	valid, duration := checkLastPricelistValidity(db)
	if valid {
		return nil, duration
	}
	resp, err := http.Get(travelPricesURL)
	if err != nil {
		return err, 0
	}
	defer resp.Body.Close()

	var list structs.Pricelist
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&list); err != nil {
		return err, 0
	}

	if err := database.InsertPricelistData(db, list); err != nil {
		return err, 0
	}
	if err := database.CleanCache(db, list.ID); err != nil {
		return err, 0
	}
	valid, duration = checkLastPricelistValidity(db)
	return nil, duration
}

// Handle "/api/get/:from/:destination" endpoint
func handleGetAPI(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	from := vars["from"]
	destination := vars["destination"]
	if !checkURLParams(from, destination) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data, err := database.GetAllPossibleRoutes(db, from, destination)
	if err != nil {
		if err == database.ErrNoProviders {
			http.Error(w, "No providers found", http.StatusNotFound)
		} else {
			log.Println("error here: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Println("error: ", err)
		}
	}
}

func handlePostAPI(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var booking structs.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	err = database.AddBooking(db, booking)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func checkURLParams(from string, destination string) bool {
	if slices.Contains(validPlanets, from) == false || slices.Contains(validPlanets, destination) == false {
		return false
	}
	return true
}

func main() {
	var db *sql.DB
	// Check if the database file exists
	if _, err := os.ReadFile(dbPath); err != nil {
		log.Printf("Database file does not exist. Creating tables...")

		createTableSQL, err := os.ReadFile(tables)
		if err != nil {
			log.Fatal(err)
		}

		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec(string(createTableSQL))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Tables created successfully.")
	} else {
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
	}
	go func() {
		for {
			err, duration := fetchAndStoreTravelPrices(db)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Minute)
			} else {
				time.Sleep(duration)
			}
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/api/get/{from}/{destination}", func(w http.ResponseWriter, r *http.Request) {
		handleGetAPI(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/api/post", func(w http.ResponseWriter, r *http.Request) {
		handlePostAPI(w, r, db)
	}).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8085"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
