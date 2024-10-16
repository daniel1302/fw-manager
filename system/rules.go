package system

import "github.com/daniel1302/fw-manager/types"

type (
	RuleIP       string
	RulePort     int
	FirewallRule struct {
		IP      RuleIP
		Port    RulePort
		RawRule string
	}
)

const (
	LogstashPort      RulePort = 5141
	NodeExporterPort  RulePort = 9100
	MySQLExportedPort RulePort = 9104
	MySQLPort         RulePort = 3306
)

// Logic to prepare rules is hardcoded as following:
//   - 5141 - Logstash rsyslog port on logs.*, required access from ALL hosts.
//   - 9100 - Node exporter on ALL hosts, required access by metrics.*.
//   - 9104 - MySQL exporter on app.* hosts, required access by metrics.*.
//   - 3306 - MySQL database on app.*, requires access by backups.*.
func PrepareFirewallRules(thisComputer types.FleetItem, fleetCatalog *types.FleetCatalog) []FirewallRule {
	if fleetCatalog == nil {
		return []FirewallRule{}
	}

	result := []FirewallRule{}
	// Allow 9100 for metrics.*
	for _, fleetItem := range (*fleetCatalog)[types.FleetMetrics] {
		if fleetItem.ID == thisComputer.ID {
			continue // Ignore localhost
		}

		result = append(result, FirewallRule{
			IP:   RuleIP(fleetItem.Address),
			Port: NodeExporterPort,
		})
	}

	switch thisComputer.Type {
	case types.FleetLogs:
		// Allow 5141 for all hosts
		for fletType := range *fleetCatalog {
			for _, fleetItem := range (*fleetCatalog)[fletType] {
				if fleetItem.ID == thisComputer.ID {
					continue
				}

				result = append(result, FirewallRule{
					IP:   RuleIP(fleetItem.Address),
					Port: LogstashPort,
				})
			}
		}

	case types.FleetApp:
		// Allow 9104 for all metrics.*
		// Allow 3306 for backups.*
		for _, fleetItem := range (*fleetCatalog)[types.FleetMetrics] {
			if fleetItem.ID == thisComputer.ID {
				continue
			}

			result = append(result, FirewallRule{
				IP:   RuleIP(fleetItem.Address),
				Port: MySQLExportedPort,
			})
		}

		for _, fleetItem := range (*fleetCatalog)[types.FleetBackups] {
			if fleetItem.ID == thisComputer.ID {
				continue
			}

			result = append(result, FirewallRule{
				IP:   RuleIP(fleetItem.Address),
				Port: MySQLPort,
			})
		}
	}

	return result
}
