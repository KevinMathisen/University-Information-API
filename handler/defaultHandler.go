package handler

import (
	"fmt"
	"net/http"
)

/*
Handler for default endpoint
*/
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	// Send error if request is not GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

	// Set content type
	w.Header().Set("content-type", "text/html")

	// Information to display to user on root path
	outout := "This website offers three services.<br><a href=\"" +
		DIAG_PATH + "\">" + DIAG_PATH + "</a> - Diagnostic interface<br> <a href=\"" + UNIINFI_PATH +
		"\">" + UNIINFI_PATH + "</a> - University information by name<br> <a href=\"" + NEIGHBOURUNIS_PATH +
		"\">" + NEIGHBOURUNIS_PATH + "</a> - University information by neighbour country and name<br> <a href=\"" + COUNTRYUNIS_PATH +
		"\">" + COUNTRYUNIS_PATH + "</a> - University information by country and name<br> <a href=\"" + COUNTRYALLUNIS_PATH +
		"\">" + COUNTRYALLUNIS_PATH + "</a> - University information by country"

	// Write information to client
	_, err := fmt.Fprintf(w, "%v", outout)

	// Deal with potential errors
	if err != nil {
		http.Error(w, "Error when writing to client", http.StatusInternalServerError)
	}

}
