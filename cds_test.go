package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func itemName() string {
	return "cow"
}

func newItem() services {
	return services{
		Items: []item{
			item{
				Metadata: metadata{
					Name: itemName(),
				},
			},
		},
	}
}

func label() string {
	return "label"
}

func namespace() string {
	return "kubeNamespace"
}

func createKubeServiceMux(path string, t *testing.T) *http.ServeMux {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		s := r.URL.Query().Get("labelSelector")
		if s != label() {
			t.Error("Incorrect label selector param given")
		}

		b, err := json.Marshal(newItem())
		if err != nil {
			t.Error(err)
		}

		w.Write(b)
	})
	return mux
}

func makeKubeRequest(path string, mux *http.ServeMux, t *testing.T) []byte {
	t.Helper()

	kubeServer := httptest.NewServer(mux)
	defer kubeServer.Close()

	kubeURL, err := url.Parse(kubeServer.URL)
	if err != nil {
		t.Error(err)
	}

	s := ServiceHandler{
		Namespace:    namespace(),
		KubeHost:     kubeURL.Host,
		ServiceLabel: label(),
	}

	discoverServer := httptest.NewServer(createHandlers(s))

	u := fmt.Sprintf("%s/%s", discoverServer.URL, path)
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		t.Error(err)
	}

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Did not receive a %v code back. Got %v", http.StatusOK, res.StatusCode)
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	return b
}

func TestGetClusters(t *testing.T) {

	path := fmt.Sprintf(servicesPath, namespace())
	mux := createKubeServiceMux(path, t)

	body := makeKubeRequest("v1/clusters", mux, t)
	var s Services
	if err := json.Unmarshal(body, &s); err != nil {
		t.Error(err)
	}

	if len(s.Clusters) != 1 {
		t.Fatal("Did not create the correct number of clusters.")
	}

	clusters := s.Clusters[0]
	expected := Cluster{
		Name:             itemName(),
		Type:             "sds",
		ConnectTimeoutMs: 250,
		LBType:           "round_robin",
		ServiceName:      itemName(),
		Features:         "http2",
	}

	if !reflect.DeepEqual(clusters, expected) {
		t.Error("Expected cluster to did equal the given")
	}
}
