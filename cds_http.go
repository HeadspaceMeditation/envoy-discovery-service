package main

import (
	"net/http"
	"time"
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
	start := time.Now()
	rs, err := getServices(h.kubeProxyEndpoint, h.namespace, h.serviceLabelSelector)
	serveHTTP(start, rs, err, w, r)
}
