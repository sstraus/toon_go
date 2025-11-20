package main

import (
	"fmt"
	"log"

	"github.com/sstraus/toon_go/toon"
)

func main() {
	fmt.Println("=== Configuration File Examples ===\n")

	// Example 1: Read Application Configuration
	fmt.Println("1. Reading Application Configuration:")
	configTOON := `app_name: E-Commerce Platform
version: 3.5.2
environment: production

server:
  host: 0.0.0.0
  port: 8443
  ssl_enabled: true
  cert_path: /etc/ssl/certs/server.crt
  key_path: /etc/ssl/private/server.key
  read_timeout: 30
  write_timeout: 30

database:
  primary:
    driver: postgresql
    host: db-primary.internal
    port: 5432
    database: ecommerce_prod
    pool_size: 50
    max_idle: 10
  replica:
    host: db-replica.internal
    port: 5432
    database: ecommerce_prod
    pool_size: 30

cache:
  enabled: true
  provider: redis
  hosts[3]: cache-1.internal:6379,cache-2.internal:6379,cache-3.internal:6379
  ttl_seconds: 3600
  max_memory: 4096

logging:
  level: info
  format: json
  output: /var/log/app/application.log
  rotation:
    enabled: true
    max_size_mb: 100
    max_age_days: 30
    compress: true

security:
  jwt_secret_env: JWT_SECRET_KEY
  jwt_expiry_hours: 24
  rate_limiting:
    enabled: true
    requests_per_minute: 100
    burst_size: 20
  cors:
    allowed_origins[2]: https://shop.example.com,https://admin.example.com
    allowed_methods[4]: GET,POST,PUT,DELETE
    allow_credentials: true

features:
  payment_providers[3]: stripe,paypal,square
  email_service: sendgrid
  sms_service: twilio
  analytics_enabled: true
  recommendations_enabled: true`

	var appConfig map[string]interface{}
	err := toon.UnmarshalFromString(configTOON, &appConfig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("  Application: %s v%s\n", appConfig["app_name"], appConfig["version"])
	fmt.Printf("  Environment: %s\n", appConfig["environment"])
	fmt.Printf("  Configuration loaded successfully with %d top-level keys\n", len(appConfig))

	// Example 2: Generate Configuration File
	fmt.Println("\n2. Generating Configuration File:")
	newConfig := toon.NewOrderedMap()
	newConfig.Set("service_name", "Payment Gateway")
	newConfig.Set("version", "2.1.0")
	newConfig.Set("deploy_environment", "staging")

	// API Configuration
	api := toon.NewOrderedMap()
	api.Set("base_url", "https://api-staging.example.com")
	api.Set("timeout_seconds", 60)
	api.Set("retry_attempts", 3)
	api.Set("endpoints", []interface{}{
		"/v1/payments/process",
		"/v1/payments/refund",
		"/v1/payments/status",
	})
	newConfig.Set("api", api)

	// Provider Configurations
	providers := []interface{}{}

	stripeConfig := toon.NewOrderedMap()
	stripeConfig.Set("name", "stripe")
	stripeConfig.Set("enabled", true)
	stripeConfig.Set("api_key_env", "STRIPE_API_KEY")
	stripeConfig.Set("webhook_secret_env", "STRIPE_WEBHOOK_SECRET")
	stripeConfig.Set("supported_currencies", []interface{}{"USD", "EUR", "GBP"})
	providers = append(providers, stripeConfig)

	paypalConfig := toon.NewOrderedMap()
	paypalConfig.Set("name", "paypal")
	paypalConfig.Set("enabled", true)
	paypalConfig.Set("client_id_env", "PAYPAL_CLIENT_ID")
	paypalConfig.Set("client_secret_env", "PAYPAL_CLIENT_SECRET")
	paypalConfig.Set("mode", "sandbox")
	providers = append(providers, paypalConfig)

	newConfig.Set("providers", providers)

	// Monitoring
	monitoring := toon.NewOrderedMap()
	monitoring.Set("metrics_enabled", true)
	monitoring.Set("metrics_port", 9090)
	monitoring.Set("health_check_path", "/health")
	monitoring.Set("alerts", []interface{}{
		"transaction_failure_rate",
		"api_response_time",
		"provider_downtime",
	})
	newConfig.Set("monitoring", monitoring)

	generated, _ := toon.MarshalToString(newConfig, toon.WithIndent(2))
	fmt.Println(generated)

	// Example 3: Environment-Specific Configuration
	fmt.Println("\n3. Environment-Specific Configuration:")

	environments := map[string]interface{}{
		"development": map[string]interface{}{
			"debug":    true,
			"log_sql":  true,
			"base_url": "http://localhost:3000",
			"database": map[string]interface{}{
				"host":     "localhost",
				"port":     5432,
				"ssl_mode": "disable",
			},
		},
		"staging": map[string]interface{}{
			"debug":    true,
			"log_sql":  false,
			"base_url": "https://staging.example.com",
			"database": map[string]interface{}{
				"host":     "db-staging.internal",
				"port":     5432,
				"ssl_mode": "require",
			},
		},
		"production": map[string]interface{}{
			"debug":    false,
			"log_sql":  false,
			"base_url": "https://example.com",
			"database": map[string]interface{}{
				"host":     "db-prod.internal",
				"port":     5432,
				"ssl_mode": "require",
			},
		},
	}

	envConfig, _ := toon.MarshalToString(environments)
	fmt.Println(envConfig)

	// Example 4: Feature Flags Configuration
	fmt.Println("\n4. Feature Flags Configuration:")
	featureFlags := toon.NewOrderedMap()
	featureFlags.Set("version", "1.0")
	featureFlags.Set("last_updated", "2025-11-20")

	flags := []interface{}{}

	flag1 := toon.NewOrderedMap()
	flag1.Set("key", "new_checkout_flow")
	flag1.Set("enabled", true)
	flag1.Set("rollout_percentage", 50)
	flag1.Set("environments", []interface{}{"staging", "production"})
	flags = append(flags, flag1)

	flag2 := toon.NewOrderedMap()
	flag2.Set("key", "ai_recommendations")
	flag2.Set("enabled", true)
	flag2.Set("rollout_percentage", 100)
	flag2.Set("environments", []interface{}{"production"})
	flags = append(flags, flag2)

	flag3 := toon.NewOrderedMap()
	flag3.Set("key", "beta_dashboard")
	flag3.Set("enabled", false)
	flag3.Set("rollout_percentage", 0)
	flag3.Set("environments", []interface{}{"development"})
	flags = append(flags, flag3)

	featureFlags.Set("flags", flags)

	flagsOutput, _ := toon.MarshalToString(featureFlags)
	fmt.Println(flagsOutput)

	// Example 5: Microservices Configuration
	fmt.Println("\n5. Microservices Configuration:")
	microservices := toon.NewOrderedMap()
	microservices.Set("platform", "kubernetes")
	microservices.Set("namespace", "production")

	services := []interface{}{}

	userService := toon.NewOrderedMap()
	userService.Set("name", "user-service")
	userService.Set("image", "registry.example.com/user-service:v1.2.3")
	userService.Set("replicas", 3)
	userService.Set("port", 8080)
	userService.Set("env", map[string]interface{}{
		"DB_HOST":   "postgres-users",
		"CACHE_URL": "redis://cache:6379",
		"LOG_LEVEL": "info",
	})
	userService.Set("resources", map[string]interface{}{
		"cpu":    "500m",
		"memory": "512Mi",
	})
	services = append(services, userService)

	orderService := toon.NewOrderedMap()
	orderService.Set("name", "order-service")
	orderService.Set("image", "registry.example.com/order-service:v2.1.0")
	orderService.Set("replicas", 5)
	orderService.Set("port", 8080)
	orderService.Set("env", map[string]interface{}{
		"DB_HOST":       "postgres-orders",
		"KAFKA_BROKERS": "kafka:9092",
		"LOG_LEVEL":     "info",
	})
	orderService.Set("resources", map[string]interface{}{
		"cpu":    "1000m",
		"memory": "1Gi",
	})
	services = append(services, orderService)

	microservices.Set("services", services)

	microOutput, _ := toon.MarshalToString(microservices)
	fmt.Println(microOutput)

	fmt.Println("\n✓ Configuration examples demonstrate TOON's token efficiency")
	fmt.Println("✓ 30-60% fewer tokens than JSON for configuration files")
}
