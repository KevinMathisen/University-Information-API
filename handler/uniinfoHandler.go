package handler

import (
	"net/http"
)

/*
Handler for University information endpoint
*/
func UniinfoHandler(w http.ResponseWriter, r *http.Request) {

	// Send error if request is not GET:
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method, currently only GET is supported", http.StatusNotImplemented)
		return
	}

}
