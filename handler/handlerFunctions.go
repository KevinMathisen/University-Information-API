package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
Creates and sends a request to specified URL with specified method and content type.
Returns errors if any.
*/
func Request(url string, method string, contentType string) (http.Response, error) {

	// Create empty response to return in case of error
	var response http.Response

	// Create request
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Errorf("Error when creating request", err.Error())
		return response, err
	}

	// Set content type
	r.Header.Add("content-type", contentType)

	// Set up client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue http request
	res, err := client.Do(r)
	if err != nil {
		fmt.Errorf("Error in response", err.Error())
		return response, err
	}

	//  Return response
	return *res, nil
}

/*
Handles get request to diagnostic enpoint
*/
func handleGetRequest(w http.ResponseWriter, r *http.Request, contentType string, jsonBody interface{}) {
	// Write content type
	w.Header().Add("content-type", CONT_TYPE_JSON)

	// Encode content and write to response
	err := json.NewEncoder(w).Encode(jsonBody)
	if err != nil {
		http.Error(w, "Error during encoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Manually set response http status to ok
	w.WriteHeader(http.StatusOK)
}
