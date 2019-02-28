package main

import (
	"log"
	"net/http"
)

func main() {
	// Initiailize address handler
	addressBook := NewAddressBook()
	addressHandler := NewAddressHandler(addressBook)

	// Configure and start server
	http.HandleFunc("/address/", addressHandler.Handle)
	log.Println("Starting server on port 8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
