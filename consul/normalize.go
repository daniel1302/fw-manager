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

	result := types.FleetCatalog{
		types.FleetAll: []types.FleetItem{},
	}

	for _, service := range catalog {
		serviceType := types.FleetTagsToFleetType(service.ServiceTags)

		if _, fleetTypeDefined := result[serviceType]; !fleetTypeDefined {
			result[serviceType] = []types.FleetItem{}
		}
		fleetItem := types.FleetItem{
			Node:    service.Node,
			Address: service.ServiceAddress,
		}
		result[types.FleetAll] = append(result[types.FleetAll], fleetItem)
		result[serviceType] = append(result[serviceType], fleetItem)
	}

	return result, nil
}
