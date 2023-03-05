package handler

import (
	"net/http"
	"strings"
)

/*
Handler for country all universities enpoint
*/
func CountryallunisHandler(w http.ResponseWriter, r *http.Request) {

	// Send error if request is not GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

	// Get limit from request
	limit, err := getLimitParam(w, r)
	if err != nil {
		return
	}

	// Split path into args
	args := strings.Split(r.URL.Path, "/")

	// Check if url is correctly formated
	if (len(args) != 5 && len(args) != 6) || args[4] == "" {
		http.Error(w, "Malformed URL, Expecting format "+COUNTRYALLUNIS_PATH+"name{?limit=num}", http.StatusBadRequest)
		return
	}

	// Create a slice with country provided
	countries := []string{args[4]}

	// Get universities data by requesting API
	unisReq, err := getUnisInCountry(w, r, "", countries, limit)
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
