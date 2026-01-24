# o2go

**A simple and flexible OAuth2 client for Go.**  
Handles **authorization code** and **refresh token flows**, with **context support**, **custom HTTP clients**, and **rich error handling**.  

---

## Why o2go?

Many Go OAuth2 libraries:

- Force opinions on **state**, **scopes**, or **token storage**.  
- Make it hard to inject custom HTTP clients.  
- Limit error inspection or context cancellation.  

**o2go solves these gaps:**

- Minimal, focused on token exchange flows only.  
- Full control over **state**, **scopes**, and **token persistence**.  
- Supports **context-aware** requests and **custom HTTP clients**.  
- Rich `HTTPError` including **status codes** and **response body** for debugging.  

You can also use **custom providers**:

```go
type MyProvider struct{}

func (p *MyProvider) AuthURL() string  { return "https://example.com/oauth/authorize" }
func (p *MyProvider) TokenURL() string { return "https://example.com/oauth/token" }
```

Or use ready-made providers from `o2go-providers` (like Google, GitHub, etc.) when available.

---

## Installation

```bash
go get github.com/gokhanaltun/o2go
```

---

## Quick Start

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "net/http"

    "github.com/gokhanaltun/o2go"
)

func main() {
    cfg := &o2go.Config{
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        RedirectURL:  "https://yourapp.com/callback",
    }

    client, err := o2go.New(&MyProvider{}, cfg)
    if err != nil {
        panic(err)
    }

    authURL, err := client.AuthCodeURL("random-state")
    if err != nil {
        panic(err)
    }
    fmt.Println("Open this URL:", authURL)

    code := "code-from-callback"
    resp, err := client.ExchangeAuthCode(context.Background(), code)
    if err != nil {
        var he *o2go.HTTPError
        if errors.As(err, &he) {
            fmt.Println("Status:", he.Status)
            fmt.Println("Body:", string(he.Body))
        } else {
            fmt.Println("Other error:", err)
        }
    }

    fmt.Println("Access Token:", resp.Data["access_token"])
}
```

---

## Advanced Usage

### Custom HTTP client

```go
httpClient := &http.Client{Timeout: 10 * time.Second}

client, err := o2go.New(&MyProvider{}, cfg, o2go.WithHttpClient(httpClient))
```

### Refresh token flow

```go
refreshToken := resp.Data["refresh_token"].(string)
newResp, err := client.ExchangeRefreshToken(context.Background(), refreshToken)
if err != nil {
    var he *o2go.HTTPError
    if errors.As(err, &he) {
        fmt.Println("Status:", he.Status)
        fmt.Println("Body:", string(he.Body))
    } else {
        fmt.Println("Other error:", err)
    }
}

fmt.Println("New Access Token:", newResp.Data["access_token"])
```

### Adding extra token parameters

```go
cfg.TokenParams = map[string]string{
    "audience": "my-api",
}
```

These parameters are automatically merged into the token request while respecting reserved keys like `grant_type` or `client_id`.

---

## Error Handling

`ExchangeAuthCode` and `ExchangeRefreshToken` return `HTTPError` for non-2xx responses:

```go
var he *o2go.HTTPError
if errors.As(err, &he) {
    fmt.Println("Status:", he.Status)
    fmt.Println("Body:", string(he.Body))
}
```

This allows full inspection of HTTP status codes and response bodies.

---

## API Reference

### `func New(provider Provider, cfg *Config, opts ...Option) (*OAuth2, error)`

* Creates a new OAuth2 instance.
* Optional `WithHttpClient(client *http.Client)`.

### `func (o *OAuth2) AuthCodeURL(state string) (string, error)`

Builds the authorization URL. App manages `state`.

### `func (o *OAuth2) ExchangeAuthCode(ctx context.Context, code string) (*ExchangeResponse, error)`

Exchanges **authorization code** for tokens.

### `func (o *OAuth2) ExchangeRefreshToken(ctx context.Context, refreshToken string) (*ExchangeResponse, error)`

Exchanges **refresh token** for new tokens.

### `type ExchangeResponse`

```go
type ExchangeResponse struct {
    Raw  []byte
    Data map[string]any
}
```

* `Decode(v any) error` — decode raw JSON response.

### `type HTTPError`

* `Status int` — HTTP status code
* `Body []byte` — raw response body
* `Err error` — wrapped underlying error

Supports `errors.Unwrap()`.

---

## License

MIT License — free to use, modify, and distribute.
