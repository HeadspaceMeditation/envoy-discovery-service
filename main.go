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
	"flag"
	"log"
	"net/http"
	"os"
)

const (
	endpointsPath = "/api/v1/namespaces/%s/endpoints/%s"
	servicesPath  = "/api/v1/namespaces/%s/services"
)

var (
	httpAddr string
)

func main() {
	flag.StringVar(&httpAddr, "http", "127.0.0.1:8080", "The HTTP listen address.")
	flag.Parse()

	// POD_NAMESPACE env var should be set in container spec via downward API
	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		log.Println("POD_NAMESPACE must be set in environment")
		os.Exit(2)
	}

	log.Println("Starting the Kubernetes Envoy CDS Service...")
	log.Println("Starting the Kubernetes Envoy SDS Service...")
	log.Printf("Listening on %s...", httpAddr)

	http.Handle("/v1/clusters/", clusterServer(namespace))
	http.Handle("/v1/registration/", registrationServer(namespace))
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}
