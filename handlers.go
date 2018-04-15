package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

/*
takes a request for a specific shoeName
runs the shoeName through the database and interprets and handles the results

Test in the browser with this URL:
"http://localhost/shoes/adidas Yeezy"
*/
func ShoeShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var shoeName = vars["shoeName"]
	shoe := DBFindShoe(shoeName)
	if shoe.TrueToSize > 0 { //if the Shoe is the empty shoe, its TrueToSize is 0, otherwise the value came from the db
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(shoe); err != nil {
			panic(err)
		}
		return
	}

	//if we didn't find a matching shoe, send 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

/*
parses POST data into Shoe
runs the instance through the database and interprets and handles the results

Test with this curl command:
curl -H "Content-Type: application/json" -d '{"name":"adidas Yeezy","trueToSize":4}' http://localhost/shoes
*/
func ShoeCreate(w http.ResponseWriter, r *http.Request) {
	var shoe Shoe                                                //empty shoe where the JSON will be unmarshalled
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) //limit the reader to mitigate bad actors dumping data
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &shoe); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	t, err := DBCreateShoe(shoe) //try to put the information in the database

	//errors may be caused by violations of constraints, length requirements non-null requirements, etc.
	//they are caught blanket-style and the client is sent 422
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
