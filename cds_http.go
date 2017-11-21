package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type clusterHandler struct {
	kubeProxyEndpoint    string
	namespace            string
	serviceLabelSelector string
}

func clusterServer(kubeProxyEndpoint string, namespace string, serviceLabelSelector string) http.Handler {
	return &clusterHandler{kubeProxyEndpoint, namespace, serviceLabelSelector}
}

func (h *clusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// FIXME remove debug logging
	log.Printf(r.URL.Path)

	cs, err := getServices(h.kubeProxyEndpoint, h.namespace, serviceLabelSelector)
	if err != nil {
		log.Printf("verbose error info: %#v", err)
		log.Println(cs)
		w.WriteHeader(500)
		return
	}

	data, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		log.Printf("verbose error info: %#v", err)
		log.Println(cs)
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}
