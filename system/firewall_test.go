package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRule(t *testing.T) {
	t.Run("Parse tcp", func(t *testing.T) {
		expected := &iptablesRule{
			chain:   "INPUT",
			proto:   "tcp",
			dstPort: 80,
			comment: "",
			target:  "ACCEPT",
		}
		res, err := parseRule("-A INPUT -p tcp -m tcp --dport 80 -c 0 0 -j ACCEPT")
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Parse tcp invalid port", func(t *testing.T) {
		res, err := parseRule("-A INPUT -p tcp -m tcp --dport 2h2 -c 0 0 -j ACCEPT")
		assert.Nil(t, res)
		assert.Error(t, err)
	})

	t.Run("Parse tcp with comment", func(t *testing.T) {
		expected := &iptablesRule{
			chain:   "INPUT",
			proto:   "tcp",
			dstPort: 80,
			comment: "FW-MANAGER RULE",
			target:  "ACCEPT",
		}
		res, err := parseRule(`-A INPUT -p tcp -m tcp --dport 80 -m comment --comment "FW-MANAGER RULE" -c 0 0 -j ACCEPT`)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Parse tcp with comment and source", func(t *testing.T) {
		expected := &iptablesRule{
			chain:   "INPUT",
			proto:   "tcp",
			dstPort: 9100,
			source:  "10.10.0.18/32",
			comment: "FW-MANAGER RULE",
			target:  "ACCEPT",
		}
		res, err := parseRule(`-A INPUT -s 10.10.0.18/32 -p tcp -m tcp --dport 9100 -m comment --comment "FW-MANAGER RULE" -j ACCEPT`)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}
