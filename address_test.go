package main

import "testing"

func TestNewAddressBook(t *testing.T) {
	ab := NewAddressBook()
	if ab.Addresses == nil {
		t.Error("Addresses should exist")
	}
}

func TestGetAddressesNone(t *testing.T) {
	ab := NewAddressBook()
	result, err := ab.GetAddresses()

	if err != nil {
		t.Error(err)
	}
	if result == nil {
		t.Error("Result should exist")
	}
	if len(result) > 0 {
		t.Error("Result should have zero elements")
	}
}

func TestGetAddresses(t *testing.T) {
	ab := NewAddressBook()
	ab.AddAddress("first1", "last1", "email1", "phone1")

	result, err := ab.GetAddresses()
	if err != nil {
		t.Error(err)
	}
	if len(result) != 1 {
		t.Error("Result should have 1 element")
	}

	ab.AddAddress("first2", "last2", "email2", "phone2")
	result, err = ab.GetAddresses()
	if err != nil {
		t.Error(err)
	}
	if len(result) != 2 {
		t.Error("Result should have 2 elements")
	}
}

func TestGetAddress(t *testing.T) {
	ab := NewAddressBook()
	addResult, err := ab.AddAddress("firstname", "lastname", "email", "phone")
	if err != nil {
		t.Error(err)
	}

	result, err := ab.GetAddress(addResult.ID)
	if err != nil {
		t.Error(err)
	}
	if result != addResult {
		t.Error("Result should equal expected")
	}
}

func TestAddAddress(t *testing.T) {
	ab := NewAddressBook()
	result, err := ab.AddAddress("first", "last", "email", "phone")

	if err != nil {
		t.Error(err)
	}

	if result.ID == "" {
		t.Error("Result should have an ID")
	}

	expected := Address{
		ID:          result.ID,
		FirstName:   "first",
		LastName:    "last",
		Email:       "email",
		PhoneNumber: "phone",
	}

	if result != expected {
		t.Error("Result should equal expected")
	}
}

func TestUpdateAddress(t *testing.T) {
	ab := NewAddressBook()
	addResult, err := ab.AddAddress("first", "last", "email", "phone")
	if err != nil {
		t.Error(err)
	}

	result, err := ab.UpdateAddress(addResult.ID, "newfirst", "newlast", "newemail", "newphone")
	if err != nil {
		t.Error(err)
	}

	expected := Address{
		ID:          result.ID,
		FirstName:   "newfirst",
		LastName:    "newlast",
		Email:       "newemail",
		PhoneNumber: "newphone",
	}

	if result != expected {
		t.Error("Result should equal expected")
	}
}

func TestDeleteAddress(t *testing.T) {
	ab := NewAddressBook()

	err := ab.DeleteAddress("")
	if err == nil {
		t.Error("Delete address with no ID should throw error")
	}

	err = ab.DeleteAddress("123")
	if err == nil {
		t.Error("Delete address with ID that does not exist should throw error")
	}

	addResult, err := ab.AddAddress("first", "last", "email", "phone")
	if err != nil {
		t.Error(err)
	}

	err = ab.DeleteAddress(addResult.ID)
	if err != nil {
		t.Error(err)
	}

	err = ab.DeleteAddress(addResult.ID)
	if err == nil {
		t.Error("Deleteting address again should not exist and should throw error")
	}
}
