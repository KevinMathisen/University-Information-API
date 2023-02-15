package handler

import (
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

	// Generate diagnosis reponse
	diagRes, err := createDiagRes()
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
func createDiagRes() (Diag, error) {

	// Get request from uni and country source
	resUni, err := Request(UNI_URL, http.MethodHead, "")
	if err != nil {
		return Diag{}, err
	}

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

/*
Handles get request to diagnostic enpoint

func handleGetRequestDiag(w http.ResponseWriter, r *http.Request) {
	// Write content type
	w.Header().Add("content-type", CONT_TYPE_JSON)

	// Generate diagnosis reponse
	diagRes, err := createDiagRes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode content and write to response
	err = json.NewEncoder(w).Encode(diagRes)
	if err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Manually set response http status to ok
	w.WriteHeader(http.StatusOK)
}
*/
