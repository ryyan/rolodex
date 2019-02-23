package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Address struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}

type AddressBook struct {
	Addresses []Address
}

func (ab AddressBook) AddressHandler(res http.ResponseWriter, req *http.Request) {
	// Log request
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL)

	// Parse path param
	id := strings.TrimPrefix(req.URL.Path, "/address/")
	log.Println(id)

	// Return response
	json.NewEncoder(res).Encode(ab.Addresses)
}

func (ab AddressBook) GetAddresses() ([]Address, error) {
}

func (ab AddressBook) GetAddress(id string) (Address, error) {
}

func (ab AddressBook) AddAddress(firstName string, lastName string, email string, phoneNumber string) (Address, error) {
}

func (ab AddressBook) UpdateAddress(id string, firstName string, lastName string, email string, phoneNumber string) (Address, error) {
}

func (ab AddressBook) DeleteAddress(id string) error {
}

func (ab AddressBook) ImportCsv() error {
}

func (ab AddressBook) ExportCsv() error {
}
