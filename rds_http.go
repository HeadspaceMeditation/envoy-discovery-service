package main

import (
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
	rs, err := getRoutes(h.kubeProxyEndpoint, h.namespace, serviceLabelSelector)
	serveHTTP(start, rs, err, w, r)
}
