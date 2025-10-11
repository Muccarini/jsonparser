package main

import (
	"fmt"

	"github.com/muccarini/jsonparser"
)

func main() {
	json := []byte(`{
        "name": "John",
        "age": 30,
        "active": true,
        "score": 85.5,
        "tags": ["dev", "golang"],
        "profile": {"city": "NYC"}
    }`)

	// Basic extraction
	name, _ := jsonparser.GetString(json, "name")
	age, _ := jsonparser.GetInt(json, "age")
	active, _ := jsonparser.GetBool(json, "active")
	score, _ := jsonparser.GetFloat64(json, "score")

	// Basic extraction using generic parser
	var nameGen string
	jsonparser.Get(&nameGen, json, "name")

	var ageGen int
	jsonparser.Get(&ageGen, json, "age")

	var activeGen bool
	jsonparser.Get(&activeGen, json, "active")

	var scoreGen float64
	jsonparser.Get(&scoreGen, json, "score")

	fmt.Printf("Generic - Name: %s, Age: %d, Active: %t, Score: %.1f\n", nameGen, ageGen, activeGen, scoreGen)

	// Array and nested access
	tag, _ := jsonparser.GetString(json, "tags", "1")
	city, _ := jsonparser.GetString(json, "profile", "city")

	fmt.Printf("Name: %s, Age: %d, Active: %t\n", name, age, active)
	fmt.Printf("Score: %.1f, Tag: %s, City: %s\n", score, tag, city)
}
