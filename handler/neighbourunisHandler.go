package handler

import (
	"errors"
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

	// Get limit from request
	limit, err := getLimitParam(w, r)
	if err != nil {
		return
	}

	// Get country and university name from request
	country, uniName, err := getArgsNURL(w, r)
	if err != nil {
		return
	}

	// Get contires we want to find universities in
	countries, err := getNeighboursCountryReq(w, country)
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
	handleGetRequest(w, r, CONT_TYPE_JSON, unis)

}

/*
Return a list of the names of the countries who border country specified
*/
func getNeighboursCountryReq(w http.ResponseWriter, countryName string) ([]string, error) {
	// List we want to return
	var countries []string

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
Returns all universities with given name in given countries
*/
func getUnisInCountry(w http.ResponseWriter, r *http.Request, uniName string, countries []string, limit int) ([]map[string]interface{}, error) {
	var unis []map[string]interface{}

	// Get universities for each country
	for _, country := range countries {
		// If limit set by user is reached
		if len(unis) >= limit && limit != 0 {
			return unis, nil
		}

		// Get all universities for a given country
		unisReq, err := getUnisReq(w, r, uniName, country)
		if err != nil {
			return unis, err
		}

		// Check if limit is exceded, if so we cut out part of the response
		if len(unis)+len(unisReq) > limit && limit != 0 {
			maxLen := limit - len(unis)
			unisReq = unisReq[:maxLen]
		}

		// Save universities
		unis = append(unis, unisReq...)
	}

	return unis, nil
}

/*
Return limit from URL or 0 if not set, and checks for errors
*/
func getLimitParam(w http.ResponseWriter, r *http.Request) (int, error) {
	// Parse url to get country, university name and limit
	limitstring := (r.URL.Query()).Get("limit")

	// Try to convert limit to a number if limit is specified
	limit, err := strconv.Atoi(limitstring)
	// If there was an error and limit was set by user, or if the limit is less than 0
	if err != nil && limitstring != "" || limit < 0 {
		http.Error(w, "Malformed URL, Invalid limit set ", http.StatusBadRequest)
		return -1, err
	}

	return limit, nil
}

/*
Get arguments country name and university name from URL path, and checks for errors
*/
func getArgsNURL(w http.ResponseWriter, r *http.Request) (string, string, error) {
	// Split path into args
	args := strings.Split(r.URL.Path, "/")

	// Check if URl is correctly formated
	if (len(args) != 6 && len(args) != 7) || args[4] == "" || args[5] == "" {
		http.Error(w, "Malformed URL, Expecting format "+NEIGHBOURUNIS_PATH+"country/uniName{?limit=num}", http.StatusBadRequest)
		return "", "", errors.New("malformed URL")
	}

	return args[4], args[5], nil
}
