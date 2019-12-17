package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	/*
		You can indicate that it's OK for a test to be run in concurrently
		alongside other tests by calling the t.Parallel() function at the
		start of the test.

		Tests marked using t.Parallel() will be run in parallel with — and
		only with — other parallel tests.
	*/
	t.Parallel()
	// init a new resp writer
	rr := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// mock HTTP handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// pass the mock handler to secureHeaders.
	// as secureHeaders returns a ServeHTTP() method, we can call it here
	// passing in our test rr & r
	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1, mode=block", xssProtection)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
