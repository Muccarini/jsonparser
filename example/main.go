package main

import (
	"fmt"
	"log"

	"github.com/Muccarini/jsonparser"
)

func main() {
	// Sample JSON data
	jsonData := []byte(`{
		"user": {
			"id": 1,
			"name": "Alice Johnson",
			"email": "alice@example.com",
			"address": {
				"street": "123 Main St",
				"city": "Boston",
				"coordinates": {
					"lat": 42.3601,
					"lng": -71.0589
				}
			},
			"active": true,
			"roles": ["admin", "user"],
			"metadata": null
		},
		"products": [
			{
				"id": 101,
				"name": "Laptop",
				"price": 999.99,
				"available": true
			},
			{
				"id": 102,
				"name": "Mouse",
				"price": 29.99,
				"available": false
			}
		],
		"version": "1.0.0"
	}`)

	fmt.Println("=== Basic Value Extraction ===")
	
	// Extract basic values
	name, err := jsonparser.GetString(jsonData, "user", "name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User name: %s\n", name)

	userId, err := jsonparser.GetInt(jsonData, "user", "id")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User ID: %d\n", userId)

	lat, err := jsonparser.GetFloat(jsonData, "user", "address", "coordinates", "lat")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Latitude: %f\n", lat)

	active, err := jsonparser.GetBoolean(jsonData, "user", "active")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User active: %t\n", active)

	fmt.Println("\n=== Array Navigation ===")
	
	// Access array elements
	firstRole, err := jsonparser.GetString(jsonData, "user", "roles", "0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("First role: %s\n", firstRole)

	firstProductPrice, err := jsonparser.GetFloat(jsonData, "products", "0", "price")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("First product price: $%.2f\n", firstProductPrice)

	fmt.Println("\n=== Raw Value Access ===")
	
	// Get raw values with type information
	rawValue, valueType, _, err := jsonparser.Get(jsonData, "user", "metadata")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Metadata: %s (type: %v)\n", string(rawValue), valueType)

	fmt.Println("\n=== Array Iteration ===")
	
	// Iterate over arrays
	fmt.Println("User roles:")
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		if dataType == jsonparser.String {
			// For string values, we need to parse them
			role, parseErr := jsonparser.ParseStringValue(value)
			if parseErr != nil {
				log.Printf("Error parsing string: %v", parseErr)
				return
			}
			fmt.Printf("  - %s\n", role)
		}
	}, "user", "roles")

	fmt.Println("\nProducts:")
	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		if dataType == jsonparser.Object {
			// Extract values from each product object
			name, _ := jsonparser.GetString(value, "name")
			price, _ := jsonparser.GetFloat(value, "price")
			available, _ := jsonparser.GetBoolean(value, "available")
			fmt.Printf("  - %s: $%.2f (available: %t)\n", name, price, available)
		}
	}, "products")

	fmt.Println("\n=== Object Iteration ===")
	
	// Iterate over object keys
	fmt.Println("User address fields:")
	jsonparser.ObjectEach(jsonData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		fmt.Printf("  %s: %s\n", string(key), string(value))
		return nil
	}, "user", "address")

	fmt.Println("\n=== Error Handling ===")
	
	// Demonstrate error handling
	_, err = jsonparser.GetString(jsonData, "nonexistent", "key")
	if err != nil {
		fmt.Printf("Expected error for missing key: %v\n", err)
	}

	_, err = jsonparser.GetString(jsonData, "user", "roles", "10") // Out of bounds
	if err != nil {
		fmt.Printf("Expected error for out of bounds: %v\n", err)
	}

	fmt.Println("\n=== Performance Features ===")
	fmt.Println("✓ Zero memory allocation for Get operations")
	fmt.Println("✓ Direct byte slice access without copying")
	fmt.Println("✓ Streaming-friendly design")
	fmt.Println("✓ JSONPath-like key navigation")
}