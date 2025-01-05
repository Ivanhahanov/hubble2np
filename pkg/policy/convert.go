package policy

import (
	"fmt"
	"hubble2np/pkg/graph"
	"strings"

	v1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	"github.com/cilium/cilium/pkg/policy/api"
	"github.com/urfave/cli/v3"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNetworkPolicy(cmd *cli.Command, n *graph.Node) *NetworkPolicy {
	p := new(NetworkPolicy)
	ciliumPolicy := new(CiliumNetworkPolicy)

	ciliumPolicy.TypeMeta = metav1.TypeMeta{
		APIVersion: "cilium.io/v2",
		Kind:       "CiliumNetworkPolicy",
	}

	ciliumPolicy.ObjectMeta = metav1.ObjectMeta{
		Name:      n.Name,
		Namespace: n.Namespace,
	}

	spec := new(Rule)
	spec.EndpointSelector = api.EndpointSelector{
		LabelSelector: &v1.LabelSelector{
			MatchLabels: labelSliceToMap(n.Labels),
		},
	}
	spec.Ingress = []api.IngressRule{}
	for _, ing := range *n.Ingress {
		spec.AddIngress(cmd, ing)
	}

	spec.Egress = []api.EgressRule{}
	for _, eg := range *n.Egress {
		spec.AddEgress(cmd, eg)
	}
	// add default dns rule
	if !cmd.Bool("nodns") {
		spec.AddDns()
	}

	ciliumPolicy.Spec = spec
	p.CiliumNetworkPolicy = ciliumPolicy
	return p
}

func (r *Rule) AddIngress(cmd *cli.Command, n *graph.Node) *Rule {
	selector := api.EndpointSelector{
		LabelSelector: &v1.LabelSelector{
			MatchLabels: labelSliceToMap(n.Labels),
		},
	}
	r.Ingress = append(r.Ingress, api.IngressRule{
		IngressCommonRule: api.IngressCommonRule{
			FromEndpoints: []api.EndpointSelector{
				selector,
			},
		},
		ToPorts: generatePortRules(cmd, n.Ports),
	})

	return r
}
func (r *Rule) AddEgress(cmd *cli.Command, n *graph.Node) *Rule {
	selector := api.EndpointSelector{
		LabelSelector: &v1.LabelSelector{
			MatchLabels: labelSliceToMap(n.Labels),
		},
	}
	r.Egress = append(r.Egress, api.EgressRule{
		EgressCommonRule: api.EgressCommonRule{
			ToEndpoints: []api.EndpointSelector{
				selector,
			},
		},
		ToPorts: generatePortRules(cmd, n.Ports),
	})

	return r
}

func generatePortRules(cmd *cli.Command, ports []uint32) api.PortRules {
	toPort := api.PortRules{}
	if cmd.Bool("ports") {
		var p = []api.PortProtocol{}
		for _, port := range ports {
			p = append(p, api.PortProtocol{Port: fmt.Sprintf("%d", port)})
		}
		toPort = api.PortRules{
			api.PortRule{Ports: p},
		}
	}
	return toPort
}

func labelSliceToMap(slice []string) map[string]string {
	labels := map[string]string{}
	for _, label := range slice {
		// fmt.Println(l)
		// l = strings.TrimPrefix(l, `k8s:`)
		if !strings.Contains(label, "cilium.k8s") {
			//label = strings.TrimPrefix(label, "k8s:")
			l := strings.Split(label, "=")
			labels[l[0]] = l[1]
		}

	}
	return labels
}
