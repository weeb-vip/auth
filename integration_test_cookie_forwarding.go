package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

// TestCookieForwardingIntegration tests the complete cookie forwarding flow
// from auth service through Apollo Router to verify Set-Cookie headers are properly forwarded
func TestCookieForwardingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test cases for different endpoints
	testCases := []struct {
		name        string
		endpoint    string
		description string
	}{
		{
			name:        "deployed_auth",
			endpoint:    "http://localhost:3003/graphql",
			description: "Direct auth service (deployed version)",
		},
		{
			name:        "apollo_router",
			endpoint:    "http://localhost:4000/graphql",
			description: "Apollo Router forwarding to auth service",
		},
	}

	// GraphQL mutation for creating a session
	mutation := `{
		"query": "mutation { CreateSession(input: {username: \"test\", password: \"test\"}) { id } }"
	}`

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create HTTP request
			req, err := http.NewRequest("POST", tc.endpoint, strings.NewReader(mutation))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			// Send request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Read response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Verify response status
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body))
			}

			// Check for Set-Cookie headers
			setCookieHeaders := resp.Header.Values("Set-Cookie")
			if len(setCookieHeaders) == 0 {
				// Check lowercase header (Apollo Router might normalize)
				setCookieHeaders = resp.Header.Values("set-cookie")
			}

			if len(setCookieHeaders) == 0 {
				t.Fatalf("No Set-Cookie headers found in response from %s", tc.description)
			}

			// Verify we have both access_token and refresh_token cookies
			var hasAccessToken, hasRefreshToken bool
			var accessTokenDomain, refreshTokenDomain string

			for _, cookie := range setCookieHeaders {
				if strings.Contains(cookie, "access_token=") {
					hasAccessToken = true
					// Extract domain from cookie string
					if domainIndex := strings.Index(cookie, "Domain="); domainIndex != -1 {
						domainPart := cookie[domainIndex+7:] // Skip "Domain="
						if semicolonIndex := strings.Index(domainPart, ";"); semicolonIndex != -1 {
							accessTokenDomain = domainPart[:semicolonIndex]
						} else {
							accessTokenDomain = domainPart
						}
					}
				}
				if strings.Contains(cookie, "refresh_token=") {
					hasRefreshToken = true
					// Extract domain from cookie string
					if domainIndex := strings.Index(cookie, "Domain="); domainIndex != -1 {
						domainPart := cookie[domainIndex+7:] // Skip "Domain="
						if semicolonIndex := strings.Index(domainPart, ";"); semicolonIndex != -1 {
							refreshTokenDomain = domainPart[:semicolonIndex]
						} else {
							refreshTokenDomain = domainPart
						}
					}
				}
			}

			if !hasAccessToken {
				t.Errorf("Missing access_token cookie in response from %s", tc.description)
			}

			if !hasRefreshToken {
				t.Errorf("Missing refresh_token cookie in response from %s", tc.description)
			}

			// Verify cookie domains are normalized (without leading dot)
			expectedDomain := "weeb.vip" // Go normalizes ".weeb.vip" to "weeb.vip"

			if accessTokenDomain != expectedDomain {
				t.Errorf("Expected access_token domain '%s', got '%s' from %s",
					expectedDomain, accessTokenDomain, tc.description)
			}

			if refreshTokenDomain != expectedDomain {
				t.Errorf("Expected refresh_token domain '%s', got '%s' from %s",
					expectedDomain, refreshTokenDomain, tc.description)
			}

			// Log success for debugging
			t.Logf("✓ %s: Successfully received %d Set-Cookie headers with domain=%s",
				tc.description, len(setCookieHeaders), expectedDomain)
		})
	}
}

// TestCookieDomainNormalization tests that Go's net/http normalizes cookie domains
func TestCookieDomainNormalization(t *testing.T) {
	// This test documents the expected behavior that Go's net/http package
	// normalizes cookie domains by removing leading dots per RFC 6265

	testCases := []struct {
		configDomain   string
		expectedDomain string
		description    string
	}{
		{
			configDomain:   ".weeb.vip",
			expectedDomain: "weeb.vip",
			description:    "Leading dot should be removed by Go's net/http",
		},
		{
			configDomain:   "localhost",
			expectedDomain: "localhost",
			description:    "Domain without dot should remain unchanged",
		},
		{
			configDomain:   ".example.com",
			expectedDomain: "example.com",
			description:    "Generic domain with leading dot should be normalized",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.configDomain, func(t *testing.T) {
			// Create a cookie with the config domain
			cookie := &http.Cookie{
				Name:   "test_cookie",
				Value:  "test_value",
				Domain: tc.configDomain,
			}

			// Simulate what http.SetCookie does by creating a response
			// and checking what domain actually gets set
			resp := &http.Response{
				Header: make(http.Header),
			}

			// This is what our auth service does internally
			http.SetCookie(resp, cookie)

			// Check the actual Set-Cookie header that was generated
			setCookieHeader := resp.Header.Get("Set-Cookie")

			if !strings.Contains(setCookieHeader, fmt.Sprintf("Domain=%s", tc.expectedDomain)) {
				t.Errorf("Expected domain '%s' in Set-Cookie header, got: %s (%s)",
					tc.expectedDomain, setCookieHeader, tc.description)
			}

			t.Logf("✓ %s: Cookie domain '%s' normalized to '%s'",
				tc.description, tc.configDomain, tc.expectedDomain)
		})
	}
}

func TestMain(m *testing.M) {
	// This integration test requires port forwards to be running:
	// kubectl port-forward pod/auth-xxx -n weeb 3003:3000
	// kubectl port-forward -n weeb deployment/apollo-router 4000:5000

	fmt.Println("Running cookie forwarding integration tests...")
	fmt.Println("Note: This test requires active port forwards to auth service (3003) and Apollo Router (4000)")

	code := m.Run()
	os.Exit(code)
}