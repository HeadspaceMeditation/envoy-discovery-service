package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func serveHTTP(start time.Time, result interface{}, err error, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		log.Printf("verbose error info: %#v", err)
		w.WriteHeader(500)
		return
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("verbose error info: %#v", err)
		w.WriteHeader(500)
		return
	}
	w.Write(data)

	elapsed := time.Since(start)
	log.Printf("%s %s", r.URL.Path, elapsed)
}
