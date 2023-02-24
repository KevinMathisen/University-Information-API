package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
		log.Println("Error when creating request" + err.Error())
		fmt.Errorf("Error when creating request", err.Error())
		return response, err
	}

	// Set content type
	r.Header.Add("content-type", contentType)

	// Set up client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	log.Println(url)
	// Issue http request
	res, err := client.Do(r)
	if err != nil {
		log.Println("Error in response" + err.Error())
		fmt.Errorf("Error in response", err.Error())
		return response, err
	}

	//  Return response
	return *res, nil
}

/*
Handles get request when body is of type json
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

/*
Get all universities from hipolab with arguments provided by client in request
*/
func getUnisReq(w http.ResponseWriter, r *http.Request, uniName string, country string) ([]map[string]interface{}, error) {

	// List of open-ended map structures, which we can populate with results from hipolab
	var unisReq []map[string]interface{}

	// Create url to request from:							TODO: add + "name="
	reqUrl := UNI_URL + UNI_SEARCH_PATH + "name=" + formatURLArg(uniName)

	// Check if country specified
	if country != "" {
		reqUrl += "&country=" + formatURLArg(country)
	}

	// Get uni from hoplab
	res, err := Request(reqUrl, http.MethodGet, "")
	if err != nil {
		log.Println("Error during request: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return unisReq, err
	}

	// Decode request into unis
	err = json.NewDecoder(res.Body).Decode(&unisReq)
	if err != nil {
		log.Println("Error during encoding: " + err.Error())
		http.Error(w, "Error during decoding", http.StatusInternalServerError)
		return unisReq, err
	}

	return unisReq, nil
}

/*
Creates and return a slice of all universities provided
*/
func createUnisStruct(w http.ResponseWriter, unisReq []map[string]interface{}) ([]Uni, error) {
	var unis []Uni

	// For each university we got create a uni struct, and add it to our response
	for _, value := range unisReq {
		uni, err := createUniStruct(w, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return unis, err
		}
		unis = append(unis, uni)
	}

	return unis, nil
}

/*
Create a uni struct with all information by asking restcountries
*/
func createUniStruct(w http.ResponseWriter, uniReq map[string]interface{}) (Uni, error) {

	// Get country
	country, err := getCountryReq(w, getCountryUni(uniReq), COUNTRY_SEARCH_URL, false)
	if err != nil {
		return Uni{}, err
	}

	// Assemble the struct
	uni := Uni{
		Name:      uniReq["name"].(string),
		Country:   getCountryUni(uniReq),
		Isocode:   getISOCountry(country),
		Webpages:  getWebpagesUni(uniReq),
		Languages: getLanguagesCountry(country),
		Map:       getMapCountry(country),
	}

	return uni, nil
}

/*
Returns a map of the country specified
*/
func getCountryReq(w http.ResponseWriter, countryName string, searchMethod string, isoSearch bool) (map[string]interface{}, error) {
	// Create country map
	var country map[string]interface{}

	// Create url to request from:
	reqUrl := COUNTRY_URL + searchMethod + countryName

	// Get country from restcountries
	res, err := Request(reqUrl, http.MethodGet, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return country, err
	}

	// List of open-ended map structures, which we can populate with results from hipolab
	var countriesReq []map[string]interface{}

	// Decode request into countries
	err = json.NewDecoder(res.Body).Decode(&countriesReq)
	if err != nil {
		log.Println("Error during encoding: " + err.Error())
		http.Error(w, "Error during decoding", http.StatusInternalServerError)
		return country, err
	}

	// In case of multiple countries with similar name, get the one we want from the list
	for _, countryReq := range countriesReq {
		if strings.EqualFold(getNameCountry(countryReq), countryName) || isoSearch {
			country = countryReq
		}
	}

	return country, nil
}

/*
Return the name of a country from map as a string
*/
func getNameCountry(country map[string]interface{}) string {
	return (country["name"].(map[string]interface{}))["common"].(string)
}

/*
Return the openstreetmap from a country as a string
*/
func getMapCountry(country map[string]interface{}) string {
	return (country["maps"].(map[string]interface{}))["openStreetMaps"].(string)
}

/*
Return all languages as a map of strings from a country
*/
func getLanguagesCountry(country map[string]interface{}) map[string]string {
	languages := make(map[string]string)

	// Convert each language in country map to a string and add to new map
	for key, language := range country["languages"].(map[string]interface{}) {
		languages[key] = language.(string)
	}

	return languages
}

/*
Return all webpages from a university as a list of strings
*/
func getWebpagesUni(uni map[string]interface{}) []string {
	var webpages []string

	// Convert each webpage in uni map to a string and add to new list
	for _, webpage := range uni["web_pages"].([]interface{}) {
		webpages = append(webpages, webpage.(string))
	}

	return webpages
}

/*
Return country of a university as a string
*/
func getCountryUni(uni map[string]interface{}) string {
	return uni["country"].(string)
}

/*
Returns if given value exists in array non case sensetive
*/
func find(value string, array []string) bool {
	for _, v := range array {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}

/*
Return the iso of a country from map as a string
*/
func getISOCountry(country map[string]interface{}) string {
	return country["cca2"].(string)
}

/*
Formats an url for a request
*/
func formatURLArg(url string) string {
	url = strings.TrimSpace(strings.ReplaceAll(url, "%20", " "))
	return strings.ReplaceAll(url, " ", "%20")
}
