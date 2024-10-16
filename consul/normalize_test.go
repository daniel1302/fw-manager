package consul_test

import (
	"encoding/json"
	"testing"

	"github.com/daniel1302/fw-manager/consul"
	"github.com/daniel1302/fw-manager/types"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	consulCatalog := []*consulapi.CatalogService{}
	if err := json.Unmarshal([]byte(consulCatalogData), &consulCatalog); err != nil {
		t.Fatal("failed to unmarshal consulCatalogData in TestNormalize", err)
	}

	t.Run("Normalize example data", func(t *testing.T) {
		expected := types.FleetCatalog{
			types.FleetMetrics: []types.FleetItem{
				{
					Type:    types.FleetMetrics,
					ID:      "b27a1a90-dff4-4ff8-9fe8-cc3b573a85b7",
					Node:    "node-01.eu-dc1.metrics.prod",
					Address: "10.10.0.17",
				},
				{
					Type:    types.FleetMetrics,
					ID:      "03deab88-ddd4-46ca-a38a-e75a4635c3a3",
					Node:    "node-02.eu-dc1.metrics.prod",
					Address: "10.10.0.18",
				},
				{
					Type:    types.FleetMetrics,
					ID:      "16c59e2d-7589-4c87-85a1-6550d7fd6f8c",
					Node:    "node-01.eu-dc1.metrics.test",
					Address: "10.10.0.19",
				},
			},

			types.FleetLogs: []types.FleetItem{
				{
					Type:    types.FleetLogs,
					ID:      "c98551e3-fbda-4b3a-9d83-b2a720150d2e",
					Node:    "node-01.eu-dc1.logs.prod",
					Address: "10.10.0.20",
				},
				{
					Type:    types.FleetLogs,
					ID:      "aa02244b-8015-4d04-b262-3e8dc858f6de",
					Node:    "node-01.eu-dc1.logs.test",
					Address: "10.10.0.22",
				},
			},
			types.FleetBackups: []types.FleetItem{
				{
					Type:    types.FleetBackups,
					ID:      "f2dac58a-4377-4cc2-9fe5-cbc483c82f4f",
					Node:    "node-01.eu-dc1.backups.prod",
					Address: "10.10.0.23",
				},
			},
		}

		normalizedCatalog, err := consul.NormalizeCatalog(consulCatalog)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, normalizedCatalog)
	})

	t.Run("Normalize empty catalog", func(t *testing.T) {
		normalizedCatalog, err := consul.NormalizeCatalog([]*consulapi.CatalogService{})
		assert.NoError(t, err)
		assert.EqualValues(t, normalizedCatalog, types.FleetCatalog{})
	})

	t.Run("Normalize nil catalog", func(t *testing.T) {
		normalizedCatalog, err := consul.NormalizeCatalog(nil)
		assert.NoError(t, err)
		assert.EqualValues(t, normalizedCatalog, types.FleetCatalog{})
	})
}

const consulCatalogData = `[
    {
      "ID": "b27a1a90-dff4-4ff8-9fe8-cc3b573a85b7",
      "Node": "node-01.eu-dc1.metrics.prod",
      "Address": "171.53.121.51",
      "Datacenter": "eu-dc1",
      "NodeMeta": {
        "env": "metrics",
        "stage": "prod"
      },
      "ServiceID": "wireguard",
      "ServiceName": "wireguard",
      "ServiceTags": [
        "eu-dc1",
        "metrics.prod",
        "wireguard",
        "vpn"
      ],
      "ServiceAddress": "10.10.0.17",
      "ServicePort": 51820
    },
    {
      "ID": "03deab88-ddd4-46ca-a38a-e75a4635c3a3",
      "Node": "node-02.eu-dc1.metrics.prod",
      "Address": "221.211.6.3",
      "Datacenter": "eu-dc1",
      "NodeMeta": {
        "env": "metrics",
        "stage": "prod"
      },
      "ServiceID": "wireguard",
      "ServiceName": "wireguard",
      "ServiceTags": [
        "eu-dc1",
        "metrics.prod",
        "wireguard",
        "vpn"
      ],
      "ServiceAddress": "10.10.0.18",
      "ServicePort": 51820
    },
    {
      "ID": "16c59e2d-7589-4c87-85a1-6550d7fd6f8c",
      "Node": "node-01.eu-dc1.metrics.test",
      "Address": "100.93.246.214",
      "Datacenter": "eu-dc1",
      "NodeMeta": {
        "env": "metrics",
        "stage": "test"
      },
      "ServiceID": "wireguard",
      "ServiceName": "wireguard",
      "ServiceTags": [
        "eu-dc1",
        "metrics.test",
        "wireguard",
        "vpn"
      ],
      "ServiceAddress": "10.10.0.19",
      "ServicePort": 51820
    },
    {
      "ID": "c98551e3-fbda-4b3a-9d83-b2a720150d2e",
      "Node": "node-01.eu-dc1.logs.prod",
      "Address": "80.27.3.90",
      "Datacenter": "eu-dc1",
      "NodeMeta": {
        "env": "logs",
        "stage": "prod"
      },
      "ServiceID": "wireguard",
      "ServiceName": "wireguard",
      "ServiceTags": [
        "eu-dc1",
        "logs.prod",
        "wireguard",
        "vpn"
      ],
      "ServiceAddress": "10.10.0.20",
      "ServicePort": 51820
    },
    {
      "ID": "aa02244b-8015-4d04-b262-3e8dc858f6de",
      "Node": "node-01.eu-dc1.logs.test",
      "Address": "113.193.10.185",
      "Datacenter": "eu-dc1",
      "NodeMeta": {
        "env": "logs",
        "stage": "test"
      },
      "ServiceID": "wireguard",
      "ServiceName": "wireguard",
      "ServiceTags": [
        "eu-dc1",
        "logs.test",
        "wireguard",
        "vpn"
      ],
      "ServiceAddress": "10.10.0.22",
      "ServicePort": 51820
    },
    {
      "ID": "f2dac58a-4377-4cc2-9fe5-cbc483c82f4f",
      "Node": "node-01.eu-dc1.backups.prod",
      "Address": "200.224.62.2",
      "Datacenter": "eu-dc1",
      "NodeMeta": {
        "env": "backups",
        "stage": "prod"
      },
      "ServiceID": "wireguard",
      "ServiceName": "wireguard",
      "ServiceTags": [
        "eu-dc1",
        "backups.prod",
        "wireguard",
        "vpn"
      ],
      "ServiceAddress": "10.10.0.23",
      "ServicePort": 51820
    }
]`
