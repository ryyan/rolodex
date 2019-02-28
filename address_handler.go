package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type addressHandler struct {
	*addressBook
}

type addressResponse struct {
	Rolodex []Address `json:"rolodex"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewAddressHandler(ab *addressBook) *addressHandler {
	return &addressHandler{ab}
}

// Handle satisfies the http.Handler for the address endpoint
func (ah *addressHandler) Handle(res http.ResponseWriter, req *http.Request) {
	// Log request
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL)

	// Route request
	contentType, statusCode, response := ah.route(req)

	// Write response
	writeResponse(res, contentType, statusCode, response)
}

// route routes a request to an appropriate content type router
func (ah *addressHandler) route(req *http.Request) (contentType string, statusCode int, response []byte) {
	if req.Header.Get("Content-Type") == "text/csv" {
		statusCode, response = ah.routeCsv(req)
		return "text/csv", statusCode, response
	} else {
		// By default use json router
		statusCode, response = ah.routeJson(req)
		return "application/json", statusCode, response
	}
}

// routeCsv routes a request to a csv method
func (ah *addressHandler) routeCsv(req *http.Request) (statusCode int, response []byte) {
	var err error

	switch req.Method {
	case "GET":
		response, err = ah.ExportCsv()

	case "POST":
		err = ah.ImportCsv(req.Body)
	}

	if err != nil {
		statusCode = http.StatusBadRequest
		response, _ = json.Marshal(errorResponse{err.Error()})
	} else {
		statusCode = http.StatusOK
	}
	return statusCode, response
}

// routeJson routes a request to a json method
func (ah *addressHandler) routeJson(req *http.Request) (statusCode int, response []byte) {
	var address Address
	var result []Address
	var err error

	// Parse path param
	id := strings.TrimLeft(req.URL.Path, "/address")

	switch req.Method {
	case "GET":
		if id == "" {
			result, err = ah.GetAddresses()
		} else {
			address, err = ah.GetAddress(id)
			result = []Address{address}
		}

	case "PUT":
		req.ParseForm()
		address, err = ah.UpdateAddress(id, req.FormValue("firstname"), req.FormValue("lastname"), req.FormValue("email"), req.FormValue("phonenumber"))
		result = []Address{address}

	case "POST":
		req.ParseForm()
		address, err = ah.AddAddress(req.FormValue("firstname"), req.FormValue("lastname"), req.FormValue("email"), req.FormValue("phonenumber"))
		result = []Address{address}

	case "DELETE":
		err = ah.DeleteAddress(id)
	}

	// Return response
	if err != nil {
		statusCode = http.StatusBadRequest
		response, _ = json.Marshal(errorResponse{err.Error()})
	} else {
		statusCode = http.StatusOK
		response, _ = json.Marshal(addressResponse{result})
	}
	return statusCode, response
}

// writeResponse writes a byte message to an http response writer
func writeResponse(res http.ResponseWriter, contentType string, statusCode int, message []byte) {
	res.Header().Set("Content-Type", contentType)
	res.WriteHeader(statusCode)
	res.Write(message)
}
