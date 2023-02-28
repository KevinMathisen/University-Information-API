package handler

import (
	"net/http"
	"time"
)

/*
Handler for diagnostic endpoint
*/
func DiagHandler(w http.ResponseWriter, r *http.Request, start time.Time) {

	// Send error if request is not GET:
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

	// Generate diagnosis reponse
	diagRes, err := createDiagRes(start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle get request
	handleGetRequest(w, r, CONT_TYPE_JSON, diagRes)

}

/*
Creates and returns a diag struct
*/
func createDiagRes(start time.Time) (Diag, error) {

	// Get request from uni  source
	resUni, err := Request(UNI_URL, http.MethodHead, "")
	if err != nil {
		return Diag{}, err
	}

	// Get request from country source
	resCountry, err := Request(COUNTRY_URL, http.MethodHead, "")
	if err != nil {
		return Diag{}, err
	}

	// Calulate time elapsed
	elapsed := time.Since(start)

	// Initialize the diag response struct
	diagResponse := Diag{
		UniApi:     resUni.Status,
		CountryApi: resCountry.Status,
		Version:    VERSION,
		Uptime:     float64(elapsed.Seconds()),
	}

	return diagResponse, nil
}
