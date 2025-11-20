package main

import (
	"fmt"
	"log"

	"github.com/sstraus/toon_go/toon"
)

func main() {
	fmt.Println("=== TOON Go Library Examples ===\n")

	// Example 1: Simple primitives
	fmt.Println("1. Simple Primitives:")
	examples := []interface{}{
		nil,
		true,
		42,
		3.14,
		"hello",
	}
	for _, ex := range examples {
		encoded, _ := toon.MarshalToString(ex)
		fmt.Printf("  %v -> %s\n", ex, encoded)
	}

	// Example 2: Simple object
	fmt.Println("\n2. Simple Object:")
	obj := map[string]interface{}{
		"name": "Alice",
		"age":  30,
		"city": "New York",
	}
	encoded, err := toon.MarshalToString(obj)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(encoded)

	// Example 3: Nested object
	fmt.Println("\n3. Nested Object:")
	nested := map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "Bob",
			"email": "bob@example.com",
		},
	}
	encoded, _ = toon.MarshalToString(nested)
	fmt.Println(encoded)

	// Example 4: Inline array
	fmt.Println("\n4. Inline Array (primitives):")
	inlineArr := map[string]interface{}{
		"tags": []interface{}{"go", "toon", "llm"},
	}
	encoded, _ = toon.MarshalToString(inlineArr)
	fmt.Println(encoded)

	// Example 5: Tabular array
	fmt.Println("\n5. Tabular Array (uniform objects):")
	tabular := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "age": 30},
			map[string]interface{}{"name": "Bob", "age": 25},
			map[string]interface{}{"name": "Carol", "age": 28},
		},
	}
	encoded, _ = toon.MarshalToString(tabular)
	fmt.Println(encoded)

	// Example 6: List array (mixed)
	fmt.Println("\n6. List Array (mixed types):")
	listArr := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"type": "book", "title": "Go Programming"},
			map[string]interface{}{"type": "video", "duration": 120},
		},
	}
	encoded, _ = toon.MarshalToString(listArr)
	fmt.Println(encoded)

	// Example 7: Custom options (using functional options)
	fmt.Println("\n7. Custom Options (tab delimiter, length marker):")
	withOpts := map[string]interface{}{
		"values": []interface{}{1, 2, 3, 4, 5},
	}
	encoded, _ = toon.MarshalToString(withOpts, toon.WithDelimiter("\t"), toon.WithLengthMarker("#"))
	fmt.Println(encoded)

	// Example 8: Complex nested structure
	fmt.Println("\n8. Complex Nested Structure:")
	complex := map[string]interface{}{
		"project": map[string]interface{}{
			"name":    "TOON Go",
			"version": "1.1.0",
			"authors": []interface{}{
				map[string]interface{}{"name": "Alice", "role": "Lead"},
				map[string]interface{}{"name": "Bob", "role": "Developer"},
			},
			"tags": []interface{}{"go", "toon", "encoding"},
		},
	}
	encoded, _ = toon.MarshalToString(complex)
	fmt.Println(encoded)

	// Example 9: Token efficiency comparison
	fmt.Println("\n9. Token Efficiency (TOON vs JSON):")
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "age": 30},
			map[string]interface{}{"name": "Bob", "age": 25},
		},
	}
	toonEncoded, _ := toon.MarshalToString(data)
	fmt.Printf("TOON (%d chars):\n%s\n", len(toonEncoded), toonEncoded)

	// Note: JSON comparison would require encoding/json import
	fmt.Println("\nJSON equivalent would be approximately 30-60% longer!")
}
