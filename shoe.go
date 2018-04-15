package main

//JSON keys should not be capitalized
//more information could be added to this struct in the future
type Shoe struct {
	Name       string  `json:"name"`
	TrueToSize float32 `json:"trueToSize"`
}
