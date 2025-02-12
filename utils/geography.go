// utils/geography.go

package utils

// stateCityMap maps US states to their cities from our dataset
var stateCityMap = map[string][]string{
	"washington":     {"seattle"},
	"california":     {"los angeles", "anaheim", "oakland", "del mar", "san francisco", "pico rivera", "inglewood"},
	"missouri":       {"kansas city", "st louis"},
	"texas":          {"dallas", "houston"},
	"georgia":        {"atlanta"},
	"massachusetts":  {"boston"},
	"new york":       {"brooklyn", "newark", "uniondale"},
	"illinois":       {"chicago", "berwyn", "rosemont"},
	"pennsylvania":   {"philadelphia", "pittsburgh", "hershey"},
	"michigan":       {"grand rapids", "detroit"},
	"indiana":        {"indianapolis"},
	"ohio":           {"cleveland", "cincinnati"},
	"nebraska":       {"omaha"},
	"north carolina": {"charlotte", "columbia"},
	"louisiana":      {"new orleans"},
	"wisconsin":      {"madison"},
}

// GetStateForCity returns the state that a city belongs to, or empty string if not found
func GetStateForCity(city string) string {
	for state, cities := range stateCityMap {
		for _, c := range cities {
			if c == city {
				return state
			}
		}
	}
	return ""
}

// GetCitiesInState returns all cities in a state
func GetCitiesInState(state string) []string {
	cities, exists := stateCityMap[state]
	if !exists {
		return []string{}
	}
	return cities
}
