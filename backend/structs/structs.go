package structs

import (
	"time"
)

type Location struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RouteInfo struct {
	ID       string `json:"id"`
	From     Location
	To       Location
	Distance int `json:"distance"`
}

type Provider struct {
	ID          string `json:"id"`
	Company     Company
	Price       float64   `json:"price"`
	FlightStart time.Time `json:"flightStart"`
	FlightEnd   time.Time `json:"flightEnd"`
}

type Leg struct {
	ID        string     `json:"id"`
	RouteInfo RouteInfo  `json:"routeInfo"`
	Providers []Provider `json:"providers"`
}

type Pricelist struct {
	ID         string    `json:"id"`
	ValidUntil time.Time `json:"validUntil"`
	Legs       []Leg     `json:"legs"`
}

type SimplifiedProvider struct {
	CompanyName string    `json:"companyName"`
	CompanyID   string    `json:"companyID"`
	Price       float64   `json:"price"`
	FlightStart time.Time `json:"flightStart"`
	FlightEnd   time.Time `json:"flightEnd"`
}

type PossibleRoute struct {
	TotalPrice    string               `json:"totalPrice"`
	TotalDuration string               `json:"totalDuration"`
	Providers     []SimplifiedProvider `json:"providers"`
}

type GetResponse struct {
	TotalDistance string `json:"totalDistance"`
	ValidUntil    string `json:"validUntil"`
	PricelistID   string `json:"pricelistID"`
	PossibleRoutes []PossibleRoute `json:"possibleRoutes"`
}

type Booking struct {
    CompanyNames []string    // Array of company names
    StartTime   string       // Start time of the flight
    FirstName   string       // First name of the passenger
    LastName   string       // Last name of the passenger
    TotalPrice  float64      // Total price of the booking
    TotalDuration string         // Total duration of the flight in minutes
    PricelistID  string        // ID of the pricelist for the booking
    Routes      Routes       // Route details
    ValidUntil  string       // Valid until date for the booking
}

type Routes struct {
    From    string  // Departure city
    Destination string // Destination city
}
