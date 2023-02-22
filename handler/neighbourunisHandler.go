package handler

import (
	"log"
	"net/http"
	"strconv"
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

	// Parse url to get country, university name and limit
	limitstring := (r.URL.Query()).Get("limit")

	// Try to convert limit to a number if limit is specified
	limit, err := strconv.Atoi(limitstring)
	if err != nil && limitstring != "" || limit < 0 {
		log.Println("Error limit: " + err.Error())
		http.Error(w, "Malformed URL, Invalid limit set ", http.StatusBadRequest)
		return
	}

	// Split path into args
	args := strings.Split(r.URL.Path, "/")

	// Check if URl is correctly formated
	if len(args) != 6 || args[4] == "" || args[5] == "" {
		http.Error(w, "Malformed URL, Expecting format "+NEIGHBOURUNIS_PATH+"country/uniName{?limit=num}", http.StatusBadRequest)
		return
	}

	// Get contires we want to find universeties in
	countries, err := getNeighboursCountryReq(w, args[4])
	if err != nil {
		return
	}

	// Get unisversities for each border country
	unisReq, err := getUnisInCountry(w, r, args[5], countries, limit)
	if err != nil {
		return
	}

	log.Println(unisReq)

	// Get universities by request
	unis, err := createUnisStruct(w, unisReq)
	if err != nil {
		return
	}

	// Respond with content to user
	handleGetRequest(w, r, CONT_TYPE_JSON, unis)

}

/*
Return a list of the names of the countries who border country specified
*/
func getNeighboursCountryReq(w http.ResponseWriter, countryName string) ([]string, error) {
	// List we want to return
	var countries []string

	// Add Original country to list of countries.
	// NB! If assignment only wants neighbours, not including country itself, remove this line! ----------------------
	countries = append(countries, countryName)

	// Get country specified
	country, err := getCountryReq(w, countryName, COUNTRY_SEARCH_URL, false)
	if err != nil {
		return countries, err
	}

	// For each bordering country, get country name and add to array
	for _, border := range country["borders"].([]interface{}) {
		country, err = getCountryReq(w, border.(string), ISO_SEARCH_URL, true)
		if err != nil {
			return countries, err
		}
		countries = append(countries, getNameCountry(country))
	}

	return countries, nil

}

/*
Returns all universeties with given name in given countries
*/
func getUnisInCountry(w http.ResponseWriter, r *http.Request, uniName string, countries []string, limit int) ([]map[string]interface{}, error) {
	var unis []map[string]interface{}

	// Get universeties for each country
	for i, country := range countries {
		// If limit set by user is reached
		if i >= limit && limit != 0 {
			return unis, nil
		}

		log.Println("Unireq " + uniName + " " + country)

		// Get all universeties for a given country
		unisReq, err := getUnisReq(w, r, uniName, country)
		if err != nil {
			return unis, err
		}

		// Save university
		unis = append(unis, unisReq...)
	}

	return unis, nil
}
