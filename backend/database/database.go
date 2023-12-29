package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"space-travel/calculations"
	"space-travel/structs"
	"strconv"
	"time"
)

var (
	ErrNoPricelist = errors.New("No pricelist")
	ErrNoProviders = errors.New("No providers")
	maxPricelists  = 15
)

// AddBooking inserts a new booking into the database
func AddBooking(db *sql.DB, booking structs.Booking) error {
	insertBookingSQL := `
		INSERT INTO Bookings (
			CompanyNames,
			StartTime,
			FirstName,
			LastName,
			TotalPrice,
			TotalDuration,
			PricelistID,
			FromCity,
			DestinationCity
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(
		insertBookingSQL,
		calculations.ArrayToString(booking.CompanyNames),
		booking.StartTime,
		booking.FirstName,
		booking.LastName,
		booking.TotalPrice,
		booking.TotalDuration,
		booking.PricelistID,
		booking.Routes.From,
		booking.Routes.Destination,
	)
	if err != nil {
		return fmt.Errorf("failed to insert booking: %v", err)
	}

	return nil
}

func InsertPricelistData(db *sql.DB, pricelist structs.Pricelist) error {
	if err := checkMaxPriceLists(db, pricelist); err != nil {
		return err
	}
	exists, err := pricelistExists(db, pricelist.ID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	for _, leg := range pricelist.Legs {
		fromLocation := leg.RouteInfo.From
		toLocation := leg.RouteInfo.To
		if err := insertLocation(db, fromLocation, leg.ID); err != nil {
			return err
		}
		if err := insertLocation(db, toLocation, leg.ID); err != nil {
			return err
		}
	}

	// Insert Company, RouteInfo, and Route data
	for _, leg := range pricelist.Legs {
		if err := insertCompany(db, leg.Providers[0].Company, pricelist.ID); err != nil {
			return err
		}
		if err := insertRouteInfo(db, leg.RouteInfo, leg.ID); err != nil {
			return err
		}
		if err := insertRoute(db, leg, pricelist.ID); err != nil {
			return err
		}
	}

	// Insert Provider data
	for _, leg := range pricelist.Legs {
		for _, provider := range leg.Providers {
			if err := insertProvider(db, provider, leg.ID); err != nil {
				return err
			}
		}
	}

	// Insert Pricelist data
	if err := insertPricelist(db, pricelist); err != nil {
		return err
	}

	return nil
}

func pricelistExists(db *sql.DB, pricelistID string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pricelists WHERE ID = ?", pricelistID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func checkMaxPriceLists(db *sql.DB, pricelist structs.Pricelist) error {
	// Check if the Pricelist already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pricelists").Scan(&count)
	if err != nil {
		return err
	}

	if count >= maxPricelists {
		if err := deleteLoop(db, count); err != nil {
			return err
		}
	}

	return nil
}

func deleteLoop(db *sql.DB, count int) error {
	if count >= maxPricelists {
		if err := deleteOldestPricelistAndRelatedData(db); err != nil {
			return err
		}
		deleteLoop(db, count-1)
	}
	return nil
}

func deleteOldestPricelistAndRelatedData(db *sql.DB) error {
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Step 1: Get the ID of the oldest Pricelist
	var pricelistID string
	err = tx.QueryRow("SELECT ID FROM Pricelists ORDER BY ValidUntil ASC LIMIT 1").Scan(&pricelistID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete from Companies
	_, err = tx.Exec("DELETE FROM Companies WHERE PriceListID = ?", pricelistID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM Providers WHERE LegID IN (SELECT ID FROM Legs WHERE PriceListID = ?)", pricelistID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM Locations WHERE LegID IN (SELECT ID FROM Legs WHERE PriceListID = ?)", pricelistID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM RouteInfos WHERE LegID IN (SELECT ID FROM Legs WHERE PriceListID = ?)", pricelistID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM Legs WHERE PriceListID = ?", pricelistID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM Pricelists WHERE ID = ?", pricelistID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE * FROM Bookings WHERE PricelistID = ?", pricelistID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Function to insert Location data into the database
func insertLocation(db *sql.DB, location structs.Location, legID string) error {
	// Check if the Location already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM locations WHERE id = ?", location.ID).Scan(&count)
	if err != nil {
		return err
	}

	// Insert Location if it doesn't exist
	if count == 0 {
		_, err := db.Exec("INSERT INTO locations (id, name, legID) VALUES (?, ?, ?)", location.ID, location.Name, legID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Function to insert Company data into the database
func insertCompany(db *sql.DB, company structs.Company, pricelistID string) error {
	// Check if the Company already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM companies WHERE id = ?", company.ID).Scan(&count)
	if err != nil {
		return err
	}

	// Insert Company if it doesn't exist
	if count == 0 {
		_, err := db.Exec("INSERT INTO companies (id, name, pricelistID) VALUES (?, ?, ?)", company.ID, company.Name, pricelistID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Function to insert RouteInfo data into the database
func insertRouteInfo(db *sql.DB, routeInfo structs.RouteInfo, legID string) error {
	_, err := db.Exec("INSERT INTO routeInfos (id, FromID, ToID, distance, LegID) VALUES (?, ?, ?, ?, ?)",
		routeInfo.ID, routeInfo.From.ID, routeInfo.To.ID, routeInfo.Distance, legID)
	if err != nil {
		return err
	}

	return nil
}

// Function to insert Route data into the database
func insertRoute(db *sql.DB, route structs.Leg, priceListID string) error {
	_, err := db.Exec("INSERT INTO legs (id, routeInfoId, PriceListID) VALUES (?, ?, ?)", route.ID, route.RouteInfo.ID, priceListID)
	if err != nil {
		return err
	}

	return nil
}

// Function to insert Provider data into the database
func insertProvider(db *sql.DB, provider structs.Provider, legID string) error {
	_, err := db.Exec("INSERT INTO providers (id, companyID, price, flightStart, flightEnd, legID) VALUES (?, ?, ?, ?, ?, ?)",
		provider.ID, provider.Company.ID, provider.Price, provider.FlightStart, provider.FlightEnd, legID)
	if err != nil {
		return err
	}

	return nil
}

// Function to insert Pricelist data into the database
func insertPricelist(db *sql.DB, pricelist structs.Pricelist) error {
	count := 0
	err := db.QueryRow("SELECT COUNT(*) FROM pricelists WHERE id = ?", pricelist.ID).Scan(&count)
	if err != nil {
		return err
	}

	// Insert pricelist if it doesn't exist
	if count == 0 {
		_, err = db.Exec("INSERT INTO pricelists (id, validUntil) VALUES (?, ?)", pricelist.ID, pricelist.ValidUntil)
		if err != nil {
			return err
		}
	}

	return nil
}

// Function to get simplified data from the latest Pricelist for any given route
func GetAllPossibleRoutes(db *sql.DB, from string, destination string) (structs.GetResponse, error) {
	latestPricelistID, err := getLatestPricelistID(db)
	if err != nil {
		return structs.GetResponse{}, err
	}

	cachedRoutes, err := getCachedRoutes(db, latestPricelistID, from, destination)
	if err == nil {
		var cachedData structs.GetResponse
		if err := json.Unmarshal([]byte(cachedRoutes), &cachedData); err != nil {
			return structs.GetResponse{}, err
		}
		return cachedData, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return structs.GetResponse{}, err
	}

	finalRoute := calculations.CalculateShortestRoute(from, destination)
	providers, totalDistance, err := providersAndTotalDistance(db, finalRoute, latestPricelistID)
	if len(providers) < len(finalRoute) {
		return structs.GetResponse{}, ErrNoProviders
	}

	possibleRoutes := calculations.MakeCorrectRoutes(providers)
	// Get the validUntil from the database
	var validUntil string
	err = db.QueryRow("SELECT ValidUntil FROM Pricelists WHERE ID = ?", latestPricelistID).Scan(&validUntil)
	if err != nil {
		return structs.GetResponse{}, err
	}
	totalDistanceStr := strconv.Itoa(totalDistance)
	// Construct the final response
	response := structs.GetResponse{
		TotalDistance:  totalDistanceStr,
		ValidUntil:     validUntil,
		PricelistID:    latestPricelistID,
		PossibleRoutes: possibleRoutes,
	}
	err = cacheAndUpdateResponse(db, latestPricelistID, from, destination, response)
	if err != nil {
		return structs.GetResponse{}, err
	}

	return response, nil
}

func providersAndTotalDistance(db *sql.DB, finalRoute []calculations.Route, latestPricelistID string) ([][]structs.SimplifiedProvider, int, error) {
	var providers [][]structs.SimplifiedProvider
	var totalDistance int
	for _, route := range finalRoute {
		query := `
			SELECT Providers.Price, Providers.FlightStart, Providers.FlightEnd,
				RouteInfos.Distance, Companies.ID, Companies.Name AS CompanyName
			FROM Legs
			JOIN RouteInfos ON Legs.RouteInfoID = RouteInfos.ID
			JOIN Locations LocationsFrom ON RouteInfos.FromID = LocationsFrom.ID
			JOIN Locations LocationsTo ON RouteInfos.ToID = LocationsTo.ID
			JOIN Providers ON Legs.ID = Providers.LegID
			JOIN Companies ON Providers.CompanyID = Companies.ID
			WHERE Legs.PriceListID = ? AND LocationsFrom.Name = ? AND LocationsTo.Name = ?
			ORDER BY Providers.FlightStart
		`

		rows, err := db.Query(query, latestPricelistID, route.From, route.Destination)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return [][]structs.SimplifiedProvider{}, 0, err
		}
		defer rows.Close()

		var simplifiedProviders []structs.SimplifiedProvider
		var routeDistance = 0
		for rows.Next() {
			var companyName, companyID string
			var price float64
			var flightStart, flightEnd time.Time
			var distance int

			err := rows.Scan(&price, &flightStart, &flightEnd, &distance, &companyID, &companyName)
			if err != nil {
				return [][]structs.SimplifiedProvider{}, 0, err
			}
			if routeDistance == 0 {
				routeDistance = distance
			}

			// Add Provider to the slice
			simplifiedProviders = append(simplifiedProviders, structs.SimplifiedProvider{
				CompanyName: companyName,
				CompanyID:   companyID,
				Price:       price,
				FlightStart: flightStart,
				FlightEnd:   flightEnd,
			})
		}

		totalDistance += routeDistance
		providers = append(providers, simplifiedProviders)
	}
	return providers, totalDistance, nil
}

func cacheAndUpdateResponse(db *sql.DB, latestPricelistID, from, destination string, possibleRoutes structs.GetResponse) error {
	jsonRoutes, err := json.Marshal(possibleRoutes)
	if err != nil {
		log.Println(err)
		return err
	}

	err = updateCachedRoutes(db, latestPricelistID, from, destination, string(jsonRoutes))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func getCachedRoutes(db *sql.DB, latestPricelistID string, from, destination string) (string, error) {
	var cachedRoutes string
	err := db.QueryRow("SELECT Routes FROM CachedRoutes WHERE PricelistID = ? AND FromLocation = ? AND ToLocation = ?", latestPricelistID, from, destination).Scan(&cachedRoutes)
	if err != nil {
		return "", err
	}
	return cachedRoutes, nil
}

func updateCachedRoutes(db *sql.DB, latestPricelistID string, from, destination string, routes string) error {
	_, err := db.Exec(`
		INSERT INTO CachedRoutes (ID, PricelistID, FromLocation, ToLocation, Routes)
		VALUES (?, ?, ?, ?, ?)
	`, uuid.New().String(), latestPricelistID, from, destination, routes)
	return err
}

func getLatestPricelistID(db *sql.DB) (string, error) {
	var latestPricelistID string
	err := db.QueryRow("SELECT ID FROM Pricelists ORDER BY ValidUntil DESC LIMIT 1").Scan(&latestPricelistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoPricelist
		}
		return "", err
	}
	return latestPricelistID, nil
}

func CleanCache(db *sql.DB, pricelistID string) error {
	_, err := db.Exec("DELETE FROM CachedRoutes WHERE PricelistID != ?", pricelistID)
	return err
}
