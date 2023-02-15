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

	// Get universities by request
	unis, err := getUnis(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with content to user
	handleGetRequest(w, r, CONT_TYPE_JSON, unis)

}

/*
Get all universeties from hipolab with arguments provided by client in request
*/
func getUnis(r *http.Request) ([]Uni, error) {
	args := strings.Split(r.URL.Path, "/")
	bname := args[len(args)-1]

	var uni Uni
	uni.Name = bname

	return []Uni{uni}, nil
}
