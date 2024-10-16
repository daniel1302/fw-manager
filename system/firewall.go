package system

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
	"github.com/docker/docker/libnetwork/iptables"
)

const ManagedComment = "FW-MANAGER RULE"

type FirewallManager struct {
	wrapper *iptables.IP6Tables
}

func NewFirewallManager(wrapper *iptables.IPTables) (*FirewallManager, error) {
	if wrapper == nil {
		wrapper, err = iptables.New()
		if err != nil {
			return nil, fmt.Errorf("failed to create iptables wrapper: %w", err)
		}
	}

	
	return &FirewallManager{
		wrapper: 
	}, nil
}

func (fwm *FirewallManager) ListManagedFirewallRules() []FirewallRule {

	return []FirewallRule{}
}

func ExecuteFirewall(rules []FirewallRule) error {

	return nil
}
