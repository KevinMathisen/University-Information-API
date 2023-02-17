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

	// Parse url to get country, university name and limit

	// Get universities by university name given
	/*
		unisReq, err := getUnisReq(w, r)
		if err != nil {
			return
		}

		// Filter universeties by country specified and limit if defined

		// Get universities by request
		unis, err := createUnisStruct(w, unisReq)
		if err != nil {
			return
		}

		// Respond with content to user
		handleGetRequest(w, r, CONT_TYPE_JSON, unis)
	*/
}
