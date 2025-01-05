package policy

import (
	v1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	"github.com/cilium/cilium/pkg/policy/api"
)

// Add default dns rule to allow nslookup requests to dns service
func (r *Rule) AddDns() *Rule {
	selector := api.EndpointSelector{
		LabelSelector: &v1.LabelSelector{
			MatchLabels: map[string]string{
				"k8s:io.kubernetes.pod.namespace": "kube-system",
				"k8s:k8s-app":                     "kube-dns",
			},
		},
	}
	r.Egress = append(r.Egress, api.EgressRule{
		EgressCommonRule: api.EgressCommonRule{
			ToEndpoints: []api.EndpointSelector{
				selector,
			},
		},
		ToPorts: api.PortRules{
			api.PortRule{
				Ports: []api.PortProtocol{
					api.PortProtocol{
						Port:     "53",
						Protocol: api.ProtoUDP,
					},
				},
				Rules: &api.L7Rules{
					DNS: []api.PortRuleDNS{
						api.PortRuleDNS{
							MatchPattern: "*",
						},
					},
				},
			},
		},
	})
	return r
}
