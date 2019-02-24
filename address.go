package main

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type Address struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}

type addressBook struct {
	Addresses []Address
	*sync.Mutex
}

func NewAddressBook() *addressBook {
	return &addressBook{make([]Address, 0), &sync.Mutex{}}
}

func (ab *addressBook) GetAddresses() ([]Address, error) {
	ab.Lock()
	defer ab.Unlock()

	return ab.Addresses, nil
}

func (ab *addressBook) GetAddress(id string) (Address, error) {
	ab.Lock()
	defer ab.Unlock()

	for i := range ab.Addresses {
		if ab.Addresses[i].ID == id {
			return ab.Addresses[i], nil
		}
	}

	return Address{}, errors.New("Address not found")
}

func (ab *addressBook) AddAddress(firstName string, lastName string, email string, phoneNumber string) (Address, error) {
	ab.Lock()
	defer ab.Unlock()

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
	ab.Lock()
	defer ab.Unlock()

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
	ab.Lock()
	defer ab.Unlock()

	// Verify input
	if id == "" {
		return errors.New("Must pass in ID")
	}

	// Find address
	for i := range ab.Addresses {
		if ab.Addresses[i].ID == id {

			// Delete address if found
			ab.Addresses = append(ab.Addresses[:i], ab.Addresses[i+1:]...)

			return nil
		}
	}

	return errors.New("Address not found")
}

func (ab *addressBook) ImportCsv() error {
	ab.Lock()
	defer ab.Unlock()

	return nil
}

func (ab *addressBook) ExportCsv() error {
	ab.Lock()
	defer ab.Unlock()

	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateID(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}
	return string(b)
}
