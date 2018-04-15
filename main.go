package main

import (
	"log"
	"net/http"
)

/*
main entry point
*/
func main() {
	defer db.Close() //close the database when the program exits
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":80", router)) //start listening on port 80, pass all requests to custom router
}
