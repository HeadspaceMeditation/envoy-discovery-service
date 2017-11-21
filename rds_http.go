package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type routeHandler struct {
	kubeProxyEndpoint    string
	namespace            string
	serviceLabelSelector string
}

func routeServer(kubeProxyEndpoint string, namespace string, serviceLabelSelector string) http.Handler {
	return &routeHandler{kubeProxyEndpoint, namespace, serviceLabelSelector}
}

func (h *routeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rs, err := getRoutes(h.kubeProxyEndpoint, h.namespace, serviceLabelSelector)
	if err != nil {
		log.Printf("verbose error info: %#v", err)
		log.Println(rs)
		w.WriteHeader(500)
		return
	}

	data, err := json.MarshalIndent(rs, "", "  ")
	if err != nil {
		log.Printf("verbose error info: %#v", err)
		log.Println(rs)
		w.WriteHeader(500)
		return
	}
	w.Write(data)

	elapsed := time.Since(start)
	log.Printf("%s %s", r.URL.Path, elapsed)
}
