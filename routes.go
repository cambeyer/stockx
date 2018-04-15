package main

import "net/http"

type Route struct {
	Name        string
	Method      string //GET, POST
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

//add any new custom routes here
var routes = Routes{
	Route{
		"ShoeCreate",
		"POST",
		"/shoes",
		ShoeCreate,
	},
	Route{
		"ShoeShow",
		"GET",
		"/shoes/{shoeName}",
		ShoeShow,
	},
}
