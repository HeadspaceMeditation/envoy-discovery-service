package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarshall(t *testing.T) {
	t.Parallel()

	type validPayload struct {
		Test string
	}

	data := []struct {
		Payload  interface{}
		Header   int
		TestName string
	}{
		{
			Payload:  validPayload{"cow"},
			Header:   http.StatusOK,
			TestName: "Valid JSON struct",
		},
		{
			Payload:  1234,
			Header:   http.StatusOK,
			TestName: "Valid json number error",
		},
		{
			Payload:  make(chan interface{}),
			Header:   http.StatusInternalServerError,
			TestName: "Invalid json",
		},
	}

	for _, d := range data {
		t.Run(fmt.Sprintf("%v ", d.TestName), func(t *testing.T) {
			r := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "www.google.com", nil)
			if err != nil {
				t.Fatal(fmt.Sprintf("Unable to create request. Error %v", err))
			}
			marshallData(d.Payload, r, request)

			if r.Result().StatusCode != d.Header {
				t.Errorf("Did not receive correct status code. Expected: %v. Got %v", d.Header, r.Result().StatusCode)
			}
		})
	}
}
