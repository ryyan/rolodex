package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type AddressHandler struct {
	addressBook
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (ah AddressHandler) Handle(res http.ResponseWriter, req *http.Request) {
	// Log request
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL)

	// Route request
	result, err := ah.route(req)
	if err == nil {
		writeResponse(res, http.StatusOK, result)
	} else {
		errJson, _ := json.Marshal(ErrorResponse{err.Error()})
		writeResponse(res, http.StatusBadRequest, errJson)
	}
}

func (ah AddressHandler) route(req *http.Request) ([]byte, error) {
	// Parse path param
	id := strings.TrimPrefix(req.URL.Path, "/address/")

	switch req.Method {
	case "GET":
		if id == "" {
			result, err := ah.addressBook.GetAddresses()
			resultJson, _ := json.Marshal(result)
			return resultJson, err
		} else {
			result, err := ah.addressBook.GetAddress(id)
			resultJson, _ := json.Marshal(result)
			return resultJson, err
		}

	case "PUT":
		req.ParseForm()
		result, err := ah.addressBook.UpdateAddress(id, req.FormValue("firstname"), req.FormValue("lastname"), req.FormValue("email"), req.FormValue("phonenumber"))
		resultJson, _ := json.Marshal(result)
		return resultJson, err

	case "POST":
		req.ParseForm()
		result, err := ah.addressBook.AddAddress(req.FormValue("firstname"), req.FormValue("lastname"), req.FormValue("email"), req.FormValue("phonenumber"))
		resultJson, _ := json.Marshal(result)
		return resultJson, err

	case "DELETE":
		err := ah.addressBook.DeleteAddress(id)
		return nil, err

	default:
		return nil, nil
	}
}

func writeResponse(res http.ResponseWriter, statusCode int, message []byte) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	res.Write(message)
}
