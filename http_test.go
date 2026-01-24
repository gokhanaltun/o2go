package o2go

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPostForm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("expected Content-Type application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
		}
		r.ParseForm()
		if r.FormValue("foo") != "bar" {
			t.Errorf("expected form value foo=bar, got %s", r.FormValue("foo"))
		}
		fmt.Fprint(w, `{"success": true}`)
	}))
	defer ts.Close()

	data := url.Values{}
	data.Set("foo", "bar")

	resp, err := postForm(context.Background(), ts.URL, data, http.DefaultClient)
	if err != nil {
		t.Fatalf("postForm() error = %v", err)
	}

	want := `{"success": true}`
	if string(resp) != want {
		t.Errorf("postForm() = %s, want %s", string(resp), want)
	}
}

func TestPostForm_ErrorStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request")
	}))
	defer ts.Close()

	_, err := postForm(context.Background(), ts.URL, url.Values{}, http.DefaultClient)
	if err == nil {
		t.Fatal("postForm() expected error for 400 status, got nil")
	}

	hErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("expected *HTTPError, got %T", err)
	}
	if hErr.Status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", hErr.Status)
	}
	if string(hErr.Body) != "bad request" {
		t.Errorf("expected body 'bad request', got %s", string(hErr.Body))
	}
}

func TestHTTPError_Error(t *testing.T) {
	err := &HTTPError{
		Status: 404,
		Body:   []byte("not found"),
		Err:    fmt.Errorf("some error"),
	}
	got := err.Error()
	want := "some error (status: 404, body: not found)"
	if got != want {
		t.Errorf("HTTPError.Error() = %s, want %s", got, want)
	}
}
