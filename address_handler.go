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

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewAddressHandler(ab *addressBook) *addressHandler {
	return &addressHandler{ab}
}

func (ah *addressHandler) Handle(res http.ResponseWriter, req *http.Request) {
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

func (ah *addressHandler) route(req *http.Request) ([]byte, error) {
	// Parse path param
	id := strings.TrimLeft(req.URL.Path, "/address")

	switch req.Method {
	case "GET":
		if id == "" {
			result, err := ah.GetAddresses()
			resultJson, _ := json.Marshal(result)
			return resultJson, err
		} else {
			result, err := ah.GetAddress(id)
			resultJson, _ := json.Marshal(result)
			return resultJson, err
		}

	case "PUT":
		req.ParseForm()
		result, err := ah.UpdateAddress(id, req.FormValue("firstname"), req.FormValue("lastname"), req.FormValue("email"), req.FormValue("phonenumber"))
		resultJson, _ := json.Marshal(result)
		return resultJson, err

	case "POST":
		req.ParseForm()
		result, err := ah.AddAddress(req.FormValue("firstname"), req.FormValue("lastname"), req.FormValue("email"), req.FormValue("phonenumber"))
		resultJson, _ := json.Marshal(result)
		return resultJson, err

	case "DELETE":
		err := ah.DeleteAddress(id)
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
