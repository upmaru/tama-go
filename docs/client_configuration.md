## Client Configuration

### NewClient(config Config) *Client

Creates a new Tama API client with the provided configuration.

**Parameters:**
- `config` (Config): Configuration object containing:
  - `BaseURL` (string): The base URL of the Tama API (required)
  - `APIKey` (string): Your API authentication key (required)
  - `Timeout` (time.Duration): Request timeout (optional, default: 30s)

**Returns:**
- `*Client`: Configured client instance

**Example:**
```go
config := tama.Config{
    BaseURL: "https://api.tama.io",
    APIKey:  "your-api-key",
    Timeout: 30 * time.Second,
}
client := tama.NewClient(config)
```

### Client Methods

#### SetAPIKey(apiKey string)

Updates the API key for authentication.

**Parameters:**
- `apiKey` (string): New API key

#### SetDebug(debug bool)

Enables or disables debug mode for HTTP requests.

**Parameters:**
- `debug` (bool): Enable/disable debug mode
