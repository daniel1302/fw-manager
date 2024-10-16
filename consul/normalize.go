package consul

import (
	"github.com/daniel1302/fw-manager/types"
	consulapi "github.com/hashicorp/consul/api"
)

// NormalizeCatalog prepares fleet catalog in more friendly for later operations format.
func NormalizeCatalog(catalog []*consulapi.CatalogService) (types.FleetCatalog, error) {
	if len(catalog) < 1 {
		return types.FleetCatalog{}, nil // Nothing to do?
	}

	result := types.FleetCatalog{}

	for _, service := range catalog {
		serviceType := types.FleetTagsToFleetType(service.ServiceTags)

		if _, fleetTypeDefined := result[serviceType]; !fleetTypeDefined {
			result[serviceType] = []types.FleetItem{}
		}
		result[serviceType] = append(result[serviceType], types.FleetItem{
			Type:    serviceType,
			ID:      service.ID,
			Node:    service.Node,
			Address: service.ServiceAddress,
		})
	}

	return result, nil
}
