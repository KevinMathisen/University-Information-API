package handler

import (
	"fmt"
	"net/http"
)

func Request(url string, method string, contentType string) http.Response {

	// Create request
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Errorf("Error when creating request", err.Error())
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
	}

	//  Return response
	return *res
}

func getStatus(url string, method string, contentType string) string {
	// Issue request
	return Request(url, method, contentType).Status
}
