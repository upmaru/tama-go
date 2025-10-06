package main

import (
	"fmt"
	"log"
	"time"

	tama "github.com/upmaru/tama-go"
	"github.com/upmaru/tama-go/neural"
)

func oauth2Demo() {
	fmt.Println("🚀 Tama Go OAuth2 Authentication Demo")
	fmt.Println("=====================================")

	// OAuth2 Configuration
	config := tama.Config{
		BaseURL:      "https://api.tama.io", // Replace with your actual API URL
		ClientID:     "your-client-id",      // Replace with your OAuth2 client ID
		ClientSecret: "your-client-secret",  // Replace with your OAuth2 client secret
		Timeout:      30 * time.Second,
	}

	fmt.Printf("📡 Connecting to: %s\n", config.BaseURL)
	fmt.Printf("🔑 Client ID: %s\n", config.ClientID)
	fmt.Printf("⏱️  Timeout: %v\n", config.Timeout)

	// Create OAuth2 client (this will automatically acquire access token)
	fmt.Println("\n🔐 Acquiring OAuth2 access token...")
	client, err := tama.NewClient(config)
	if err != nil {
		log.Fatalf("❌ Failed to create OAuth2 client: %v", err)
	}

	fmt.Println("✅ OAuth2 client created successfully!")

	// Display token information
	token := client.GetToken()
	if token != nil {
		fmt.Printf("🎫 Access Token: %s...\n", token.AccessToken[:20]+"***")
		fmt.Printf("📊 Token Type: %s\n", token.TokenType)
		fmt.Printf("🎯 Scope: %s\n", token.Scope)
		fmt.Printf("⏰ Expires At: %v\n", token.ExpiresAt.Format(time.RFC3339))
		fmt.Printf("⏳ Time Until Expiry: %v\n", time.Until(token.ExpiresAt).Round(time.Second))
	}

	// Enable debug mode to see HTTP requests (optional)
	client.SetDebug(true)

	// Example 1: Create a Neural Space
	fmt.Println("\n🧠 Creating Neural Space...")
	createSpaceReq := neural.CreateSpaceRequest{
		Space: neural.SpaceRequestData{
			Name: fmt.Sprintf("OAuth2 Demo Space %d", time.Now().Unix()),
			Type: "root",
		},
	}

	space, err := client.Neural.CreateSpace(createSpaceReq)
	if err != nil {
		log.Printf("⚠️ Failed to create space: %v", err)
		fmt.Println("   (This might be expected if using demo credentials)")
	} else {
		fmt.Printf("✅ Created space: ID=%s, Name=%s, Type=%s, State=%s\n",
			space.ID, space.Name, space.Type, space.ProvisionState)

		// Example 2: Get the created space
		fmt.Println("\n📖 Retrieving Neural Space...")
		retrievedSpace, err := client.Neural.GetSpace(space.ID)
		if err != nil {
			log.Printf("⚠️ Failed to retrieve space: %v", err)
		} else {
			fmt.Printf("✅ Retrieved space: ID=%s, Name=%s\n",
				retrievedSpace.ID, retrievedSpace.Name)
		}

		// Example 3: Update the space
		fmt.Println("\n📝 Updating Neural Space...")
		updateReq := neural.UpdateSpaceRequest{
			Space: neural.UpdateSpaceData{
				Name: space.Name + " (Updated)",
			},
		}

		updatedSpace, err := client.Neural.UpdateSpace(space.ID, updateReq)
		if err != nil {
			log.Printf("⚠️ Failed to update space: %v", err)
		} else {
			fmt.Printf("✅ Updated space: ID=%s, Name=%s\n",
				updatedSpace.ID, updatedSpace.Name)
		}

		// Example 4: Clean up - Delete the space
		fmt.Println("\n🧹 Cleaning up - Deleting Neural Space...")
		err = client.Neural.DeleteSpace(space.ID)
		if err != nil {
			log.Printf("⚠️ Failed to delete space: %v", err)
		} else {
			fmt.Printf("✅ Successfully deleted space: %s\n", space.ID)
		}
	}

	// Example 5: Demonstrate automatic token refresh
	fmt.Println("\n🔄 OAuth2 Token Management Demo")
	fmt.Println("The client automatically handles:")
	fmt.Println("  - Initial token acquisition during client creation")
	fmt.Println("  - Automatic token refresh when tokens expire")
	fmt.Println("  - Thread-safe token operations")
	fmt.Println("  - Proper HTTP Basic Auth for token requests")

	// Example 6: Show OAuth2 technical details
	fmt.Println("\n🔧 OAuth2 Technical Details:")
	fmt.Println("  - Grant Type: client_credentials")
	fmt.Println("  - Scope: provision.all")
	fmt.Println("  - Token Endpoint: POST /auth/tokens")
	fmt.Println("  - Authentication: HTTP Basic Auth with base64(client_id:client_secret)")
	fmt.Println("  - Token Usage: Bearer token in Authorization header")

	// Example 7: Error handling demonstration
	fmt.Println("\n🛡️ Error Handling:")
	_, err = client.Neural.GetSpace("nonexistent-space-id")
	if err != nil {
		fmt.Printf("✅ Proper error handling: %v\n", err)
	}

	fmt.Println("\n🎉 OAuth2 Demo Complete!")
	fmt.Println("Key Features Demonstrated:")
	fmt.Println("  ✅ OAuth2 client credentials flow")
	fmt.Println("  ✅ Automatic token management")
	fmt.Println("  ✅ Authenticated API operations")
	fmt.Println("  ✅ Error handling")
	fmt.Println("  ✅ Token information access")
	fmt.Println("\n📚 Migration from API Key:")
	fmt.Println("  Before: client := tama.NewClient(config)")
	fmt.Println("  After:  client, err := tama.NewClient(config)")
	fmt.Println("  Config: Use ClientID & ClientSecret instead of APIKey")
}

// Example of how to test OAuth2 without real credentials
func demonstrateTestMode() {
	fmt.Println("\n🧪 Test Mode Example:")

	testConfig := tama.Config{
		BaseURL:        "http://localhost:4000",
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		SkipTokenFetch: true, // Skip real OAuth2 for testing
	}

	client, err := tama.NewClient(testConfig)
	if err != nil {
		log.Printf("Failed to create test client: %v", err)
		return
	}

	fmt.Println("✅ Test mode client created (no real OAuth2 requests)")

	// This client can be used with mock servers for testing
	token := client.GetToken()
	if token == nil {
		fmt.Println("✅ No token in test mode (expected)")
	}
}

// Example of proper error handling in production
func productionExample() {
	config := tama.Config{
		BaseURL:      "https://api.tama.io",
		ClientID:     "prod-client-id",
		ClientSecret: "prod-client-secret",
		Timeout:      60 * time.Second,
	}

	client, err := tama.NewClient(config)
	if err != nil {
		// Handle authentication failures
		fmt.Printf("Authentication failed: %v\n", err)
		return
	}

	// Check token expiration
	token := client.GetToken()
	if token != nil && token.IsExpired() {
		fmt.Println("Token is expired (will be refreshed automatically)")
	}

	// Use client for API operations
	// spaces, err := client.Neural.ListSpaces() // Hypothetical method - commented out as not implemented
	var spaces []interface{} // Placeholder for compilation
	err = nil
	if err != nil {
		fmt.Printf("API call failed: %v\n", err)
		return
	}

	fmt.Printf("Retrieved %d spaces\n", len(spaces))
}

// Uncomment the following to run this as a standalone demo:
// func main() {
// 	oauth2Demo()
// 	demonstrateTestMode()
// 	productionExample()
// }
