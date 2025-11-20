package main

import (
	"fmt"

	"github.com/sstraus/toon_go/toon"
)

func main() {
	fmt.Println("=== API Response Examples (TOON for LLM-friendly APIs) ===\n")

	// Example 1: User Profile API Response
	fmt.Println("1. User Profile API Response:")
	userProfile := toon.NewOrderedMap()
	userProfile.Set("user_id", "usr_7x9k2p")
	userProfile.Set("username", "alice_developer")
	userProfile.Set("email", "alice@example.com")
	userProfile.Set("full_name", "Alice Johnson")
	userProfile.Set("role", "senior_engineer")
	userProfile.Set("active", true)
	userProfile.Set("created_at", "2023-01-15T10:30:00Z")
	userProfile.Set("last_login", "2025-11-20T09:15:00Z")

	profile := toon.NewOrderedMap()
	profile.Set("bio", "Full-stack developer passionate about Go and distributed systems")
	profile.Set("location", "San Francisco, CA")
	profile.Set("website", "https://alice.dev")
	profile.Set("skills", []interface{}{"Go", "Kubernetes", "PostgreSQL", "React", "gRPC"})
	userProfile.Set("profile", profile)

	stats := toon.NewOrderedMap()
	stats.Set("repos", 47)
	stats.Set("followers", 1523)
	stats.Set("following", 89)
	userProfile.Set("stats", stats)

	response, _ := toon.MarshalToString(userProfile)
	fmt.Println(response)

	// Example 2: List API Response with Pagination
	fmt.Println("\n2. List API Response with Pagination:")
	listResponse := toon.NewOrderedMap()
	listResponse.Set("status", "success")
	listResponse.Set("total_count", 247)
	listResponse.Set("page", 1)
	listResponse.Set("page_size", 20)
	listResponse.Set("total_pages", 13)

	items := []interface{}{}
	for i := 1; i <= 5; i++ {
		item := toon.NewOrderedMap()
		item.Set("id", fmt.Sprintf("prod_%d", 1000+i))
		item.Set("name", fmt.Sprintf("Product %d", i))
		item.Set("price", 99.99+float64(i)*10)
		item.Set("in_stock", i%2 == 0)
		item.Set("category", "electronics")
		items = append(items, item)
	}
	listResponse.Set("items", items)

	links := toon.NewOrderedMap()
	links.Set("self", "/api/v1/products?page=1")
	links.Set("next", "/api/v1/products?page=2")
	links.Set("last", "/api/v1/products?page=13")
	listResponse.Set("links", links)

	listResp, _ := toon.MarshalToString(listResponse)
	fmt.Println(listResp)

	// Example 3: Error API Response
	fmt.Println("\n3. Error API Response:")
	errorResponse := toon.NewOrderedMap()
	errorResponse.Set("status", "error")
	errorResponse.Set("error_code", "VALIDATION_ERROR")
	errorResponse.Set("message", "Invalid request parameters")
	errorResponse.Set("timestamp", "2025-11-20T14:23:00Z")
	errorResponse.Set("request_id", "req_abc123def456")

	errors := []interface{}{}

	err1 := toon.NewOrderedMap()
	err1.Set("field", "email")
	err1.Set("message", "Invalid email format")
	err1.Set("code", "INVALID_FORMAT")
	errors = append(errors, err1)

	err2 := toon.NewOrderedMap()
	err2.Set("field", "age")
	err2.Set("message", "Must be at least 18")
	err2.Set("code", "MIN_VALUE_REQUIRED")
	errors = append(errors, err2)

	errorResponse.Set("errors", errors)
	errorResponse.Set("documentation_url", "https://docs.example.com/errors/validation")

	errResp, _ := toon.MarshalToString(errorResponse)
	fmt.Println(errResp)

	// Example 4: Search API Response
	fmt.Println("\n4. Search API Response:")
	searchResponse := toon.NewOrderedMap()
	searchResponse.Set("query", "golang microservices")
	searchResponse.Set("results_found", 3)
	searchResponse.Set("search_time_ms", 45)

	results := []interface{}{}

	result1 := toon.NewOrderedMap()
	result1.Set("id", "doc_1")
	result1.Set("title", "Building Microservices with Go")
	result1.Set("excerpt", "Learn how to build scalable microservices using Go...")
	result1.Set("relevance_score", 0.95)
	result1.Set("author", "John Doe")
	result1.Set("tags", []interface{}{"go", "microservices", "architecture"})
	results = append(results, result1)

	result2 := toon.NewOrderedMap()
	result2.Set("id", "doc_2")
	result2.Set("title", "Go Concurrency Patterns for Distributed Systems")
	result2.Set("excerpt", "Advanced patterns for handling concurrency...")
	result2.Set("relevance_score", 0.87)
	result2.Set("author", "Jane Smith")
	result2.Set("tags", []interface{}{"go", "concurrency", "distributed-systems"})
	results = append(results, result2)

	result3 := toon.NewOrderedMap()
	result3.Set("id", "doc_3")
	result3.Set("title", "Deploying Go Services to Kubernetes")
	result3.Set("excerpt", "Step-by-step guide to containerization and deployment...")
	result3.Set("relevance_score", 0.82)
	result3.Set("author", "Bob Johnson")
	result3.Set("tags", []interface{}{"go", "kubernetes", "deployment"})
	results = append(results, result3)

	searchResponse.Set("results", results)

	searchResp, _ := toon.MarshalToString(searchResponse)
	fmt.Println(searchResp)

	// Example 5: Analytics Dashboard API Response
	fmt.Println("\n5. Analytics Dashboard API Response:")
	analytics := toon.NewOrderedMap()
	analytics.Set("dashboard", "website_traffic")
	analytics.Set("period", "last_7_days")
	analytics.Set("generated_at", "2025-11-20T15:00:00Z")

	summary := toon.NewOrderedMap()
	summary.Set("total_visits", 125430)
	summary.Set("unique_visitors", 45872)
	summary.Set("page_views", 312456)
	summary.Set("avg_session_duration_seconds", 245)
	summary.Set("bounce_rate_percent", 42.3)
	analytics.Set("summary", summary)

	traffic := []interface{}{}
	days := []string{"2025-11-14", "2025-11-15", "2025-11-16", "2025-11-17", "2025-11-18", "2025-11-19", "2025-11-20"}
	visits := []int{15234, 16892, 18456, 20123, 19876, 17654, 17195}

	for i, day := range days {
		dayData := toon.NewOrderedMap()
		dayData.Set("date", day)
		dayData.Set("visits", visits[i])
		dayData.Set("unique", visits[i]*70/100)
		traffic = append(traffic, dayData)
	}
	analytics.Set("daily_traffic", traffic)

	topPages := []interface{}{
		map[string]interface{}{"path": "/", "views": 45234},
		map[string]interface{}{"path": "/products", "views": 38901},
		map[string]interface{}{"path": "/blog", "views": 28456},
		map[string]interface{}{"path": "/about", "views": 15234},
		map[string]interface{}{"path": "/contact", "views": 12345},
	}
	analytics.Set("top_pages", topPages)

	sources := map[string]interface{}{
		"direct":       35.2,
		"google":       28.7,
		"social_media": 18.5,
		"referral":     12.1,
		"email":        5.5,
	}
	analytics.Set("traffic_sources_percent", sources)

	analyticsResp, _ := toon.MarshalToString(analytics)
	fmt.Println(analyticsResp)

	// Example 6: Webhook Payload
	fmt.Println("\n6. Webhook Payload:")
	webhook := toon.NewOrderedMap()
	webhook.Set("event_id", "evt_9x7k2m")
	webhook.Set("event_type", "order.completed")
	webhook.Set("timestamp", "2025-11-20T16:30:00Z")
	webhook.Set("api_version", "2025-11")

	payload := toon.NewOrderedMap()
	payload.Set("order_id", "ord_abc123")
	payload.Set("customer_id", "cust_xyz789")
	payload.Set("status", "completed")
	payload.Set("total_amount", 299.97)
	payload.Set("currency", "USD")

	orderItems := []interface{}{}
	item1 := map[string]interface{}{"sku": "PROD-001", "quantity": 2, "price": 99.99}
	item2 := map[string]interface{}{"sku": "PROD-002", "quantity": 1, "price": 99.99}
	orderItems = append(orderItems, item1, item2)
	payload.Set("items", orderItems)

	shipping := toon.NewOrderedMap()
	shipping.Set("method", "express")
	shipping.Set("address", "123 Main St, San Francisco, CA 94105")
	shipping.Set("tracking_number", "TRK-1234567890")
	payload.Set("shipping", shipping)

	webhook.Set("data", payload)

	webhookResp, _ := toon.MarshalToString(webhook)
	fmt.Println(webhookResp)

	// Example 7: GraphQL-style Response
	fmt.Println("\n7. GraphQL-style Nested Response:")
	gqlResponse := toon.NewOrderedMap()

	organization := toon.NewOrderedMap()
	organization.Set("name", "TechCorp Inc")
	organization.Set("industry", "Software")

	teams := []interface{}{}

	engTeam := toon.NewOrderedMap()
	engTeam.Set("name", "Engineering")
	engTeam.Set("size", 45)

	engMembers := []interface{}{}
	for i := 1; i <= 3; i++ {
		member := toon.NewOrderedMap()
		member.Set("id", fmt.Sprintf("emp_%d", i))
		member.Set("name", fmt.Sprintf("Engineer %d", i))
		member.Set("role", "Software Engineer")
		member.Set("projects", []interface{}{fmt.Sprintf("Project %d", i), fmt.Sprintf("Project %d", i+1)})
		engMembers = append(engMembers, member)
	}
	engTeam.Set("members", engMembers)
	teams = append(teams, engTeam)

	organization.Set("teams", teams)
	gqlResponse.Set("data", map[string]interface{}{"organization": organization})

	gqlResp, _ := toon.MarshalToString(gqlResponse)
	fmt.Println(gqlResp)

	fmt.Println("\n✓ TOON provides 30-60% token savings for LLM-based API interactions")
	fmt.Println("✓ More readable than JSON while maintaining full type support")
	fmt.Println("✓ Perfect for AI agents and LLM-powered applications")
}
