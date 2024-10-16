package system

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/coreos/go-iptables/iptables"
)

const (
	IptablesTableFilter = "filter"
	IptablesChainInput  = "INPUT"

	ManagedComment = "FW-MANAGER RULE"
)

type FirewallManager struct {
	wrapper *iptables.IPTables
}

type iptablesRule struct {
	chain   string
	proto   string
	source  string
	dstPort int
	comment string
	target  string
}

func NewFirewallManager(wrapper *iptables.IPTables) (*FirewallManager, error) {
	if wrapper == nil {
		var err error
		wrapper, err = iptables.New()
		if err != nil {
			return nil, fmt.Errorf("failed to create iptables wrapper: %w", err)
		}
	}

	return &FirewallManager{
		wrapper: wrapper,
	}, nil
}

func (fwm *FirewallManager) ListManagedFirewallRules() ([]FirewallRule, error) {
	rawRules, err := fwm.wrapper.List(IptablesTableFilter, IptablesChainInput)
	if err != nil {
		return nil, fmt.Errorf("failed to list all iptables rules: %w", err)
	}

	result := []FirewallRule{}
	for _, rawRule := range rawRules {
		rule, err := parseRule(rawRule)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rule(%s): %w", rawRule, err)
		}

		if rule.comment == ManagedComment {
			result = append(result, FirewallRule{
				IP:   RuleIP(rule.source),
				Port: RulePort(rule.dstPort),
			})
		}
	}

	return result, nil
}

func (fwm *FirewallManager) ExecuteRules(add []FirewallRule, delete []FirewallRule) error {
	// sudo iptables -D INPUT -p tcp -m tcp --dport 80 -m comment --comment "FW-MANAGER RULE" -j ACCEPT
	for _, rule := range delete {
		err := fwm.wrapper.Delete(
			IptablesTableFilter,
			IptablesChainInput,
			"-p", "tcp",
			"-m", "tcp",
			"--dport", fmt.Sprintf("%d", rule.Port),
			"-s", string(rule.IP),
			"-m", "comment", "--comment", ManagedComment,
			"-j", "ACCEPT",
		)

		if err != nil {
			return fmt.Errorf("failed to delete rule with port %d and user %s: %w", rule.Port, rule.IP, err)
		}
	}

	// sudo iptables -A INPUT -p tcp -m tcp --dport 80 -m comment --comment "FW-MANAGER RULE" -j ACCEPT
	for _, rule := range add {
		err := fwm.wrapper.AppendUnique(
			IptablesTableFilter,
			IptablesChainInput,
			"-p", "tcp",
			"-m", "tcp",
			"--dport", fmt.Sprintf("%d", rule.Port),
			"-s", string(rule.IP),
			"-m", "comment", "--comment", ManagedComment,
			"-j", "ACCEPT",
		)
		if err != nil {
			return fmt.Errorf("failed to add rule with port %d and user %s: %w", rule.Port, rule.IP, err)
		}
	}

	return nil
}

// PrepareRulesExecutionPlan checks existing and new rules and determine which needs to be added and which removed
// It returns rules to delete, rules to add and optionally error
func PrepareRulesExecutionPlan(existingRules []FirewallRule, newRules []FirewallRule) ([]FirewallRule, []FirewallRule, error) {
	rulesToDelete := []FirewallRule{}
	rulesToAdd := []FirewallRule{}

	// Find rules that exists in iptables but they should not be added anymore
	for idx, rule := range existingRules {
		if !slices.ContainsFunc(newRules, func(nR FirewallRule) bool {
			return rule.IP == nR.IP && rule.Port == nR.Port
		}) {
			rulesToDelete = append(rulesToDelete, existingRules[idx])
		}
	}

	// Find rules that should be added but they do not exist in the iptables anymore
	for idx, rule := range newRules {
		if slices.ContainsFunc(existingRules, func(nR FirewallRule) bool {
			return rule.IP == nR.IP && rule.Port == nR.Port
		}) {
			// rule already exists
			continue
		}

		rulesToAdd = append(rulesToAdd, newRules[idx])
	}

	return rulesToDelete, rulesToAdd, nil
}

type tokenT string

func parseRule(rule string) (*iptablesRule, error) {
	const (
		tokenEmpty          tokenT = ""
		tokenChain          tokenT = "chain"
		tokenProto          tokenT = "proto"
		tokenSource         tokenT = "source"
		tokenDstPort        tokenT = "dstPort"
		tokenComment        tokenT = "comment"
		tokenTarget         tokenT = "target"
		tokenCommentContent tokenT = "commentContent"
	)

	ruleSlice := strings.Split(rule, " ")

	result := &iptablesRule{}

	currentToken := tokenEmpty
	for _, part := range ruleSlice {
		switch currentToken {
		// check if current flag is one of the potentially usefull
		case tokenEmpty:
			switch part {
			default:
				currentToken = tokenEmpty
			case "-A", "--append":
				currentToken = tokenChain
			case "-p", "--protocol":
				currentToken = tokenProto
			case "--dport", "--destination-port":
				currentToken = tokenDstPort
			case "--comment":
				currentToken = tokenComment
			case "-j", "--jump":
				currentToken = tokenTarget
			case "-s", "--source":
				currentToken = tokenSource
			}

		case tokenChain:
			result.chain = part
			currentToken = tokenEmpty

		case tokenProto:
			result.proto = part
			currentToken = tokenEmpty

		case tokenDstPort:
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("failed to parse dst port(%s) to int: %w", part, err)
			}

			currentToken = tokenEmpty
			result.dstPort = port

		case tokenComment:
			result.comment = result.comment + " " + part

			// keep adding comment as long as not finished
			if strings.HasSuffix(part, "\"") {
				result.comment = strings.Trim(result.comment, " \"")
				currentToken = tokenEmpty
			}

		case tokenTarget:
			result.target = part
			currentToken = tokenEmpty

		case tokenSource:
			result.source = part
			currentToken = tokenEmpty
		}
	}

	return result, nil
}
