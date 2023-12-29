package calculations
import (
	"space-travel/structs"
	"time"
	"fmt"
	"strings"
)
// Struct to represent a route
type Route struct {
	From        string
	Destination string
}

var planetRoutes = map[string][]string{
	"Mercury": {"Venus"},
	"Venus":   {"Earth", "Mercury"},
	"Earth":   {"Jupiter", "Uranus"},
	"Mars":    {"Venus"},
	"Jupiter": {"Mars", "Venus"},
	"Saturn":  {"Earth", "Neptune"},
	"Uranus":  {"Saturn", "Neptune"},
	"Neptune": {"Uranus", "Mercury"},
}
// Recursive function to calculate all possible routes between two planets
func CalculateShortestRoute(from string, destination string) ([]Route) {
	var allRoutes [][]string
	calculateRoutes(from, destination, &allRoutes, []string{})
	return shortestRoute(allRoutes)
}

func calculateRoutes(currentPlanet string, destination string, allRoutes *[][]string, path []string) {
	pathCopy := make([]string, len(path))
	copy(pathCopy, path)
	
	for _, planet := range pathCopy {
		if planet == currentPlanet {
			return
		}
	}
	pathCopy = append(pathCopy, currentPlanet)
	if currentPlanet == destination {
		*allRoutes = append(*allRoutes, pathCopy)
		return
	}

	possibleDestinations, ok := planetRoutes[currentPlanet]
	if !ok {
		return
	}

	for _, nextPlanet := range possibleDestinations {
		calculateRoutes(nextPlanet, destination, allRoutes, pathCopy)
	}
}

func shortestRoute(allRoutes [][]string) []Route {
	if len(allRoutes) == 0 {
		return nil
	}

	shortest := allRoutes[0]

	for _, route := range allRoutes {
		if len(route) < len(shortest) {
			shortest = route
		}
	}

	return convertToRoute(shortest)
}

func convertToRoute(planetNames []string) []Route {
	var route []Route

	for i := 0; i < len(planetNames)-1; i++ {
		route = append(route, Route{
			From:        planetNames[i],
			Destination: planetNames[i+1],
		})
	}

	return route
}

func MakeCorrectRoutes(providers [][]structs.SimplifiedProvider) []structs.PossibleRoute {
	var possiblePermutations = Loop(providers)
	nrOfJumps := len(providers)
	routes := []structs.PossibleRoute{}

	for _, permutation := range possiblePermutations {
		route := structs.PossibleRoute{}
		var totalPrice float64
		var totalDuration time.Duration
		var firstTakeoff time.Time
		var lastLanding time.Time

		for i := 0; i < nrOfJumps; i++ {
			provider := providers[i][permutation[i]]
			route.Providers = append(route.Providers, provider)
			totalPrice += provider.Price

			if i == 0 || provider.FlightStart.Before(firstTakeoff) {
				firstTakeoff = provider.FlightStart
			}

			landingTime := provider.FlightEnd
			if i == nrOfJumps-1 || landingTime.After(lastLanding) {
				lastLanding = landingTime
			}
		}

		// Calculate total duration
		totalDuration = lastLanding.Sub(firstTakeoff).Round(time.Minute)

		// Format total duration as days, hours, and minutes
		days := totalDuration / (24 * time.Hour)
		totalDuration = totalDuration % (24 * time.Hour)
		hours := totalDuration / time.Hour
		totalDuration = totalDuration % time.Hour
		minutes := totalDuration / time.Minute

		durationString := ""
		if days > 0 {
			durationString += fmt.Sprintf("%d days, ", days)
		}
		if hours > 0 {
			durationString += fmt.Sprintf("%d hours, ", hours)
		}
		durationString += fmt.Sprintf("%d minutes", minutes)

		route.TotalPrice = fmt.Sprintf("%.2f", totalPrice)
		route.TotalDuration = durationString
		routes = append(routes, route)
	}

	return routes
}

func Loop(providers [][]structs.SimplifiedProvider) [][]int {
	var nrOfJumps = len(providers)
	var lengthsOfProviders = make([]int, nrOfJumps)
	for i := 0; i < nrOfJumps; i++ {
		lengthsOfProviders[i] = len(providers[i])
	}

	var permutationsArray = make([][]int, 0)
	var counterArray = make([]int, nrOfJumps)

	var generatePermutations func(int)
	generatePermutations = func(pos int) {
		if pos == nrOfJumps {
			// Create a copy of the counterArray and append it to permutationsArray
			tmp := make([]int, nrOfJumps)
			copy(tmp, counterArray)
			permutationsArray = append(permutationsArray, tmp)
			return
		}

		for i := 0; i < lengthsOfProviders[pos]; i++ {
			// Check if the flights match the time condition
			if pos == 0 || timesMatch(providers[pos-1][counterArray[pos-1]], providers[pos][i]) {
				counterArray[pos] = i
				generatePermutations(pos + 1)
			}
		}
	}

	generatePermutations(0)
	return permutationsArray
}

func timesMatch(providerA structs.SimplifiedProvider, providerB structs.SimplifiedProvider) bool{
	return providerA.FlightEnd.Before(providerB.FlightStart)
}

func ArrayToString(arr []string) string {
	return "'" + strings.Join(arr, ",") + "'"
}