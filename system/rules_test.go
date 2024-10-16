package system_test

import (
	"testing"

	"github.com/daniel1302/fw-manager/system"
	"github.com/daniel1302/fw-manager/types"
	"github.com/stretchr/testify/assert"
)

var ExampleFleet types.FleetCatalog = types.FleetCatalog{
	types.FleetMetrics: []types.FleetItem{
		{Type: types.FleetMetrics, ID: "m1", Node: "m1.metrics.prod", Address: "10.10.10.1"},
		{Type: types.FleetMetrics, ID: "m2", Node: "m2.metrics.test", Address: "10.10.10.2"},
	},

	types.FleetApp: []types.FleetItem{
		{Type: types.FleetApp, ID: "s1", Node: "s1.app.prod", Address: "10.10.0.1"},
		{Type: types.FleetApp, ID: "s2", Node: "s2.app.prod", Address: "10.10.0.2"},
		{Type: types.FleetApp, ID: "s3", Node: "s3.app.prod", Address: "10.10.0.3"},
		{Type: types.FleetApp, ID: "s4", Node: "s4.app.prod", Address: "10.10.0.4"},
	},

	types.FleetBackups: []types.FleetItem{
		{Type: types.FleetBackups, ID: "b1", Node: "b1.backups.prod", Address: "10.10.20.1"},
		{Type: types.FleetBackups, ID: "b2", Node: "b2.backups.prod", Address: "10.10.20.2"},
		{Type: types.FleetBackups, ID: "b3", Node: "b3.backups.test", Address: "10.10.20.3"},
	},

	types.FleetLogs: []types.FleetItem{
		{Type: types.FleetLogs, ID: "l1", Node: "l1.Logs.prod", Address: "10.10.30.1"},
		{Type: types.FleetLogs, ID: "l2", Node: "l2.Logs.prod", Address: "10.10.30.2"},
		{Type: types.FleetLogs, ID: "l3", Node: "l3.Logs.test", Address: "10.10.30.3"},
	},
}

// Logic to prepare rules is hardcoded as following:
//   - 5141 - Logstash rsyslog port on logs.*, required access from ALL hosts.
//   - 9100 - Node exporter on ALL hosts, required access by metrics.*.
//   - 9104 - MySQL exporter on app.* hosts, required access by metrics.*.
//   - 3306 - MySQL database on app.*, requires access by backups.*.
func TestPrepareRules(t *testing.T) {
	t.Run("Prepare rules for monitoring server", func(t *testing.T) {
		thisComputer := types.FleetItem{Type: types.FleetMetrics, ID: "m1", Node: "m1.metrics.prod", Address: "10.10.10.1"}

		res := system.PrepareFirewallRules(thisComputer, &ExampleFleet)
		expected := []system.FirewallRule{
			// another metrics servers can access node exported on current computer
			{IP: "10.10.10.2", Port: system.NodeExporterPort},
		}

		assert.ElementsMatch(t, expected, res)
	})

	t.Run("Prepare rules for backups server", func(t *testing.T) {
		thisComputer := types.FleetItem{Type: types.FleetBackups, ID: "b1", Node: "b1.backups.prod", Address: "10.10.20.1"}

		res := system.PrepareFirewallRules(thisComputer, &ExampleFleet)
		expected := []system.FirewallRule{
			// all the metrics servers can access the node exporter running on the current server.
			{IP: "10.10.10.1", Port: system.NodeExporterPort},
			{IP: "10.10.10.2", Port: system.NodeExporterPort},
		}

		assert.ElementsMatch(t, expected, res)
	})

	t.Run("Prepare rules for logs server", func(t *testing.T) {
		thisComputer := types.FleetItem{Type: types.FleetLogs, ID: "l1", Node: "l1.Logs.prod", Address: "10.10.30.1"}

		res := system.PrepareFirewallRules(thisComputer, &ExampleFleet)
		expected := []system.FirewallRule{
			// All metrics server can access node-exporter on the current server
			{IP: "10.10.10.1", Port: system.NodeExporterPort},
			{IP: "10.10.10.2", Port: system.NodeExporterPort},

			// all applications can send logs to the logstash on the current computer
			{IP: "10.10.0.1", Port: system.LogstashPort},
			{IP: "10.10.0.2", Port: system.LogstashPort},
			{IP: "10.10.0.3", Port: system.LogstashPort},
			{IP: "10.10.0.4", Port: system.LogstashPort},

			// all metrics servers can send logs to the logstash on the current computer
			{IP: "10.10.10.1", Port: system.LogstashPort},
			{IP: "10.10.10.2", Port: system.LogstashPort},

			// All backup servers can send logs to logstash on the current computers
			{IP: "10.10.20.1", Port: system.LogstashPort},
			{IP: "10.10.20.2", Port: system.LogstashPort},
			{IP: "10.10.20.3", Port: system.LogstashPort},

			// Other logs servers can send logs to logstash on current computer
			{IP: "10.10.30.2", Port: system.LogstashPort},
			{IP: "10.10.30.3", Port: system.LogstashPort},
		}

		assert.ElementsMatch(t, expected, res)
	})

	t.Run("Prepare rules for apps server", func(t *testing.T) {
		thisComputer := types.FleetItem{Type: types.FleetApp, ID: "s2", Node: "s2.app.prod", Address: "10.10.0.2"}

		res := system.PrepareFirewallRules(thisComputer, &ExampleFleet)
		expected := []system.FirewallRule{
			// All metrics servers can access node-exporter on the current server
			{IP: "10.10.10.1", Port: system.NodeExporterPort},
			{IP: "10.10.10.2", Port: system.NodeExporterPort},

			// All metrics servers can access MYSQLExporter running on this server
			{IP: "10.10.10.1", Port: system.MySQLExportedPort},
			{IP: "10.10.10.2", Port: system.MySQLExportedPort},

			// all backups server can access mysql running on current server
			{IP: "10.10.20.1", Port: system.MySQLPort},
			{IP: "10.10.20.2", Port: system.MySQLPort},
			{IP: "10.10.20.3", Port: system.MySQLPort},
		}

		assert.ElementsMatch(t, expected, res)
	})
}
