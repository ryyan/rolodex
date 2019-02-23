package main

import (
	"log"
	"net/http"
)

func main() {
	ab := AddressBook{}
	http.HandleFunc("/address/", ab.AddressHandler)
	log.Println("Starting serving on port 8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
