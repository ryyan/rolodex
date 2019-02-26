package main

import (
	"encoding/json"
	"io/ioutil"
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

func (ah *addressHandler) Handle(res http.ResponseWriter, req *http.Request) {
	// Log request
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL)

	// Route request
	contentType, statusCode, response := ah.route(req)

	// Write response
	writeResponse(res, contentType, statusCode, response)
}

func (ah *addressHandler) route(req *http.Request) (contentType string, statusCode int, response []byte) {
	if req.Header.Get("Content-Type") == "application/csv" {
		statusCode, response = ah.routeCsv(req)
		contentType = "application/csv"
	} else {
		// By default use json router
		statusCode, response = ah.routeJson(req)
		contentType = "application/json"
	}

	return contentType, statusCode, response
}

func (ah *addressHandler) routeCsv(req *http.Request) (statusCode int, response []byte) {
	var err error

	switch req.Method {
	case "GET":
		response, err = ah.ExportCsv()

	case "POST":
		body, _ := ioutil.ReadAll(req.Body)
		err = ah.ImportCsv(body)
	}

	statusCode = http.StatusOK
	if err != nil {
		statusCode = http.StatusBadRequest
		response, _ = json.Marshal(errorResponse{err.Error()})
	}
	return statusCode, response
}

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

func writeResponse(res http.ResponseWriter, contentType string, statusCode int, message []byte) {
	res.Header().Set("Content-Type", contentType)
	res.WriteHeader(statusCode)
	res.Write(message)
}
