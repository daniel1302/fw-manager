package consul

import (
	"encoding/json"
	"fmt"
	"os"

	consulapi "github.com/hashicorp/consul/api"
)

// Function read local catalog and parse it's content, catalog may be
// delivered by hand or fetched with curl before the fw-manager binary is executed?
func ReadLocalCatalog(filePath string) ([]*consulapi.CatalogService, error) {
	// todo: think about use streams
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read local catalog file: %w", err)
	}

	res := []*consulapi.CatalogService{}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal catalog rom local file: %w", err)
	}

	return res, nil
}
