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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// https://www.envoyproxy.io/docs/envoy/latest/configuration/http_conn_man/route_config/route_config.html
type Routes struct {
	// An array of virtual hosts that make up the route table
	VirtualHosts []VirtualHost `json:"virtual_hosts"`
}

// https://www.envoyproxy.io/docs/envoy/latest/configuration/http_conn_man/route_config/vhost.html
type VirtualHost struct {
	// TThe logical name of the virtual host
	Name string `json:"name"`

	// A list of domains (host/authority header) that will be matched to this virtual host
	Domains []string `json:"domains"`

	// The list of routes that will be matched, in order, for incoming requests
	Routes []Route `json:"routes"`
}

// https://www.envoyproxy.io/docs/envoy/latest/configuration/http_conn_man/route_config/route.html
type Route struct {
	// The upstream cluster to which the request should be forwarded to
	Cluster string `json:"cluster"`
	// Specifies a set of headers that the route should match on
	Headers []Header `json:"headers,omitempty"`
	// If specified, the route is a prefix rule meaning that the prefix must match the beginning of the :path header
	Prefix string `json:"prefix"`
	// Specifies the timeout for the route
	TimeoutMS int32 `json:"timeout_ms"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func getRoutes(kubeProxyEndpoint string, namespace string, serviceLabelSelector string) (*Routes, error) {
	path := fmt.Sprintf(servicesPath, namespace)
	query := url.Values{}
	query.Set("labelSelector", serviceLabelSelector)

	r := &http.Request{
		Header: make(http.Header),
		Method: http.MethodGet,
		URL: &url.URL{
			Host:     kubeProxyEndpoint,
			Path:     path,
			Scheme:   "http",
			RawQuery: query.Encode(),
		},
	}

	r.Header.Set("Accept", "application/json, */*")

	ctx := context.Background()
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var svcs services
	err = json.NewDecoder(resp.Body).Decode(&svcs)
	if err != nil {
		return nil, err
	}

	routes := make([]Route, 0)

	for _, item := range svcs.Items {
		// TODO Extract the literal values below into flags and constants
		prefix := fmt.Sprintf("/%s", item.Metadata.Name)
		// TODO Not every service will have both an http and grpc endpoint
		routes = append(routes, []Route{Route{
			Cluster:   item.Metadata.Name,
			Prefix:    prefix,
			TimeoutMS: 0,
			Headers: []Header{Header{
				Name:  "content-type",
				Value: "application/grpc",
			}},
		}, Route{
			Cluster:   item.Metadata.Name,
			Prefix:    prefix,
			TimeoutMS: 0,
		},
		}...)
	}

	return &Routes{VirtualHosts: []VirtualHost{VirtualHost{
		Domains: []string{"*"},
		Name:    "egress",
		// TODO Not every service will have both an http and grpc endpoint
		Routes: routes}},
	}, nil
}
