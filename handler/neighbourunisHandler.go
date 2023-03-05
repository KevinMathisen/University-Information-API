package handler

import (
	"net/http"
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

	// Get limit from request
	limit, err := getLimitParam(w, r)
	if err != nil {
		return
	}

	// Get country and university name from request
	country, uniName, err := getArgsCountryUniURL(w, r)
	if err != nil {
		return
	}

	// Get contires we want to find universities in
	countries, err := getNeighboursCountry(w, country)
	if err != nil {
		return
	}

	// Get unisversities for each border country
	unisReq, err := getUnisInCountry(w, r, uniName, countries, limit)
	if err != nil {
		return
	}

	// Create universities struct
	unis, err := createUnisStruct(w, unisReq)
	if err != nil {
		return
	}

	// Respond with content to user
	respondToGetRequest(w, r, CONT_TYPE_JSON, unis)

}

/*
Return a list of the names of the countries who border country specified
*/
func getNeighboursCountry(w http.ResponseWriter, countryName string) ([]string, error) {
	// List of neighbouring countries we want to return
	var countries []string

	// Get country specified
	country, err := getCountryData(w, countryName, COUNTRY_SEARCH_URL, false)
	if err != nil {
		return countries, err
	}

	// For each bordering country, get country name and add to array
	for _, border := range country["borders"].([]interface{}) {
		country, err = getCountryData(w, border.(string), ISO_SEARCH_URL, true)
		if err != nil {
			return countries, err
		}
		countries = append(countries, getNameCountry(country))
	}

	return countries, nil

}
