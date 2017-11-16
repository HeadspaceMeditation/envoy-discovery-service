FROM scratch
ADD envoy-discovery-service /envoy-discovery-service
ENTRYPOINT ["/envoy-discovery-service"]
