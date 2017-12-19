package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func marshallData(result interface{}, w http.ResponseWriter, r *http.Request) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("verbose error info: %#v", err)
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}

func makeRequest(path, kubeProxyEndpoint, namespace, serviceLabelSelector string) ([]byte, error) {
	query := url.Values{}
	// Kinda hacky, but oh well...
	// Doing this becasue we need to add a service name, otherwise
	// it is just a query param
	if endpointsPath == path {
		path = fmt.Sprintf(path, namespace, serviceLabelSelector)
	} else {
		query.Set("labelSelector", serviceLabelSelector)
		path = fmt.Sprintf(path, namespace)
	}

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

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return ioutil.ReadAll(resp.Body)
}
