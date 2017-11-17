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

type Services struct {
	// A list of clusters that will be dynamically added/modified within the cluster manager
	Clusters []Cluster `json:"clusters"`
}

// https://www.envoyproxy.io/docs/envoy/latest/configuration/cluster_manager/cluster.html
type Cluster struct {
	// The name of the cluster which must be unique across all clusters
	Name string `json:"name"`

	// The service discovery type to use for resolving the cluster. Possible options are static, strict_dns, logical_dns, *original_dst*, and sds
	Type string `json:"type"`

	// The timeout for new network connections to hosts in the cluster specified in milliseconds
	ConnectTimeoutMs int32 `json:"connect_timeout_ms"`

	// The load balancer type to use when picking a host in the cluster. Possible options are round_robin, least_request, ring_hash, random, and original_dst_lb
	LBType string `json:"lb_type"`

	// This parameter is required if the service discovery type is sds. It will be passed to the SDS API when fetching cluster members
	ServiceName string `json:"service_name"`

	// A comma delimited list of features that the upstream cluster supports
	Features string `json:"features,omitempty"`
}

func getServices(namespace string) (*Services, error) {
	path := fmt.Sprintf(servicesPath, namespace)
	query := url.Values{}
	query.Set("labelSelector", "tier=microservices")

	r := &http.Request{
		Header: make(http.Header),
		Method: http.MethodGet,
		URL: &url.URL{
			Host:     "127.0.0.1:8001",
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

	clusters := make([]Cluster, 0)

	var cs services
	err = json.NewDecoder(resp.Body).Decode(&cs)
	if err != nil {
		return nil, err
	}

	for _, item := range cs.Items {
		// TODO Extract the literal values below into flags and constants
		clusters = append(clusters, Cluster{
			Name:             item.Metadata.Name,
			Type:             "sds",
			ConnectTimeoutMs: 250,
			LBType:           "round_robin",
			ServiceName:      item.Metadata.Name,
			Features:         "http2",
		})
	}

	return &Services{Clusters: clusters}, nil
}
