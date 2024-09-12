package rules

import "github.com/xorcare/pointer"

type Rule struct {
	Name   string
	Kind   string
	Probes []string
	IfNot  []IfNot
}

type IfNot struct {
	SupervisorHealthy *bool
}

func ParseRules(configRules []map[string]interface{}) []*Rule {
	var rules []*Rule
	for _, r := range configRules {
		rule := Rule{
			Name:   r["name"].(string),
			Kind:   r["kind"].(string),
			Probes: make([]string, 0),
			IfNot:  make([]IfNot, 0),
		}
		if probes, ok := r["probes"].([]interface{}); ok {
			for _, p := range probes {
				rule.Probes = append(rule.Probes, p.(string))
			}
		}
		if ifnot, ok := r["ifnot"].([]interface{}); ok {
			for _, i := range ifnot {
				if i.(map[interface{}]interface{})["supervisorHealthy"] != nil {
					rule.IfNot = append(rule.IfNot, IfNot{SupervisorHealthy: pointer.Bool(i.(map[interface{}]interface{})["supervisorHealthy"].(bool))})
				}
			}
		}
		rules = append(rules, &rule)
	}
	return rules
}
