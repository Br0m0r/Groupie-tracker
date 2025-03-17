package utils

// stateCityMap maps US states to their cities from our dataset
var StateCityMap = map[string][]string{
	"Washington, Usa":     {"Seattle, Usa"},
	"California, Usa":     {"Los Angeles, Usa", "Anaheim, Usa", "Oakland, Usa", "Del Mar, Usa", "San Francisco, Usa", "Pico Rivera, Usa", "Inglewood, Usa"},
	"Missouri, Usa":       {"Kansas City, Usa", "St Louis, Usa"},
	"Texas, Usa":          {"Dallas, Usa", "Houston, Usa"},
	"Georgia, Usa":        {"Atlanta, Usa"},
	"Massachusetts, Usa":  {"Boston, Usa"},
	"New York, Usa":       {"Brooklyn, Usa", "Newark, Usa", "Uniondale, Usa"},
	"Illinois, Usa":       {"Chicago, Usa", "Berwyn, Usa", "Rosemont, Usa"},
	"Pennsylvania, Usa":   {"Philadelphia, Usa", "Pittsburgh, Usa", "Hershey, Usa"},
	"Michigan, Usa":       {"Grand Rapids, Usa", "Detroit, Usa"},
	"Indiana, Usa":        {"Indianapolis, Usa"},
	"Ohio, Usa":           {"Cleveland, Usa", "Cincinnati, Usa"},
	"Nebraska, Usa":       {"Omaha, Usa"},
	"North Carolina, Usa": {"Charlotte, Usa", "Columbia, Usa"},
	"Louisiana, Usa":      {"New Orleans, Usa"},
	"Wisconsin, Usa":      {"Madison, Usa"},
	"Nevada, Usa":         {"Las Vegas, Usa"},
	// Australian states and cities
	"Victoria, Australia":          {"Melbourne, Australia", "West Melbourne, Australia"},
	"New South Wales, Australia":   {"Sydney, Australia"},
	"Queensland, Australia":        {"Brisbane, Australia"},
	"Western Australia, Australia": {"Burswood, Australia"},
}
