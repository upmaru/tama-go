# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **BREAKING CHANGE**: Migrated from API Key authentication to OAuth2 with client credentials flow
- `NewClient()` now returns `(*Client, error)` instead of `*Client` to handle authentication failures
- `Config` struct now requires `ClientID` and `ClientSecret` instead of `APIKey`
- Client automatically handles token acquisition and refresh using `/auth/tokens` endpoint

### Added
- OAuth2 client credentials flow with automatic token refresh
- `SkipTokenFetch` config option for testing to avoid real HTTP requests
- Request interceptor that ensures valid tokens before API calls
- `GetToken()` method to retrieve current token information
- Token expiration tracking with 30-second buffer for refresh
- HTTP Basic Authentication with base64 encoded client credentials for token requests

### Removed
- **BREAKING CHANGE**: Removed `SetAPIKey()` method
- **BREAKING CHANGE**: Removed `APIKey` field from `Config` struct

### Security
- Tokens are automatically refreshed when expired
- Client credentials are securely encoded using HTTP Basic Auth
- Token information is thread-safe with mutex protection

### Migration Guide

#### Before (API Key)
```go
config := tama.Config{
    BaseURL: "https://api.tama.io",
    APIKey:  "your-api-key",
    Timeout: 30 * time.Second,
}

client := tama.NewClient(config)
```

#### After (OAuth2)
```go
config := tama.Config{
    BaseURL:      "https://api.tama.io", 
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    Timeout:      30 * time.Second,
}

client, err := tama.NewClient(config)
if err != nil {
    // Handle authentication error
    panic(err)
}
```

#### For Testing
```go
config := tama.Config{
    BaseURL:        "https://api.tama.io",
    ClientID:       "test-client-id", 
    ClientSecret:   "test-client-secret",
    SkipTokenFetch: true, // Skip real authentication for tests
}

client, err := tama.NewClient(config)
```

### OAuth2 Flow Details
- **Token Endpoint**: `POST /auth/tokens`
- **Grant Type**: `client_credentials`
- **Scope**: `provision.all`
- **Authentication**: HTTP Basic Auth with `Authorization: Bearer base64(client_id:client_secret)`
- **Token Response**: JSON with `access_token`, `token_type`, `scope`, and `expires_in`
- **Token Refresh**: Automatic refresh when token expires or within 30 seconds of expiration