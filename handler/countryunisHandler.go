package handler

import (
	"net/http"
)

/*
Handler for country universities endpoint
*/
func CountryunisHandler(w http.ResponseWriter, r *http.Request) {

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

	// Create a slice with country provided
	countries := []string{country}

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
