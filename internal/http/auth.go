package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/rahulkj/cp-discovery/internal/model"
)

// AuthConfig represents authentication configuration for HTTP clients
type AuthConfig interface {
	ApplyAuth(req *http.Request)
}

// AuthParams holds all authentication parameters for a component
type AuthParams struct {
	BasicAuthUsername string
	BasicAuthPassword string
	BearerToken       string
	APIKey            string
	APIKeyHeader      string
	LDAPEnabled       bool
	LDAPServer        string
	LDAPUsername      string
	LDAPPassword      string
	LDAPBaseDN        string
	OAuthEnabled      bool
	OAuthClientID     string
	OAuthClientSecret string
	OAuthTokenURL     string
	OAuthScopes       string
}

// OAuth token cache to avoid repeated token requests
var (
	oauthTokenCache   = make(map[string]*cachedToken)
	oauthTokenCacheMu sync.RWMutex
)

type cachedToken struct {
	accessToken string
	expiresAt   time.Time
}

// OAuthTokenResponse represents the OAuth token response
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// ApplySchemaRegistryAuth applies authentication to Schema Registry HTTP requests
func ApplySchemaRegistryAuth(req *http.Request, config model.SchemaRegistryConfig) {
	params := AuthParams{
		BasicAuthUsername: config.BasicAuthUsername,
		BasicAuthPassword: config.BasicAuthPassword,
		BearerToken:       config.BearerToken,
		APIKey:            config.APIKey,
		APIKeyHeader:      config.APIKeyHeader,
		LDAPEnabled:       config.LDAPEnabled,
		LDAPServer:        config.LDAPServer,
		LDAPUsername:      config.LDAPUsername,
		LDAPPassword:      config.LDAPPassword,
		LDAPBaseDN:        config.LDAPBaseDN,
		OAuthEnabled:      config.OAuthEnabled,
		OAuthClientID:     config.OAuthClientID,
		OAuthClientSecret: config.OAuthClientSecret,
		OAuthTokenURL:     config.OAuthTokenURL,
		OAuthScopes:       config.OAuthScopes,
	}
	applyAuth(req, params)
}

// ApplyKafkaConnectAuth applies authentication to Kafka Connect HTTP requests
func ApplyKafkaConnectAuth(req *http.Request, config model.KafkaConnectConfig) {
	params := AuthParams{
		BasicAuthUsername: config.BasicAuthUsername,
		BasicAuthPassword: config.BasicAuthPassword,
		BearerToken:       config.BearerToken,
		APIKey:            config.APIKey,
		APIKeyHeader:      config.APIKeyHeader,
		LDAPEnabled:       config.LDAPEnabled,
		LDAPServer:        config.LDAPServer,
		LDAPUsername:      config.LDAPUsername,
		LDAPPassword:      config.LDAPPassword,
		LDAPBaseDN:        config.LDAPBaseDN,
		OAuthEnabled:      config.OAuthEnabled,
		OAuthClientID:     config.OAuthClientID,
		OAuthClientSecret: config.OAuthClientSecret,
		OAuthTokenURL:     config.OAuthTokenURL,
		OAuthScopes:       config.OAuthScopes,
	}
	applyAuth(req, params)
}

// ApplyKsqlDBAuth applies authentication to ksqlDB HTTP requests
func ApplyKsqlDBAuth(req *http.Request, config model.KsqlDBConfig) {
	params := AuthParams{
		BasicAuthUsername: config.BasicAuthUsername,
		BasicAuthPassword: config.BasicAuthPassword,
		BearerToken:       config.BearerToken,
		APIKey:            config.APIKey,
		APIKeyHeader:      config.APIKeyHeader,
		LDAPEnabled:       config.LDAPEnabled,
		LDAPServer:        config.LDAPServer,
		LDAPUsername:      config.LDAPUsername,
		LDAPPassword:      config.LDAPPassword,
		LDAPBaseDN:        config.LDAPBaseDN,
		OAuthEnabled:      config.OAuthEnabled,
		OAuthClientID:     config.OAuthClientID,
		OAuthClientSecret: config.OAuthClientSecret,
		OAuthTokenURL:     config.OAuthTokenURL,
		OAuthScopes:       config.OAuthScopes,
	}
	applyAuth(req, params)
}

// ApplyRestProxyAuth applies authentication to REST Proxy HTTP requests
func ApplyRestProxyAuth(req *http.Request, config model.RestProxyConfig) {
	params := AuthParams{
		BasicAuthUsername: config.BasicAuthUsername,
		BasicAuthPassword: config.BasicAuthPassword,
		BearerToken:       config.BearerToken,
		APIKey:            config.APIKey,
		APIKeyHeader:      config.APIKeyHeader,
		LDAPEnabled:       config.LDAPEnabled,
		LDAPServer:        config.LDAPServer,
		LDAPUsername:      config.LDAPUsername,
		LDAPPassword:      config.LDAPPassword,
		LDAPBaseDN:        config.LDAPBaseDN,
		OAuthEnabled:      config.OAuthEnabled,
		OAuthClientID:     config.OAuthClientID,
		OAuthClientSecret: config.OAuthClientSecret,
		OAuthTokenURL:     config.OAuthTokenURL,
		OAuthScopes:       config.OAuthScopes,
	}
	applyAuth(req, params)
}

// ApplyControlCenterAuth applies authentication to Control Center HTTP requests
func ApplyControlCenterAuth(req *http.Request, config model.ControlCenterConfig) {
	params := AuthParams{
		BasicAuthUsername: config.BasicAuthUsername,
		BasicAuthPassword: config.BasicAuthPassword,
		BearerToken:       config.BearerToken,
		APIKey:            config.APIKey,
		APIKeyHeader:      config.APIKeyHeader,
		LDAPEnabled:       config.LDAPEnabled,
		LDAPServer:        config.LDAPServer,
		LDAPUsername:      config.LDAPUsername,
		LDAPPassword:      config.LDAPPassword,
		LDAPBaseDN:        config.LDAPBaseDN,
		OAuthEnabled:      config.OAuthEnabled,
		OAuthClientID:     config.OAuthClientID,
		OAuthClientSecret: config.OAuthClientSecret,
		OAuthTokenURL:     config.OAuthTokenURL,
		OAuthScopes:       config.OAuthScopes,
	}
	applyAuth(req, params)
}

// ApplyPrometheusAuth applies authentication to Prometheus HTTP requests
func ApplyPrometheusAuth(req *http.Request, config model.PrometheusConfig) {
	params := AuthParams{
		BasicAuthUsername: config.BasicAuthUsername,
		BasicAuthPassword: config.BasicAuthPassword,
		BearerToken:       config.BearerToken,
		APIKey:            config.APIKey,
		APIKeyHeader:      config.APIKeyHeader,
		LDAPEnabled:       config.LDAPEnabled,
		LDAPServer:        config.LDAPServer,
		LDAPUsername:      config.LDAPUsername,
		LDAPPassword:      config.LDAPPassword,
		LDAPBaseDN:        config.LDAPBaseDN,
		OAuthEnabled:      config.OAuthEnabled,
		OAuthClientID:     config.OAuthClientID,
		OAuthClientSecret: config.OAuthClientSecret,
		OAuthTokenURL:     config.OAuthTokenURL,
		OAuthScopes:       config.OAuthScopes,
	}
	applyAuth(req, params)
}

// ApplyAlertmanagerAuth applies authentication to Alertmanager HTTP requests
func ApplyAlertmanagerAuth(req *http.Request, config model.AlertmanagerConfig) {
	params := AuthParams{
		BasicAuthUsername: config.BasicAuthUsername,
		BasicAuthPassword: config.BasicAuthPassword,
		BearerToken:       config.BearerToken,
		APIKey:            config.APIKey,
		APIKeyHeader:      config.APIKeyHeader,
		LDAPEnabled:       config.LDAPEnabled,
		LDAPServer:        config.LDAPServer,
		LDAPUsername:      config.LDAPUsername,
		LDAPPassword:      config.LDAPPassword,
		LDAPBaseDN:        config.LDAPBaseDN,
		OAuthEnabled:      config.OAuthEnabled,
		OAuthClientID:     config.OAuthClientID,
		OAuthClientSecret: config.OAuthClientSecret,
		OAuthTokenURL:     config.OAuthTokenURL,
		OAuthScopes:       config.OAuthScopes,
	}
	applyAuth(req, params)
}

// applyAuth applies authentication headers to an HTTP request based on configured auth type
// Priority: OAuth > LDAP > Bearer Token > API Key > Basic Auth
func applyAuth(req *http.Request, params AuthParams) {
	// OAuth/SSO Authentication (highest priority)
	if params.OAuthEnabled && params.OAuthClientID != "" && params.OAuthClientSecret != "" && params.OAuthTokenURL != "" {
		token, err := getOAuthToken(params.OAuthClientID, params.OAuthClientSecret, params.OAuthTokenURL, params.OAuthScopes)
		if err == nil && token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
			return
		}
		// If OAuth fails, fall through to other auth methods
	}

	// LDAP Authentication
	if params.LDAPEnabled && params.LDAPServer != "" && params.LDAPUsername != "" {
		// For LDAP, we authenticate using the LDAP credentials as Basic Auth
		// In a real-world scenario, you might need to perform LDAP bind and obtain a session token
		// For HTTP services, LDAP credentials are typically validated on the server side
		// and we pass them as Basic Auth or obtain a token first
		token, err := authenticateLDAP(params.LDAPServer, params.LDAPUsername, params.LDAPPassword, params.LDAPBaseDN)
		if err == nil && token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
			return
		}
		// If LDAP token retrieval fails, fall back to Basic Auth with LDAP credentials
		if params.LDAPUsername != "" {
			req.SetBasicAuth(params.LDAPUsername, params.LDAPPassword)
			return
		}
	}

	// Bearer Token Authentication (pre-configured token, JWT, etc.)
	if params.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+params.BearerToken)
		return
	}

	// API Key Authentication
	if params.APIKey != "" {
		headerName := "X-API-Key" // Default header
		if params.APIKeyHeader != "" {
			headerName = params.APIKeyHeader
		}
		req.Header.Set(headerName, params.APIKey)
		return
	}

	// Basic Authentication (lowest priority)
	if params.BasicAuthUsername != "" {
		req.SetBasicAuth(params.BasicAuthUsername, params.BasicAuthPassword)
		return
	}

	// No authentication configured
}

// getOAuthToken retrieves an OAuth access token using client credentials flow
func getOAuthToken(clientID, clientSecret, tokenURL, scopes string) (string, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("%s:%s:%s", clientID, tokenURL, scopes)

	// Check cache first
	oauthTokenCacheMu.RLock()
	cached, exists := oauthTokenCache[cacheKey]
	oauthTokenCacheMu.RUnlock()

	if exists && cached.expiresAt.After(time.Now().Add(5*time.Minute)) {
		return cached.accessToken, nil
	}

	// Prepare token request
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	if scopes != "" {
		data.Set("scope", scopes)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("requesting OAuth token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OAuth token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("decoding OAuth token response: %w", err)
	}

	// Cache the token
	oauthTokenCacheMu.Lock()
	oauthTokenCache[cacheKey] = &cachedToken{
		accessToken: tokenResp.AccessToken,
		expiresAt:   time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}
	oauthTokenCacheMu.Unlock()

	return tokenResp.AccessToken, nil
}

// authenticateLDAP performs LDAP authentication and returns a token
// Note: This is a simplified implementation. In production, you may need:
// 1. To use a proper LDAP library (e.g., go-ldap/ldap) for LDAP bind operations
// 2. To implement token exchange with your authentication service
// 3. To handle different LDAP configurations (SSL/TLS, search filters, etc.)
func authenticateLDAP(ldapServer, username, password, baseDN string) (string, error) {
	// For HTTP-based services that support LDAP, they typically:
	// 1. Accept Basic Auth with LDAP credentials, OR
	// 2. Provide a separate authentication endpoint that validates LDAP credentials and returns a token

	// This is a placeholder implementation that returns empty string
	// signaling the caller to fall back to Basic Auth with LDAP credentials
	//
	// To implement full LDAP support, you would:
	// - Use go-ldap/ldap library to perform LDAP bind
	// - Connect to the LDAP server
	// - Authenticate the user
	// - Optionally exchange credentials for a service token

	// Example implementation would look like:
	// import "github.com/go-ldap/ldap/v3"
	//
	// l, err := ldap.DialURL(ldapServer)
	// if err != nil {
	//     return "", err
	// }
	// defer l.Close()
	//
	// userDN := fmt.Sprintf("uid=%s,%s", username, baseDN)
	// err = l.Bind(userDN, password)
	// if err != nil {
	//     return "", fmt.Errorf("LDAP bind failed: %w", err)
	// }
	//
	// Exchange credentials for a service token (service-specific)
	// token, err := exchangeForToken(username, password)
	// return token, err

	return "", fmt.Errorf("LDAP authentication not fully implemented, falling back to Basic Auth")
}
