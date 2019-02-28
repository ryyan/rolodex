package main

import (
	"encoding/csv"
	"errors"
	"io"
	"math/rand"
	"strings"
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

// addressBook holds the list of addresses
type addressBook struct {
	Addresses []Address
	*sync.Mutex
}

func NewAddressBook() *addressBook {
	return &addressBook{make([]Address, 0), &sync.Mutex{}}
}

// GetAddresses returns the list of all addresses
func (ab *addressBook) GetAddresses() ([]Address, error) {
	ab.Lock()
	defer ab.Unlock()

	return ab.Addresses, nil
}

// GetAddress returns an address by ID
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

// AddAddress adds an address and returns the newly created address
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

// UpdateAddress updates an address and returns the updated address
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

// DeleteAddress deletes an address
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

// ImportCsv adds a list of addresses based on the passed in csv
func (ab *addressBook) ImportCsv(csvFile io.ReadCloser) error {
	ab.Lock()
	defer ab.Unlock()

	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) != 4 {
			return errors.New("Expected 4 fields when importing CSV")
		}
		if record[0] == "firstname" && record[1] == "lastname" {
			// Skip header row
			continue
		}

		address := Address{
			FirstName:   record[0],
			LastName:    record[1],
			Email:       record[2],
			PhoneNumber: record[3],
		}
		ab.Addresses = append(ab.Addresses, address)
	}

	return nil
}

// ExportCsv returns a csv file of the list of all addresses
func (ab *addressBook) ExportCsv() ([]byte, error) {
	ab.Lock()
	defer ab.Unlock()

	var sb strings.Builder
	sb.WriteString("firstname,lastname,email,phonenumber\n")

	for _, address := range ab.Addresses {
		sb.WriteString(address.FirstName)
		sb.WriteString(",")
		sb.WriteString(address.LastName)
		sb.WriteString(",")
		sb.WriteString(address.Email)
		sb.WriteString(",")
		sb.WriteString(address.PhoneNumber)
		sb.WriteString("\n")
	}

	return []byte(sb.String()), nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generateID returns a randomly generated string
func generateID(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}
	return string(b)
}
