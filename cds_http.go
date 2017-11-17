package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type clusterHandler struct {
	namespace string
}

func clusterServer(namespace string) http.Handler {
	return &clusterHandler{namespace}
}

func (h *clusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// FIXME remove debug logging
	log.Printf(r.URL.Path)

	cs, err := getServices(h.namespace)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := json.MarshalIndent(cs, "", "  ")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}
