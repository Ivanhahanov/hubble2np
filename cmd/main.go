package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"hubble2np/pkg/graph"
	"hubble2np/pkg/policy"
	"log"
	"os"

	observerpb "github.com/cilium/cilium/api/v1/observer"
	"github.com/urfave/cli/v3"
)

var g = graph.NewGraph()

func main() {
	cmd := &cli.Command{
		Name:  "hubble2np",
		Usage: "generate Cilium Network Policies from Hubble flow",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "nodns",
				Usage: "disable dns",
			},
			&cli.BoolFlag{
				Name:    "ports",
				Aliases: []string{"p"},
				Usage:   "enable ports",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := ReadStdin(cmd); err != nil {
				return err
			}
			for _, n := range *g {
				fmt.Println("---")
				fmt.Println(policy.CreateNetworkPolicy(cmd, n).Yaml())
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "graph",
				Usage: "Show graph",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if err := ReadStdin(cmd); err != nil {
						return err
					}

					g.Show()
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// read flow from stdin
func ReadStdin(cmd *cli.Command) error {
	// create stdin scanner
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var fl = observerpb.GetFlowsResponse{}
		json.Unmarshal(scanner.Bytes(), &fl)
		flow := fl.GetFlow()

		if flow != nil {
			// filter flow
			if flow.Source.Workloads != nil &&
				flow.Destination.Workloads != nil &&
				flow.IsReply.Value == false {
				// add flow data to graph
				g.AddNode(cmd, flow.Source).AddLink(cmd, flow.Destination, flow.L4)
			}
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	g.GenerateIngresses()
	return nil
}
