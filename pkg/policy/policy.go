package policy

import (
	"strings"

	cilium_api_v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	"github.com/cilium/cilium/pkg/policy/api"
	"sigs.k8s.io/yaml"
)

type NetworkPolicy struct {
	CiliumNetworkPolicy *CiliumNetworkPolicy
}

type CiliumNetworkPolicy struct {
	cilium_api_v2.CiliumNetworkPolicy
	Spec *Rule `json:"spec,omitempty"`

	Status *cilium_api_v2.NetworkPolicyCondition `json:"status,omitempty"`
}

type Rule struct {
	api.Rule
}

// Convert Network Policy to Yaml
func (np *NetworkPolicy) Yaml() string {
	policy, err := yaml.Marshal(np.CiliumNetworkPolicy)
	if err != nil {
		panic(err)
	}
	// Dirty hack.
	fixedPolicy := strings.ReplaceAll(string(policy), "any:k8s:", "")
	fixedPolicy = strings.ReplaceAll(fixedPolicy, "k8s:io:", "io.")
	return fixedPolicy
}
