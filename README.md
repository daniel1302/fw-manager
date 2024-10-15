# fw-manager

## Replicate multi-cluster consul infrastructure with docker-compose.yaml

```shell
# If you are experimenting you may want to remove consul homes and recreate clusters from the scratch.
rm -rf ./consul-home

docker-compose up -d
```

It may take up to 1 minute to startup entire infrastructure as there are some dependencies in the containers.

The above will:
    - start the multiple consul servers in "different data centers"
    - start multiple agents simulating real vms with registered service
    - Expose the following ports to your host environment:
        - `8500` - eu-dc
        - `8501` - eu-dc2
        - `8502` - us-dc
        - `8503` - asia-dc



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