package graph

import (
	"log"

	"github.com/cilium/cilium/api/v1/flow"
	"github.com/urfave/cli/v3"
)

type Nodes []*Node

func NewGraph() *Nodes {
	return new(Nodes)
}

func (g *Nodes) Show() {
	for _, svc := range *g {
		svc.Show()
	}
}

func (s *Nodes) AddNode(cmd *cli.Command, endpoint *flow.Endpoint) *Node {
	if len(endpoint.Workloads) > 1 {
		log.Println("DEBUG:", endpoint.String())
	}
	name := endpoint.Workloads[0].Name
	svc := findNode(endpoint.Namespace, name, s)
	if svc == nil {
		svc = &Node{
			Name:      name,
			Namespace: endpoint.Namespace,
			Labels:    endpoint.Labels,
			Egress:    new(Nodes),
			Ingress:   new(Nodes),
		}
		*s = append(*s, svc)
	}
	return svc
}

func (s *Nodes) GenerateIngresses() {
	for _, svc := range *s {
		s.generateIngressesForService(svc.Namespace, svc.Name)
	}
}

func (s *Nodes) generateIngressesForService(ns, name string) {
	var ingresses = new(Nodes)
	for _, svc := range *s {
		for _, p := range *svc.Egress {
			if p.Namespace == ns && p.Name == name {
				ingSvc := findNode(ns, name, ingresses)
				if ingSvc == nil {
					svc.Ports = p.Ports
					*ingresses = append(*ingresses, svc)
				}
			}
		}
	}
	for i, svc := range *s {
		if svc.Namespace == ns && svc.Name == name {
			(*s)[i].Ingress = ingresses
		}
	}
}

func findNode(ns, name string, nodes *Nodes) *Node {
	for _, svc := range *nodes {
		if svc.Name == name && svc.Namespace == ns {
			return svc
		}
	}
	return nil
}
