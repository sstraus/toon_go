package main

import (
	"fmt"

	"github.com/sstraus/toon_go/toon"
)

func main() {
	fmt.Println("=== Advanced TOON Encoding Examples ===\n")

	// Example 1: Inline Array Format (primitives)
	fmt.Println("1. Inline Array Format (primitives):")
	inlineData := map[string]interface{}{
		"languages": []interface{}{"Go", "Python", "JavaScript", "Rust"},
		"ports":     []interface{}{8080, 8081, 8082},
		"flags":     []interface{}{true, false, true, true},
	}
	result, _ := toon.MarshalToString(inlineData)
	fmt.Println(result)

	// Example 2: Tabular Array Format (uniform objects)
	fmt.Println("\n2. Tabular Array Format (uniform objects):")
	tabularData := map[string]interface{}{
		"employees": []interface{}{
			map[string]interface{}{"id": 101, "name": "Alice Johnson", "department": "Engineering", "salary": 95000},
			map[string]interface{}{"id": 102, "name": "Bob Smith", "department": "Marketing", "salary": 75000},
			map[string]interface{}{"id": 103, "name": "Carol Williams", "department": "Engineering", "salary": 105000},
			map[string]interface{}{"id": 104, "name": "David Brown", "department": "Sales", "salary": 80000},
		},
	}
	result, _ = toon.MarshalToString(tabularData)
	fmt.Println(result)

	// Example 3: List Array Format (mixed/nested objects)
	fmt.Println("\n3. List Array Format (mixed types):")
	listData := map[string]interface{}{
		"resources": []interface{}{
			map[string]interface{}{
				"type":     "database",
				"name":     "postgres-prod",
				"host":     "db.example.com",
				"port":     5432,
				"replicas": []interface{}{"db-replica-1", "db-replica-2"},
			},
			map[string]interface{}{
				"type":   "cache",
				"name":   "redis-cache",
				"host":   "cache.example.com",
				"port":   6379,
				"memory": "4GB",
			},
			map[string]interface{}{
				"type":      "queue",
				"name":      "rabbitmq",
				"exchanges": []interface{}{"events", "tasks", "notifications"},
			},
		},
	}
	result, _ = toon.MarshalToString(listData)
	fmt.Println(result)

	// Example 4: Custom Delimiter (Tab-separated)
	fmt.Println("\n4. Custom Delimiter (tab-separated values):")
	tsvData := map[string]interface{}{
		"measurements": []interface{}{
			map[string]interface{}{"sensor": "temp-01", "value": 22.5, "unit": "celsius"},
			map[string]interface{}{"sensor": "temp-02", "value": 23.1, "unit": "celsius"},
			map[string]interface{}{"sensor": "temp-03", "value": 21.8, "unit": "celsius"},
		},
	}
	result, _ = toon.MarshalToString(tsvData, toon.WithDelimiter("\t"))
	fmt.Println(result)

	// Example 5: Custom Indentation (4 spaces)
	fmt.Println("\n5. Custom Indentation (4 spaces):")
	deepNested := map[string]interface{}{
		"application": map[string]interface{}{
			"server": map[string]interface{}{
				"host": "localhost",
				"port": 8080,
				"ssl": map[string]interface{}{
					"enabled": true,
					"cert":    "/path/to/cert.pem",
					"key":     "/path/to/key.pem",
				},
			},
		},
	}
	result, _ = toon.MarshalToString(deepNested, toon.WithIndent(4))
	fmt.Println(result)

	// Example 6: Length Markers with Custom Prefix
	fmt.Println("\n6. Length Markers with Custom Prefix:")
	markedData := map[string]interface{}{
		"versions": []interface{}{"v1.0.0", "v1.1.0", "v1.2.0", "v2.0.0"},
		"features": []interface{}{
			map[string]interface{}{"name": "auth", "status": "stable"},
			map[string]interface{}{"name": "api", "status": "beta"},
			map[string]interface{}{"name": "ui", "status": "alpha"},
		},
	}
	result, _ = toon.MarshalToString(markedData, toon.WithLengthMarker("#"))
	fmt.Println(result)

	// Example 7: Complex Real-World Structure
	fmt.Println("\n7. Complex Structure (microservices architecture):")
	architecture := map[string]interface{}{
		"platform": "Kubernetes",
		"version":  "1.28.0",
		"services": []interface{}{
			map[string]interface{}{
				"name":     "api-gateway",
				"image":    "api-gateway:2.1.0",
				"replicas": 3,
				"ports":    []interface{}{8080, 8443},
				"env": map[string]interface{}{
					"LOG_LEVEL":     "info",
					"RATE_LIMIT":    "1000",
					"TIMEOUT":       "30s",
					"CORS_ORIGINS":  []interface{}{"https://app.example.com", "https://admin.example.com"},
				},
			},
			map[string]interface{}{
				"name":     "user-service",
				"image":    "user-service:1.5.2",
				"replicas": 2,
				"database": map[string]interface{}{
					"type": "postgres",
					"host": "users-db",
					"port": 5432,
				},
			},
		},
		"monitoring": map[string]interface{}{
			"prometheus": true,
			"grafana":    true,
			"alerts": []interface{}{
				"high_cpu_usage",
				"memory_threshold",
				"error_rate_spike",
			},
		},
	}
	result, _ = toon.MarshalToString(architecture, toon.WithIndent(2))
	fmt.Println(result)

	// Example 8: Combining Multiple Options
	fmt.Println("\n8. Multiple Options Combined:")
	combined := map[string]interface{}{
		"metrics": []interface{}{1.23, 4.56, 7.89, 10.11},
		"labels":  []interface{}{"alpha", "beta", "gamma", "delta"},
	}
	result, _ = toon.MarshalToString(combined,
		toon.WithDelimiter("|"),
		toon.WithLengthMarker("COUNT:"),
		toon.WithIndent(3),
	)
	fmt.Println(result)
}
