package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	s := ServiceHandler{
		Namespace:    namespace(),
		KubeHost:     "anything",
		ServiceLabel: label(),
	}

	discoverServer := httptest.NewServer(createHandlers(s))

	paths := []string{
		"v1/clusters",
		"v1/registration",
		"v1/routes",
	}

	for _, path := range paths {
		t.Run(fmt.Sprintf("%v test", path), func(t *testing.T) {
			u := fmt.Sprintf("%s/%s", discoverServer.URL, path)
			request, err := http.NewRequest("PUT", u, nil)
			if err != nil {
				t.Error(err)
			}

			res, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Error(err)
			}

			if res.StatusCode != http.StatusInternalServerError {
				t.Errorf("Did not receive a %v code back. Got %v", http.StatusInternalServerError, res.StatusCode)
			}
		})

	}

}
