package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	PrimaryServiceName = "wireguard"
)

var ErrMissingConsulClient error = fmt.Errorf("consul api client is nil")

type ConsulAPIClient struct {
	client *consulapi.Client
}

func NewConsulAPIClient(client *consulapi.Client) (*ConsulAPIClient, error) {
	if client != nil {
		return &ConsulAPIClient{
			client: client,
		}, nil
	}

	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &ConsulAPIClient{
		client: consul,
	}, nil
}

// GetDataCenters fetches all data centers available in the consul cluster
func (api *ConsulAPIClient) GetDataCenters() ([]string, error) {
	if api.client == nil {
		return nil, ErrMissingConsulClient
	}

	return api.client.Catalog().Datacenters()
}

func (api *ConsulAPIClient) GetFleetCatalog(datacenters []string) ([]*consulapi.CatalogService, error) {
	if api.client == nil {
		return nil, ErrMissingConsulClient
	}

	allServices := []*consulapi.CatalogService{}

	for _, dc := range datacenters {
		resp, _, err := api.client.Catalog().Service(PrimaryServiceName, "", &consulapi.QueryOptions{
			Datacenter: dc,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get catalog for the \"%s\" service in %s DC: %w", PrimaryServiceName, dc, err)
		}

		allServices = append(allServices, resp...)
	}

	return allServices, nil
}
