package main

import (
	"fmt"

	"github.com/muccarini/jsonparser"
)

func main() {
	json := []byte(`{
		"user": {
			"id": 12345,
			"name": "John Doe",
			"isActive": true,
			"location": {
				"city": "SF",
				"coords": [37.7749, -122.4194]
			}
		},
		"items": [
			{"id": 1, "tags": ["go", "test"]},
			{"id": 2, "featured": false}
		],
		"categories": ["Tech", "Programming"],
		"emptyArray": [],
		"emptyObject": {},
		"nullValue": null,
		"booleanTrue": true,
		"booleanFalse": false,
		"largeNumber": 9223372036854775807,
		"floatingPoint": 3.14159,
		"scientificNotation": 1.23e-4,
		"unicodeString": "Hello 世界 🌍",
		"escapedChars": "Quote: \"Hi\", Tab:\t, Newline:\n, Backslash: \\",
		"specialChars": "Forward slash: /"
	}`)

	res, err := jsonparser.Get[string](json, "user", "name")

	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	fmt.Printf("RESPONSE: %v", res)
}
