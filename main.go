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

func validateRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		next(w, r)
	})
}

type ServiceHandler struct {
	Namespace    string
	KubeHost     string
	ServiceLabel string
}

func (s *ServiceHandler) HandleClusters(w http.ResponseWriter, r *http.Request) {
	data, err := makeRequest(servicesPath, s.KubeHost, s.Namespace, s.ServiceLabel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	clusters, err := makeClusters(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshallData(clusters, w, r)
}

func (s *ServiceHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	data, err := makeRequest(endpointsPath, s.KubeHost, s.Namespace, serviceFromURL(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hosts, err := makeHosts(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	marshallData(hosts, w, r)
}

func (s *ServiceHandler) HandleRoutes(w http.ResponseWriter, r *http.Request) {
	data, err := makeRequest(servicesPath, s.KubeHost, s.Namespace, s.ServiceLabel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	routes, err := makeRoutes(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshallData(routes, w, r)
}

func createHandlers(s ServiceHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/clusters/", validateRequest(s.HandleClusters))
	mux.HandleFunc("/v1/registration/", validateRequest(s.HandleRegistration))
	mux.HandleFunc("/v1/routes/", validateRequest(s.HandleRoutes))

	return mux
}

func main() {
	s := ServiceHandler{}

	flag.StringVar(&httpAddr, "http", "127.0.0.1:8080", "The HTTP listen address")
	flag.StringVar(&s.KubeHost, "kube-proxy-endpoint", "127.0.0.1:9090", "A kubectl reverse-proxy URL")
	flag.StringVar(&s.ServiceLabel, "service-label-selector", "envoyTier=ingress", "The label selector to filter services for CDS")
	flag.Parse()

	// POD_NAMESPACE env var should be set in container spec via downward API
	s.Namespace = os.Getenv("POD_NAMESPACE")
	if s.Namespace == "" {
		log.Println("POD_NAMESPACE must be set in environment")
		os.Exit(2)
	}

	log.Println("Starting the Kubernetes Envoy Discovery Service...")
	log.Printf("Listening on %s...", httpAddr)

	mux := createHandlers(s)
	log.Fatal(http.ListenAndServe(httpAddr, mux))
}
