FROM golang:1.22-bookworm AS builder

COPY . /fw-manager

RUN cd /fw-manager && go build -o /consul-svc  cmd/consul-svc/main.go
RUN cd /fw-manager && go build -o /consul-config-gen  cmd/consul-config-gen/main.go

FROM consul:1.15.4

RUN apk add bash

COPY --from=builder /consul-svc /consul-svc
COPY --from=builder /consul-config-gen /consul-config-gen
COPY ./docker/entrypoint.sh /svc-entrypoint.sh

ENTRYPOINT [ "/svc-entrypoint.sh" ]