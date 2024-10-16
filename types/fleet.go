package types

import (
	"fmt"
	"net"
	"strings"
)

type FleetType string

const (
	FleetUnknown FleetType = "unknown"
	FleetAll     FleetType = "all"
	FleetLogs    FleetType = "logs"
	FleetMetrics FleetType = "metrics"
	FleetApp     FleetType = "app"
	FleetBackups FleetType = "backups"
)

type FleetItem struct {
	Node    string
	Address string
}

type FleetCatalog map[FleetType][]FleetItem

func (fleet FleetCatalog) FindItemByIP(ip net.IP) *FleetItem {
	for fleetType, fleetItems := range fleet {
		for itemIdx, fleetItem := range fleetItems {
			itemIP := net.ParseIP(fleetItem.Address)
			// invalid ip format
			if itemIP == nil {
				continue
			}

			if ip.Equal(itemIP) {
				return &fleet[fleetType][itemIdx]
			}
		}
	}

	return nil
}

// FleetTagsToFleetType check all the tags assigned to the service and checks if any of them matches to given wildcard:
// `<fleet_type>.*`, e.g: `logs.prod“ -> `logs, `apps.test` -> `apps“
func FleetTagsToFleetType(tags []string) FleetType {
	for _, tag := range tags {
		if strings.HasPrefix(tag, fmt.Sprintf("%s.", FleetLogs)) {
			return FleetLogs
		} else if strings.HasPrefix(tag, fmt.Sprintf("%s.", FleetMetrics)) {
			return FleetMetrics
		} else if strings.HasPrefix(tag, fmt.Sprintf("%s.", FleetApp)) {
			return FleetApp
		} else if strings.HasPrefix(tag, fmt.Sprintf("%s.", FleetBackups)) {
			return FleetBackups
		}
	}

	return FleetUnknown
}
