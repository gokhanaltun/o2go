package o2go

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HTTPError struct {
	Status int
	Body   []byte
	Err    error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%v (status: %d, body: %s)", e.Err, e.Status, string(e.Body))
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func postForm(ctx context.Context, url string, data url.Values, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, &HTTPError{Err: fmt.Errorf("postForm: failed to create request: %w", err)}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return Do(req, httpClient)
}

func Do(req *http.Request, client *http.Client) ([]byte, error) {
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, &HTTPError{Err: fmt.Errorf("request failed: %w", err)}
	}

	defer resp.Body.Close()

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &HTTPError{Err: fmt.Errorf("failed to read response body: %w", err)}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{
			Status: resp.StatusCode,
			Body:   bodyByte,
			Err:    fmt.Errorf("unexpected status code %d", resp.StatusCode),
		}
	}

	return bodyByte, nil
}
