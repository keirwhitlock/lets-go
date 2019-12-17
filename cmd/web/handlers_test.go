package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	/*
		You can indicate that it's OK for a test to be run in concurrently
		alongside other tests by calling the t.Parallel() function at the
		start of the test.

		Tests marked using t.Parallel() will be run in parallel with — and
		only with — other parallel tests.
	*/
	t.Parallel()

	// init a resp writer
	rr := httptest.NewRecorder()

	// init a new request
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// call ping handler func
	ping(rr, r)

	// get response recorded by http.ResponseRecorder
	rs := rr.Result()

	// check resp code
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// check resp body
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}

}
