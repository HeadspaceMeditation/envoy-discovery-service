# Kubernetes Envoy Service Discovery Service

The `envoy-discovery-service` service implements the [Envoy Service Discovery REST API](https://lyft.github.io/envoy/docs/configuration/cluster_manager/sds_api.html) and [Envoy Cluster Discovery REST API](https://www.envoyproxy.io/docs/envoy/latest/configuration/cluster_manager/cds#config-cluster-manager-cds-api) on top of the [Kubernetes Services API](https://kubernetes.io/docs/concepts/services-networking/service).

The variable POD_NAMESPACE should be set in the environment via the [Downward API](https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/#the-downward-api). Each Kubernetes service can then be referenced by its unqualified name.

## Usage

```
kubernetes-envoy-sds -h
```

```
Usage of kubernetes-envoy-sds:
  -http string
    	The HTTP listen address. (default "127.0.0.1:8080")
```
