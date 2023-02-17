package handler

import (
	"net/http"
	"strings"
)

/*
Handler for Neighbourunis endpoint
*/
func NeighbourunisHandler(w http.ResponseWriter, r *http.Request) {

	// Send error if request is not GET:
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

	limit := 0
	// Parse url to get country, university name and limit
	// If user specified limit set limit, if not set -1
	if value := (r.URL.Query()).Get("limit"); value {
		limit = value		// Convert to int, TODO
	}

	// Split path into args
	args := strings.Split(r.URL.Path, "/")

	// Check if URl is correctly formated
	if len(args) != 4 || args[2] == "" || args[3] == "" {
		http.Error(w, "Malformed URL, Expecting format "+NEIGHBOURUNIS_PATH+"country/uniName{?limit=num}", http.StatusBadRequest)
		return
	}

	// Get universities by university name given
	unisReq, err := getUnisReq(w, r, args[3])
	if err != nil {
		return
	}

	// Get contires we want to find universeties in
	countries, err := getNeighboursCountryReq(w, args[2])
	if err != nil {
		return
	}

	var unisFiltered []map[string]interface{} 

	// For all universeites, filter away those not in target countries
	//  and set a max amount of allowed unis based on limit specified
	for i, uni := range unisReq {
		if find(getCountryUni(uni), countries) && i < limit {
			unisFiltered = append(unisFiltered, uni)
		}
	}

	// Get universities by request
	unis, err := createUnisStruct(w, unisFiltered)
	if err != nil {
		return
	}

	// Respond with content to user
	handleGetRequest(w, r, CONT_TYPE_JSON, unis)

}

/*
Return a list of the names of the countries who border country specified
*/
func getNeighboursCountryReq(w http.ResponseWriter, countryName string) ([]string, err) {
	// List we want to return
	var countries []string

	// Add Original country to list of countries.
	// NB! If assignment only wants neighbours, not including country itself, remove this line! ----------------------
	countries = append(countries, countryName)

	// Get country specified 
	country, err := getCountryReq(w, countryName, COUNTRY_SEARCH_URL)
	if err != nil {

	}

	// For each bordering country, get country name and add to array
	for _, border := range country["borders"].([]interface{}) {
		country, err = getCountryReq(w, border, ISO_SEARCH_URL)
		if err != nil {
			return countries, err
		}
		countries = append(countries, getNameCountry(country))
	}

	return countries, nil

	
}