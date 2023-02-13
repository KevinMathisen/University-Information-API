package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

var start time.Time = time.Now()

/*
Handler for diagnostic endpoint
*/
func DiagHandler(w http.ResponseWriter, r *http.Request) {

	// Send error if request is not GET:
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

	// Handle get request
	handleGetRequest(w, r)

}

/*
Creates and returns a diag struct
*/
func createDiagRes() Diag {
	// Get status from uni and country source
	uniStatus := getStatus(UNI_URL, http.MethodHead, "")
	countryStatus := getStatus(COUNTRY_URL, http.MethodHead, "")

	// Calulate time elapsed
	elapsed := time.Since(start)

	diagResponse := Diag{
		UniApi:     uniStatus,
		CountryApi: countryStatus,
		Version:    VERSION,
		Uptime:     float64(elapsed),
	}

	return diagResponse
}

/*
Handles get request to diagnostic enpoint
*/
func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	// Write content type
	w.Header().Add("content-type", CONT_TYPE_JSON)

	// Encode content
	err := json.NewEncoder(w).Encode(createDiagRes())
	if err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
