package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
Creates and sends a request to specified URL with specified method and content type.
Returns errors if any.
*/
func requestGetFromUrl(url string, method string, contentType string) (http.Response, error) {
	// Create empty response to return in case of error
	var response http.Response

	// Create request
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
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
		return response, err
	}

	//  Return response
	return *res, nil
}

/*
Handles get request when body is of type json
*/
func respondToGetRequest(w http.ResponseWriter, r *http.Request, contentType string, jsonBody interface{}) {
	// Write content type
	w.Header().Add("content-type", contentType)

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
Returns all universities with given name in given countries
*/
func getUnisInCountry(w http.ResponseWriter, r *http.Request, uniName string, countries []string, limit int) ([]map[string]interface{}, error) {
	// List of universeties we want to return from specified country
	var unis []map[string]interface{}

	// Get universities for each country
	for _, country := range countries {
		// If limit set by user is reached
		if len(unis) >= limit && limit != 0 {
			return unis, nil
		}

		// Get all universities for a given country
		unisReq, err := getUnisData(w, r, uniName, country)
		if err != nil {
			return unis, err
		}

		// Check if limit is exceded, if so we cut out part of the response
		if len(unis)+len(unisReq) > limit && limit != 0 {
			maxLen := limit - len(unis)
			unisReq = unisReq[:maxLen]
		}

		// Save universities
		unis = append(unis, unisReq...)
	}

	return unis, nil
}

/*
Get data for all universities from universeties API which satisfy arguments provided.
Can specify university name, country name, or both when searching
Returns a slice of maps which contain data for each university found
*/
func getUnisData(w http.ResponseWriter, r *http.Request, uniName string, country string) ([]map[string]interface{}, error) {

	// List of open-ended map structures, which we can populate with results from hipolab
	var unisReq []map[string]interface{}
	// URL to request data from
	var reqUrl string

	// Create url to request from depending on the arguments provided:
	if uniName != "" && country != "" { // Uni and country specified
		reqUrl = UNI_URL + UNI_SEARCH_PATH + "?name=" + formatURLArg(uniName)
		reqUrl += "&country=" + formatURLArg(country)

	} else if uniName != "" { // Uni specified
		reqUrl = UNI_URL + UNI_SEARCH_PATH + "?name=" + formatURLArg(uniName)

	} else if country != "" { // Country specified
		reqUrl = UNI_URL + UNI_SEARCH_PATH + "?country=" + formatURLArg(country)

	} else { // No argument specified
		http.Error(w, "Malformed URL, required fields not specified", http.StatusBadRequest)
		return unisReq, errors.New("malformed URL")
	}

	log.Println(reqUrl)

	// Get uni from hoplab
	res, err := requestGetFromUrl(reqUrl, http.MethodGet, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return unisReq, err
	}

	// Decode request into unis
	err = json.NewDecoder(res.Body).Decode(&unisReq)
	if err != nil {
		http.Error(w, "Error during decoding", http.StatusInternalServerError)
		return unisReq, err
	}

	return unisReq, nil
}

/*
Creates and return a slice of all universities provided
*/
func createUnisStruct(w http.ResponseWriter, unisReq []map[string]interface{}) ([]Uni, error) {
	// List of university structs we want to return
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

	// Get country uni is located in
	country, err := getCountryData(w, getCountryUni(uniReq), COUNTRY_SEARCH_URL, false)
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
Returns a map of the country specified by requsting the country API
*/
func getCountryData(w http.ResponseWriter, countryName string, searchMethod string, isoSearch bool) (map[string]interface{}, error) {
	// Create country map
	var country map[string]interface{}

	// Create url to request from:
	reqUrl := COUNTRY_URL + searchMethod + countryName

	// Get country from restcountries
	res, err := requestGetFromUrl(reqUrl, http.MethodGet, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return country, err
	}

	// List of open-ended map structures, which we can populate with results from hipolab
	var countriesReq []map[string]interface{}

	// Decode request into countries
	err = json.NewDecoder(res.Body).Decode(&countriesReq)
	if err != nil {
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
	// List of webpages for the provided university we want to return
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
Return the iso of a country from map as a string
*/
func getISOCountry(country map[string]interface{}) string {
	return country["cca2"].(string)
}

/*
Formats an url for a request
*/
func formatURLArg(url string) string {
	// Remove all white space in front of and behind argument
	url = strings.TrimSpace(strings.ReplaceAll(url, "%20", " "))

	// Return argument with all spaces replaced with %20
	return strings.ReplaceAll(url, " ", "%20")
}

/*
Return limit from URL or 0 if not set, and checks for errors
*/
func getLimitParam(w http.ResponseWriter, r *http.Request) (int, error) {
	// Parse url to get country, university name and limit
	limitstring := (r.URL.Query()).Get("limit")

	// Try to convert limit to a number if limit is specified
	limit, err := strconv.Atoi(limitstring)
	// If there was an error and limit was set by user, or if the limit is less than 0
	if err != nil && limitstring != "" || limit < 0 {
		http.Error(w, "Malformed URL, Invalid limit set ", http.StatusBadRequest)
		return -1, err
	}

	return limit, nil
}

/*
Get arguments country name and university name from URL path, and checks for errors
*/
func getArgsCountryUniURL(w http.ResponseWriter, r *http.Request) (string, string, error) {
	// Split path into args
	args := strings.Split(r.URL.Path, "/")

	// Check if URl is correctly formated
	if (len(args) != 6 && len(args) != 7) || args[4] == "" || args[5] == "" {
		http.Error(w, "Malformed URL, Expecting format "+NEIGHBOURUNIS_PATH+"country/uniName{?limit=num}", http.StatusBadRequest)
		return "", "", errors.New("malformed URL")
	}

	// Return name of country and university name provided in the URL path
	return args[4], args[5], nil
}
