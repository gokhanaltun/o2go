# o2go

[![Go Reference](https://pkg.go.dev/badge/github.com/gokhanaltun/o2go.svg)](https://pkg.go.dev/github.com/gokhanaltun/o2go)
[![Go Report Card](https://goreportcard.com/badge/github.com/gokhanaltun/o2go)](https://goreportcard.com/report/github.com/gokhanaltun/o2go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

**A lightweight and flexible OAuth2 core library for Go.**

o2go provides the **core OAuth2 mechanics** (authorization code and refresh token flows)  
without forcing provider-specific logic, storage decisions, or opinionated abstractions.

---

## Why o2go?

Most OAuth2 libraries in the Go ecosystem tend to:

- Be heavy and opinionated
- Mix provider-specific behavior into the core
- Enforce how you manage state, scopes, or tokens
- Hide HTTP details and make debugging harder

**o2go takes a different approach:**

- Minimal and focused on OAuth2 mechanics only
- Provider-agnostic by design
- Full control over parameters and request behavior
- Context-aware requests and custom HTTP clients
- Rich error inspection with access to status codes and response bodies

o2go is **not a framework**  it is a clean and flexible building block.

---

## Installation

```bash
go get github.com/gokhanaltun/o2go
```

---

## Quick Start (Minimal)

```go
package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/gokhanaltun/o2go"
)

type MyProvider struct{}

func (p *MyProvider) AuthURL() string  { return "https://example.com/oauth/authorize" }
func (p *MyProvider) TokenURL() string { return "https://example.com/oauth/token" }

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

	authURL, _ := client.AuthCodeURL("state123")
	fmt.Println("Open:", authURL)

	code := "code-from-callback"
	resp, err := client.ExchangeAuthCode(context.Background(), code)
	if err != nil {
		var he *o2go.HTTPError
		if errors.As(err, &he) {
			fmt.Println("Status:", he.Status)
			fmt.Println("Body:", string(he.Body))
		}
		return
	}

	fmt.Println("Access Token:", resp.Data["access_token"])
}
```

---

## Advanced Usage

### Custom HTTP Client

```go
httpClient := &http.Client{
	Timeout: 10 * time.Second,
}

client, err := o2go.New(provider, cfg, o2go.WithHttpClient(httpClient))
```

If no client is provided, `http.DefaultClient` is used.

---

### Refresh Token Flow

```go
refreshToken := resp.Data["refresh_token"].(string)

newResp, err := client.ExchangeRefreshToken(context.Background(), refreshToken)
if err != nil {
	var he *o2go.HTTPError
	if errors.As(err, &he) {
		fmt.Println("Status:", he.Status)
		fmt.Println("Body:", string(he.Body))
	}
	return
}

fmt.Println("New Access Token:", newResp.Data["access_token"])
```

---

### Provider-Specific Parameters (Scopes, PKCE, Extras)

o2go does not introduce special helpers for extensions like PKCE.

OAuth2 extensions are expressed **purely through parameters**:

```go
cfg.AuthURLParams = map[string]string{
	"scope": "profile email",
	"code_challenge": "...",
	"code_challenge_method": "S256",
}

cfg.TokenParams = map[string]string{
	"code_verifier": "...",
}
```

This keeps the core simple and avoids opinionated abstractions.

---

## Error Handling

All non-2xx HTTP responses are returned as `HTTPError`.

```go
type HTTPError struct {
	Status int
	Body   []byte
	Err    error
}
```

`HTTPError` implements `Unwrap`, making it fully compatible with `errors.As`:

```go
var he *o2go.HTTPError
if errors.As(err, &he) {
	fmt.Println("Status:", he.Status)
	fmt.Println("Body:", string(he.Body))
}
```

---

## Design Principles

* **Core vs Provider separation**
  Core handles OAuth2 mechanics. Provider-specific behavior stays isolated.

* **No opinionated helpers**
  OAuth2 extensions are parameters, not features.

* **Minimal surface area**
  Small API, predictable behavior, easy to reason about.

* **Opt-in providers**
  Use only what you need — nothing more.

---

## Summary

* Lightweight OAuth2 core
* Provider-agnostic design
* Full control over flows and parameters
* Rich, inspectable errors
* No hidden magic

---

## Providers

o2go itself is provider-agnostic. You can either:

- Implement your own provider
- Or use optional, ready-made provider modules (Google, GitHub, etc.) from:  
  **[o2go-providers](https://github.com/gokhanaltun/o2go-providers)**

Provider packages exist for convenience and isolation —  
they do not change or pollute the core.

