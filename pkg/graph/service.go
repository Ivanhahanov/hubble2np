package graph

import (
	"fmt"
	"strings"

	"github.com/cilium/cilium/api/v1/flow"
	"github.com/urfave/cli/v3"
)

// Node of the graph with workload fields
type Node struct {
	Name      string   // Workload name
	Namespace string   // Workload namespace
	Labels    []string // Workload labels
	Ingress   *Nodes   // List of ingress workloads
	Egress    *Nodes   // List of egress workloads
	Ports     []uint32 // List of workload ports
}

// Show cli graph
func (n *Node) Show() {
	egresses, ingresses := []string{}, []string{}
	for _, p := range *n.Egress {
		egresses = append(egresses, fmt.Sprintf("%s/%s", p.Namespace, p.Name))
	}
	for _, p := range *n.Ingress {
		ingresses = append(ingresses, fmt.Sprintf("%s/%s", p.Namespace, p.Name))
	}
	fmt.Printf("[%s] -> %s/%s -> [%s]\n", strings.Join(ingresses, ","), n.Namespace, n.Name, strings.Join(egresses, ","))
}

// Add link to node
func (n *Node) AddLink(cmd *cli.Command, endpoint *flow.Endpoint, l4 *flow.Layer4) *Node {
	// TODO: need to test workloads logic
	name := endpoint.Workloads[0].Name

	// get port info from flow
	var port uint32
	switch l4.Protocol.(type) {
	case *flow.Layer4_TCP:
		port = l4.Protocol.(*flow.Layer4_TCP).TCP.DestinationPort
	}

	// add new node if not exists
	svc := findNode(endpoint.Namespace, name, n.Egress)
	if svc == nil {
		svc = &Node{
			Name:      name,
			Namespace: endpoint.Namespace,
			Labels:    endpoint.Labels,
		}
		if port != 0 {
			svc.addPort(port)
		}
		*n.Egress = append(*n.Egress, svc)
	}
	if port != 0 {
		svc.addPort(port)
	}
	return svc
}

// Append port to node if not exists
func (n *Node) addPort(port uint32) *Node {
	for _, p := range n.Ports {
		// if port already exists
		if p == port {
			return n
		}
	}
	n.Ports = append(n.Ports, port)
	return n
}
