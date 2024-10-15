# fw-manager




## Binaries

### consul-config-gen

Simple helper binary used to bootstrap node in docker.

Usage:

```shell
consul-config-gen \
    --data-center "eu-dc1" \
    --data-dir "/consul-home" \
    --node-name "node01.eu-dc.metrics.prod" \
    --retry-join-ip "10.5.0.1" \
    --node-env "metrics" \
    --node-stage "prod"
```

### consul-svg

Binary that simulate consul service.

Usage:

```shell
consul-svc \
    --name wireguard \
    --env "metrics" \
    --stage "prod" \
    --address "localhost" \
    --port 8081
```