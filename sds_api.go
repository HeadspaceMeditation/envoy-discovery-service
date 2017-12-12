// Copyright 2017 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"strings"
)

type Service struct {
	// A list of hosts that make up the service.
	Hosts []Host `json:"hosts"`
}

type Host struct {
	// The IP address of the upstream host.
	IPAddress string `json:"ip_address"`

	// The port of the upstream host.
	Port int32 `json:"port"`

	Tags *Tags `json:"tags,omitempty"`
}

type Tags struct {
	// The optional zone of the upstream host. Envoy uses the zone
	// for various statistics and load balancing tasks.
	AZ string `json:"az,omitempty"`

	// The optional canary status of the upstream host. Envoy uses
	// the canary status for various statistics and load balancing
	// tasks.
	Canary bool `json:"canary,omitempty"`

	// The optional load balancing weight of the upstream host, in
	// the range 1 - 100. Envoy uses the load balancing weight in
	// some of the built in load balancers.
	LoadBalancingWeight int32 `json:"load_balancing_weight,omitempty"`
}

func makeHosts(data []byte) (*Service, error) {
	if data == nil || len(data) == 0 {
		return &Service{Hosts: make([]Host, 0)}, nil
	}

	var eps endpoints
	err := json.Unmarshal(data, &eps)
	if err != nil {
		return nil, err
	}

	// Envoy backends only support a single port. The backend will use the first
	// port found on the Kubernetes endpoint.
	// Open questions around named ports and services with multiple ports.
	subset := eps.Subsets[0]
	hosts := make([]Host, 0)
	for _, address := range subset.Addresses {
		hosts = append(hosts, Host{IPAddress: address.IP, Port: subset.Ports[0].Port})
	}

	return &Service{Hosts: hosts}, nil
}

func serviceFromURL(path string) string {
	s := strings.Split(path, "/")
	if len(s) < 3 {
		return ""
	}
	return s[3]
}
