package main

import (
	"fmt"
	"log"

	"github.com/sstraus/toon_go/toon"
)

func main() {
	fmt.Println("=== TOON Decoding Examples ===")

	// Example 1: Simple Object Decoding
	fmt.Println("1. Decoding Simple Object:")
	input1 := `name: Production Server
host: prod.example.com
port: 443
ssl_enabled: true`

	var config map[string]interface{}
	err := toon.UnmarshalFromString(input1, &config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Decoded: %+v\n", config)

	// Example 2: Nested Object Decoding
	fmt.Println("\n2. Decoding Nested Object:")
	input2 := `database:
  driver: postgresql
  host: db.example.com
  port: 5432
  credentials:
    username: admin
    password_encrypted: true`

	var dbConfig map[string]interface{}
	err = toon.UnmarshalFromString(input2, &dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Database config: %+v\n", dbConfig)

	// Example 3: Inline Array Decoding
	fmt.Println("\n3. Decoding Inline Array:")
	input3 := `allowed_origins[3]: https://app.example.com,https://admin.example.com,https://api.example.com`

	var origins map[string]interface{}
	err = toon.UnmarshalFromString(input3, &origins)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Origins: %+v\n", origins)

	// Example 4: Tabular Array Decoding
	fmt.Println("\n4. Decoding Tabular Array:")
	input4 := `users[3]{id,username,role}:
  1,alice,admin
  2,bob,user
  3,carol,moderator`

	var users map[string]interface{}
	err = toon.UnmarshalFromString(input4, &users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Users: %+v\n", users)

	// Example 5: List Array Decoding
	fmt.Println("\n5. Decoding List Array:")
	input5 := `endpoints[2]:
  - path: /api/v1/users
    method: GET
    auth: required
  - path: /api/v1/health
    method: GET
    auth: none`

	var api map[string]interface{}
	err = toon.UnmarshalFromString(input5, &api)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  API endpoints: %+v\n", api)

	// Example 6: Strict Mode Disabled
	fmt.Println("\n6. Decoding with Strict Mode Disabled:")
	input6 := `name: Test Config
value: 123


extra_newlines: allowed`  // Extra blank lines

	var flexible map[string]interface{}
	err = toon.UnmarshalFromString(input6, &flexible, toon.WithStrictDecoding(false))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Flexible decode: %+v\n", flexible)

	// Example 7: Complex Nested Structure
	fmt.Println("\n7. Decoding Complex Structure:")
	input7 := `application:
  name: MyApp
  version: 2.0.0
  features[3]: auth,api,analytics
  environments[2]:
    - name: development
      url: https://dev.example.com
      debug: true
    - name: production
      url: https://prod.example.com
      debug: false
  security:
    jwt_enabled: true
    session_timeout: 3600
    allowed_ips[2]: 10.0.0.0/8,192.168.0.0/16`

	var app map[string]interface{}
	err = toon.UnmarshalFromString(input7, &app)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Application config: %+v\n", app)

	// Example 8: Round-trip Encoding/Decoding
	fmt.Println("\n8. Round-trip Encoding/Decoding:")
	original := map[string]interface{}{
		"service": "payment-processor",
		"version": "3.2.1",
		"settings": map[string]interface{}{
			"timeout":       30,
			"retry_count":   3,
			"fallback_mode": true,
		},
		"providers": []interface{}{"stripe", "paypal", "square"},
	}

	// Encode
	encoded, err := toon.MarshalToString(original)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Encoded:")
	fmt.Println(encoded)

	// Decode back
	var decoded map[string]interface{}
	err = toon.UnmarshalFromString(encoded, &decoded)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n  Decoded back: %+v\n", decoded)

	// Example 9: Handling Different Data Types
	fmt.Println("\n9. Decoding Various Data Types:")
	input9 := `int_value: 42
float_value: 3.14159
bool_true: true
bool_false: false
null_value: null
string_value: Hello TOON
array_mixed[4]: 1,true,hello,null`

	var types map[string]interface{}
	err = toon.UnmarshalFromString(input9, &types)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Decoded types:")
	for k, v := range types {
		fmt.Printf("    %s: %v (type: %T)\n", k, v, v)
	}

	// Example 10: Decoding to interface{}
	fmt.Println("\n10. Decoding to Generic Interface:")
	input10 := `status: active
count: 100
items[2]: item1,item2`

	var generic interface{}
	err = toon.UnmarshalFromString(input10, &generic)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Generic decode: %+v\n", generic)
}
