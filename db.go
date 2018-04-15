package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" //postgres driver for golang
)

var db *sql.DB

const (
	DB_USER     = "postgres" //default db user for simplicity, replace with non superuser in prod
	DB_PASSWORD = "password" //storing this in source for simplicity, replace with env variable in prod
	DB_NAME     = "postgres" //default db name for simplicity, replace with app-specific in prod
)

/*
opens a connection to the postgres database
creates an empty shoeinfo table if it doesn't exist
*/
func init() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, _ = sql.Open("postgres", dbinfo)

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS shoeinfo (
			shoe character varying(100) NOT NULL,
			trueToSize integer,
			CONSTRAINT CHK_Shoe CHECK (trueToSize>=1 AND trueToSize<=5)
		) WITH (OIDS=FALSE);`) //constraint added to ensure croudsourced data is within acceptable ranges
	checkErr(err)

	fmt.Println("# DB opened")
}

/*
convenience function to panic if non-nil error
lots of things can go wrong when querying/reading db, so most actions are funneled here
*/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/*
given the name of a shoe, get the average trueToSize value from the database
return a Shoe with the information to be JSON-encoded back to the client
*/
func DBFindShoe(shoeName string) Shoe {
	fmt.Println("# Finding")
	stmt, err := db.Prepare("SELECT AVG(trueToSize) FROM shoeinfo WHERE shoe=$1") //prepare the statement first...
	checkErr(err)

	rows, err := stmt.Query(shoeName) //...then query the database
	checkErr(err)

	//even though there will always be a single row as the result of SELECT AVG(),
	//this will move the scanner to the proper position
	for rows.Next() {
		var trueToSize float32 //even though individual values are stored as integers, the average is desired as float
		if err := rows.Scan(&trueToSize); err == nil {
			fmt.Printf("# Found trueToSize: %v\n", trueToSize)
			return Shoe{Name: shoeName, TrueToSize: trueToSize}
		}
		fmt.Println("# Shoe not found") //the scanner encountered an incompatible type (nil) and generated an error
	}

	return Shoe{} //default to an empty Shoe
}

/*
given a Shoe, insert the record into the database
pass back any errors (constraint violations)
*/
func DBCreateShoe(shoe Shoe) (Shoe, error) {
	shoe.TrueToSize = float32(int(shoe.TrueToSize)) //truncate any non-int values for TrueToSize
	fmt.Println("# Inserting", shoe.Name, shoe.TrueToSize)

	//even though the struct is floating-point (used as a dual-purpose for incoming shoes and outgoing results),
	//the database is expecting an integer for TrueToSize
	_, err := db.Exec("INSERT INTO shoeinfo(shoe,trueToSize) VALUES($1,$2);", shoe.Name, int(shoe.TrueToSize))
	return shoe, err
}
