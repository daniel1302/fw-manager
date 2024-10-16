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
- expose the following ports to your host environment:
    - `8500` - eu-dc
    - `8501` - eu-dc2
    - `8502` - us-dc
    - `8503` - asia-dc



## Binaries

### fw-manager

The main binary in this repository responsible for setting iptables rules based on the business logic for specific infrastructure. The binary uses consul catalog to discover infrastructure.

Flags:

- `--dry-run` - Disables the execution step. Program just prints what rules will be deleted and added.
- `--consul-catalog-file-path` - Specify local file for the consul catalog. If empty catalog will be collected from `https://localhost:8500/...`.
- `--network-cidr` - Specify the network CIDR which this binary will manage.
- `--ip-override` - Useful for testing. If empty the binary will search IP assigned to any local interface that belongs to the network specified in `--network-cidr` subnet.

#### Build

```shell
go build -o ./fw-firewall cmd/fw-manager/main.go
```

#### Usage

With Consul API - Usage on the production

```shell
./fw-manager
```

Locally with docker-compose.yaml - Useful for local development

With consul-api:

```shell
docker-compose up -d

./fw-manager --ip-override 10.10.20.1
```

With local consul catalog
```shell
docker-compose up -d

./fw-manager --ip-override 10.10.0.17 --consul-catalog-file-path", "${workspaceFolder}/services.json",
```

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