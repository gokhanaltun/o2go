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
		return nil, &HTTPError{Err: fmt.Errorf("postForm: failed to create request: %v", err)}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, &HTTPError{Err: fmt.Errorf("postForm: request failed: %v", err)}
	}

	defer resp.Body.Close()

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &HTTPError{Err: fmt.Errorf("postForm: failed to read response body: %v", err)}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{
			Status: resp.StatusCode,
			Body:   bodyByte,
			Err:    fmt.Errorf("postForm: unexpected status code %d", resp.StatusCode),
		}
	}

	return bodyByte, nil
}
