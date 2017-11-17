# Kubernetes Envoy Service Discovery Service

The `envoy-discovery-service` service implements the [Envoy Service Discovery REST API](https://lyft.github.io/envoy/docs/configuration/cluster_manager/sds_api.html) and [Envoy Cluster Discovery REST API](https://www.envoyproxy.io/docs/envoy/latest/configuration/cluster_manager/cds#config-cluster-manager-cds-api) on top of the [Kubernetes Services API](https://kubernetes.io/docs/concepts/services-networking/service).

The variable POD_NAMESPACE should be set in the environment via the [Downward API](https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/#the-downward-api). Each Kubernetes service can then be referenced by its unqualified name.

## Usage

```
envoy-discovery-service -h
```

```
Usage of envoy-discovery-service:
  -http string
    	The HTTP listen address (default "127.0.0.1:8080")
  -service-label-selector string
    	The label selector to filter services for CDS (default "envoyTier=ingress")
```

## Building

Local testing:

```bash
$ GOOS=darwin ./build
$ kubectl proxy --port=8001
$ POD_NAMESPACE=alpha ./envoy-discovery-service -http 127.0.0.1:8081
```

Pushing new image:


```bash
$ TAG=X.X.X ./build-container
```
