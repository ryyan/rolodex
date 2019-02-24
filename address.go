package main

import (
	"errors"
	"math/rand"
)

type Address struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}

type addressBook struct {
	Addresses []Address `json:"addresses"`
}

func NewAddressBook() addressBook {
	return addressBook{make([]Address, 0)}
}

func (ab *addressBook) GetAddresses() ([]Address, error) {
	return ab.Addresses, nil
}

func (ab *addressBook) GetAddress(id string) (Address, error) {
	for _, address := range ab.Addresses {
		if address.ID == id {
			return address, nil
		}
	}

	return Address{}, errors.New("Address not found")
}

func (ab *addressBook) AddAddress(firstName string, lastName string, email string, phoneNumber string) (Address, error) {
	// Verify input
	if firstName == "" && lastName == "" {
		return Address{}, errors.New("Must include either first or last name")
	}

	// Create new address
	address := Address{
		ID:          generateID(12),
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phoneNumber,
	}

	// Add address
	ab.Addresses = append(ab.Addresses, address)
	return address, nil
}

func (ab *addressBook) UpdateAddress(id string, firstName string, lastName string, email string, phoneNumber string) (Address, error) {
	// Verify input
	if id == "" {
		return Address{}, errors.New("Must pass in ID")
	}
	if firstName == "" && lastName == "" {
		return Address{}, errors.New("Must include either first or last name")
	}

	// Find address
	for i := range ab.Addresses {

		// Update address if found
		if ab.Addresses[i].ID == id {
			ab.Addresses[i] = Address{
				ID:          ab.Addresses[i].ID,
				FirstName:   firstName,
				LastName:    lastName,
				Email:       email,
				PhoneNumber: phoneNumber,
			}

			return ab.Addresses[i], nil
		}
	}

	return Address{}, errors.New("Address not found")
}

func (ab *addressBook) DeleteAddress(id string) error {
	// Verify input
	if id == "" {
		return errors.New("Must pass in ID")
	}

	// Find address
	for i := range ab.Addresses {
		if ab.Addresses[i].ID == id {

			// Delete address if found
			ab.Addresses = append(ab.Addresses[:i], ab.Addresses[i+1:]...)
		}
	}

	return nil
}

func (ab *addressBook) ImportCsv() error {
	return nil
}

func (ab *addressBook) ExportCsv() error {
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateID(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
