package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	// Initialize handler
	addressBook := NewAddressBook()
	addressHandler := NewAddressHandler(addressBook)

	// Initialize server
	server := httptest.NewServer(http.HandlerFunc(addressHandler.Handle))
	defer server.Close()
	time.Sleep(100)

	// Initialize client
	client := http.Client{}

	// Variables used for json unmarshalling
	var addressResponse addressResponse
	var errorResponse errorResponse

	// GET addresses should return empty array
	response, err := http.Get(server.URL + "/address")
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Error("Should return 200")
	}

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &addressResponse)
	if addressResponse.Rolodex == nil || len(addressResponse.Rolodex) > 0 {
		t.Error("Address list should be an empty")
	}

	// GET an address that does not exist should return error
	response, err = http.Get(server.URL + "/address/thisShouldNotExist")
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 400 {
		t.Error("Should return 400")
	}

	body, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &errorResponse)
	if errorResponse.Error != "Address not found" {
		t.Error("Should have error message")
	}

	// POST to create an address
	form := url.Values{}
	form.Add("firstname", "myFirstName")
	form.Add("lastname", "myLastName")
	form.Add("email", "myEmail")
	form.Add("phonenumber", "myPhoneNumber")
	request, _ := http.NewRequest("POST", server.URL+"/address", strings.NewReader(form.Encode()))
	request.PostForm = form
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err = client.Do(request)
	if err != nil {
		t.Error(err)
	}

	body, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &addressResponse)

	expected := Address{
		ID:          addressResponse.Rolodex[0].ID,
		FirstName:   "myFirstName",
		LastName:    "myLastName",
		Email:       "myEmail",
		PhoneNumber: "myPhoneNumber",
	}
	if addressResponse.Rolodex[0] != expected {
		t.Error("Address returned does not equal expected")
	}

	// GET all addresses should have 1
	response, err = http.Get(server.URL + "/address")
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Error("Should return 200")
	}

	body, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &addressResponse)
	if addressResponse.Rolodex == nil || len(addressResponse.Rolodex) != 1 {
		t.Error("Address list should have 1 elements")
	}
	firstID := addressResponse.Rolodex[0].ID

	// GET the address created
	response, err = http.Get(server.URL + "/address/" + firstID)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Error("Should return 200")
	}

	body, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &addressResponse)

	if addressResponse.Rolodex[0] != expected {
		t.Error("Address returned does not equal expected")
	}

	// POST to create another address
	form = url.Values{}
	form.Add("firstname", "myFirstName2")
	form.Add("lastname", "myLastName2")
	form.Add("email", "myEmail2")
	form.Add("phonenumber", "myPhoneNumber2")
	request, _ = http.NewRequest("POST", server.URL+"/address", strings.NewReader(form.Encode()))
	request.PostForm = form
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err = client.Do(request)
	if err != nil {
		t.Error(err)
	}

	body, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &addressResponse)
	secondID := addressResponse.Rolodex[0].ID

	// PUT to update the address just created
	form = url.Values{}
	form.Add("firstname", "myFirstName3")
	form.Add("lastname", "myLastName3")
	form.Add("email", "myEmail3")
	form.Add("phonenumber", "myPhoneNumber3")
	request, _ = http.NewRequest("PUT", server.URL+"/address/"+secondID, strings.NewReader(form.Encode()))
	request.PostForm = form
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err = client.Do(request)
	if err != nil {
		t.Error(err)
	}

	body, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &addressResponse)
	expected = Address{
		ID:          secondID,
		FirstName:   "myFirstName3",
		LastName:    "myLastName3",
		Email:       "myEmail3",
		PhoneNumber: "myPhoneNumber3",
	}

	if addressResponse.Rolodex[0] != expected {
		t.Error("Address returned does not equal expected")
	}

	// GET all addresses should have 2
	response, err = http.Get(server.URL + "/address")
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Error("Should return 200")
	}

	body, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &addressResponse)
	if addressResponse.Rolodex == nil || len(addressResponse.Rolodex) != 2 {
		t.Error("Address list should have 2 elements")
	}

	// DELETE the second created address
	request, _ = http.NewRequest("DELETE", server.URL+"/address/"+secondID, nil)
	response, err = client.Do(request)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Error("Should return 200")
	}

	// GET all addresses should have 1
	response, err = http.Get(server.URL + "/address")
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Error("Should return 200")
	}

	body, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &addressResponse)
	if addressResponse.Rolodex == nil || len(addressResponse.Rolodex) != 1 {
		t.Error("Address list should have 1 elements")
	}

	// GET an address that does not exist should return error
	response, err = http.Get(server.URL + "/address/" + secondID)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 400 {
		t.Error("Should return 400")
	}

	body, _ = ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &errorResponse)
	if errorResponse.Error != "Address not found" {
		t.Error("Should have error message")
	}
}

func TestHandlerCsv(t *testing.T) {
	// Initialize handler
	addressBook := NewAddressBook()
	addressHandler := NewAddressHandler(addressBook)

	// Initialize server
	server := httptest.NewServer(http.HandlerFunc(addressHandler.Handle))
	defer server.Close()
	time.Sleep(100)

	// Initialize client
	client := http.Client{}

	// Variables used for json unmarshalling
	var addressResponse addressResponse

	// POST to import a CSV
	csvImport := `firstname,lastname,email,phonenumber
firstname1,lastname1,email1,phoneumber1
firstname2,lastname2,email2,phoneumber2
firstname3,lastname3,email3,phoneumber3
firstname4,lastname4,email4,phoneumber4
`
	request, _ := http.NewRequest("POST", server.URL+"/address", strings.NewReader(csvImport))
	request.Header.Add("Content-Type", "text/csv")
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Error("Should return 200")
	}

	// GET all addresses should have 4
	response, err = http.Get(server.URL + "/address")
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Error("Should return 200")
	}

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &addressResponse)
	if addressResponse.Rolodex == nil || len(addressResponse.Rolodex) != 4 {
		t.Error("Address list should have 4 elements")
	}

	// GET to export a CSV
	request, _ = http.NewRequest("GET", server.URL+"/address", nil)
	request.Header.Add("Content-Type", "text/csv")
	response, err = client.Do(request)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != 200 {
		t.Error("Should return 200")
	}

	r := csv.NewReader(response.Body)
	records, err := r.ReadAll()
	if err != nil {
		t.Error(err)
	}
	if len(records) != 5 {
		t.Error("Records should have 5 elements")
	}
}
