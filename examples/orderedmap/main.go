package main

import (
	"fmt"
	"log"

	"github.com/sstraus/toon_go/toon"
)

func main() {
	fmt.Println("=== OrderedMap Examples ===\n")

	// Example 1: Basic OrderedMap Usage
	fmt.Println("1. Basic OrderedMap (preserves insertion order):")
	om := toon.NewOrderedMap()
	om.Set("zulu", "last alphabetically")
	om.Set("alpha", "first alphabetically")
	om.Set("mike", "middle alphabetically")

	result, _ := toon.MarshalToString(om)
	fmt.Println(result)
	fmt.Println("  ^ Keys appear in insertion order, not alphabetical")

	// Example 2: Configuration with Logical Order
	fmt.Println("\n2. Configuration File (logical section ordering):")
	config := toon.NewOrderedMap()

	// Set in logical reading order
	config.Set("version", "1.0.0")
	config.Set("name", "Production API")

	// Server section
	server := toon.NewOrderedMap()
	server.Set("host", "0.0.0.0")
	server.Set("port", 8080)
	server.Set("timeout", 30)
	config.Set("server", server)

	// Database section
	database := toon.NewOrderedMap()
	database.Set("driver", "postgresql")
	database.Set("host", "localhost")
	database.Set("port", 5432)
	database.Set("database", "production")
	config.Set("database", database)

	// Logging section
	logging := toon.NewOrderedMap()
	logging.Set("level", "info")
	logging.Set("format", "json")
	logging.Set("output", "/var/log/app.log")
	config.Set("logging", logging)

	result, _ = toon.MarshalToString(config)
	fmt.Println(result)

	// Example 3: HTTP Headers (Order Matters)
	fmt.Println("\n3. HTTP Headers (order preserved):")
	headers := toon.NewOrderedMap()
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", "Bearer token123")
	headers.Set("Accept", "application/json")
	headers.Set("User-Agent", "TOON-Client/1.0")
	headers.Set("X-Request-ID", "abc-123-def-456")

	result, _ = toon.MarshalToString(headers)
	fmt.Println(result)

	// Example 4: Build Steps (Sequential Order)
	fmt.Println("\n4. Build Pipeline Steps (sequential order):")
	pipeline := toon.NewOrderedMap()
	pipeline.Set("name", "CI/CD Pipeline")

	steps := []interface{}{}

	step1 := toon.NewOrderedMap()
	step1.Set("step", 1)
	step1.Set("name", "checkout")
	step1.Set("command", "git clone")
	steps = append(steps, step1)

	step2 := toon.NewOrderedMap()
	step2.Set("step", 2)
	step2.Set("name", "dependencies")
	step2.Set("command", "go mod download")
	steps = append(steps, step2)

	step3 := toon.NewOrderedMap()
	step3.Set("step", 3)
	step3.Set("name", "test")
	step3.Set("command", "go test ./...")
	steps = append(steps, step3)

	step4 := toon.NewOrderedMap()
	step4.Set("step", 4)
	step4.Set("name", "build")
	step4.Set("command", "go build -o app")
	steps = append(steps, step4)

	step5 := toon.NewOrderedMap()
	step5.Set("step", 5)
	step5.Set("name", "deploy")
	step5.Set("command", "kubectl apply -f deployment.yaml")
	steps = append(steps, step5)

	pipeline.Set("steps", steps)

	result, _ = toon.MarshalToString(pipeline)
	fmt.Println(result)

	// Example 5: OrderedMap Methods
	fmt.Println("\n5. OrderedMap Methods:")
	demo := toon.NewOrderedMap()
	demo.Set("first", 1)
	demo.Set("second", 2)
	demo.Set("third", 3)

	fmt.Println("  Keys:", demo.Keys())
	fmt.Println("  Length:", demo.Len())

	if val, ok := demo.Get("second"); ok {
		fmt.Printf("  Get('second'): %v\n", val)
	}

	demo.Delete("second")
	fmt.Println("  After Delete('second'):", demo.Keys())

	// Example 6: Nested OrderedMaps
	fmt.Println("\n6. Nested OrderedMaps (preserving structure):")
	document := toon.NewOrderedMap()
	document.Set("title", "User Documentation")
	document.Set("version", "2.0")

	// Table of contents in specific order
	toc := toon.NewOrderedMap()
	toc.Set("1_introduction", "Getting Started")
	toc.Set("2_installation", "Installation Guide")
	toc.Set("3_configuration", "Configuration Options")
	toc.Set("4_advanced", "Advanced Topics")
	toc.Set("5_troubleshooting", "Troubleshooting")
	document.Set("table_of_contents", toc)

	// Metadata
	metadata := toon.NewOrderedMap()
	metadata.Set("author", "Documentation Team")
	metadata.Set("date", "2025-11-20")
	metadata.Set("status", "published")
	document.Set("metadata", metadata)

	result, _ = toon.MarshalToString(document)
	fmt.Println(result)

	// Example 7: Comparison with Regular Map
	fmt.Println("\n7. OrderedMap vs Regular Map:")

	// Regular map (alphabetically sorted)
	regularMap := map[string]interface{}{
		"zebra":    "z",
		"apple":    "a",
		"mango":    "m",
		"banana":   "b",
	}
	regular, _ := toon.MarshalToString(regularMap)
	fmt.Println("  Regular map (alphabetical):")
	fmt.Println(regular)

	// OrderedMap (insertion order)
	orderedMap := toon.NewOrderedMap()
	orderedMap.Set("zebra", "z")
	orderedMap.Set("apple", "a")
	orderedMap.Set("mango", "m")
	orderedMap.Set("banana", "b")
	ordered, _ := toon.MarshalToString(orderedMap)
	fmt.Println("\n  OrderedMap (insertion order):")
	fmt.Println(ordered)

	// Example 8: Form Fields (Display Order)
	fmt.Println("\n8. Form Fields (logical display order):")
	form := toon.NewOrderedMap()
	form.Set("form_id", "user_registration")

	fields := []interface{}{}

	f1 := toon.NewOrderedMap()
	f1.Set("name", "full_name")
	f1.Set("label", "Full Name")
	f1.Set("type", "text")
	f1.Set("required", true)
	fields = append(fields, f1)

	f2 := toon.NewOrderedMap()
	f2.Set("name", "email")
	f2.Set("label", "Email Address")
	f2.Set("type", "email")
	f2.Set("required", true)
	fields = append(fields, f2)

	f3 := toon.NewOrderedMap()
	f3.Set("name", "phone")
	f3.Set("label", "Phone Number")
	f3.Set("type", "tel")
	f3.Set("required", false)
	fields = append(fields, f3)

	f4 := toon.NewOrderedMap()
	f4.Set("name", "password")
	f4.Set("label", "Password")
	f4.Set("type", "password")
	f4.Set("required", true)
	fields = append(fields, f4)

	form.Set("fields", fields)

	result, err := toon.MarshalToString(form, toon.WithIndent(2))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
