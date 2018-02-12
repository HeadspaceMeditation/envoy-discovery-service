package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func ip() string {
	return "127.0.0.2"
}

func getPort() int32 {
	return 12345
}

func newEndpoints() endpoints {

	return endpoints{
		Subsets: []subset{
			subset{
				Addresses: []address{
					address{
						IP: ip(),
					},
				},
				Ports: []port{
					port{
						Port: getPort(),
					},
				},
			},
		},
	}
}

func createKubeHostMux(path string, t *testing.T) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		b, err := json.Marshal(newEndpoints())
		if err != nil {
			t.Error(err)
		}
		w.Write(b)
	})
	return mux
}

func TestGetHosts(t *testing.T) {
	var (
		serviceName = "cowService"
	)
	path := fmt.Sprintf(endpointsPath, namespace(), serviceName)
	mux := createKubeHostMux(path, t)

	body := makeKubeRequest(fmt.Sprintf("v1/registration/%v", serviceName), mux, t)
	var s Service
	if err := json.Unmarshal(body, &s); err != nil {
		t.Error(err)
	}

	if len(s.Hosts) != 1 {
		t.Fatal("Incorrect number of hosts returned")
	}

	h := s.Hosts[0]
	expected := Host{
		IPAddress: ip(),
		Port:      getPort(),
	}

	if !reflect.DeepEqual(h, expected) {
		t.Error("Response did not match expected")
	}
}
