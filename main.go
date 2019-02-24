package main

import (
	"log"
	"net/http"
)

func main() {
	addressBook := NewAddressBook()
	addressHandler := AddressHandler{addressBook}
	http.HandleFunc("/address/", addressHandler.Handle)
	log.Println("Starting serving on port 8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
