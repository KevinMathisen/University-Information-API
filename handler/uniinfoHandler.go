package handler

import (
	"net/http"
	"strings"
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

	// Split path into args
	args := strings.Split(r.URL.Path, "/")

	// Check if url is correctly formated
	if (len(args) != 5 && len(args) != 6) || args[4] == "" {
		http.Error(w, "Malformed URL, Expecting format "+UNIINFI_PATH+"name", http.StatusBadRequest)
		return
	}

	// Get universities data by requesting API
	unisReq, err := getUnisData(w, r, args[4], "")
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
