package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestGetRoutes(t *testing.T) {
	path := fmt.Sprintf(servicesPath, namespace())
	mux := createKubeServiceMux(path, t)

	body := makeKubeRequest("v1/routes", mux, t)
	var s Routes
	if err := json.Unmarshal(body, &s); err != nil {
		t.Error(err)
	}
	if len(s.VirtualHosts) != 1 {
		t.Fatal("Did not create the correct number of virtual hosts.")
	}

	if len(s.VirtualHosts[0].Routes) != 2 {
		t.Fatal("Did not create the correct number of virtual host routes.")
	}

	prefix := fmt.Sprintf("/%s", itemName())

	vHost := s.VirtualHosts[0]
	if !reflect.DeepEqual(vHost.Domains, []string{"*"}) {
		t.Error("Domains not correctly made")
	}

	if vHost.Name != "egress" {
		t.Error("Vhost name is incorrect.")
	}

	routeHeaders := s.VirtualHosts[0].Routes[0]
	expectedHeaders := Route{
		Cluster:   itemName(),
		Prefix:    prefix,
		TimeoutMS: 0,
		Headers: []Header{Header{
			Name:  "content-type",
			Value: "application/grpc",
		}},
	}

	if !reflect.DeepEqual(routeHeaders, expectedHeaders) {
		t.Error("Expected route header did not match expected")
	}

	routeBasic := s.VirtualHosts[0].Routes[1]
	expected := Route{
		Cluster:   itemName(),
		Prefix:    prefix,
		TimeoutMS: 0,
	}

	if !reflect.DeepEqual(routeBasic, expected) {
		t.Error("Expected route did not match expected")
	}
}
