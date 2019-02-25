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
	Error   string    `json:"error"`
}

func NewAddressHandler(ab *addressBook) *addressHandler {
	return &addressHandler{ab}
}

func NewAddressResponse(response []Address, err error) addressResponse {
	if err != nil {
		return addressResponse{response, err.Error()}
	}
	return addressResponse{response, ""}
}

func (ah *addressHandler) Handle(res http.ResponseWriter, req *http.Request) {
	// Log request
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL)

	// Route request
	result, err := ah.route(req)

	// Return json response
	resultJson, _ := json.Marshal(NewAddressResponse(result, err))
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusBadRequest
	}
	writeJsonResponse(res, statusCode, resultJson)
}

func (ah *addressHandler) route(req *http.Request) ([]Address, error) {
	// Parse path param
	id := strings.TrimLeft(req.URL.Path, "/address")

	switch req.Method {
	case "GET":
		if id == "" {
			return ah.GetAddresses()
		} else {
			result, err := ah.GetAddress(id)
			return []Address{result}, err
		}

	case "PUT":
		req.ParseForm()
		result, err := ah.UpdateAddress(id, req.FormValue("firstname"), req.FormValue("lastname"), req.FormValue("email"), req.FormValue("phonenumber"))
		return []Address{result}, err

	case "POST":
		req.ParseForm()
		result, err := ah.AddAddress(req.FormValue("firstname"), req.FormValue("lastname"), req.FormValue("email"), req.FormValue("phonenumber"))
		return []Address{result}, err

	case "DELETE":
		err := ah.DeleteAddress(id)
		return nil, err

	default:
		return nil, nil
	}
}

func writeJsonResponse(res http.ResponseWriter, statusCode int, message []byte) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	res.Write(message)
}
