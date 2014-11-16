/*
	Get a list of all nodes.
*/

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func listNodes(c *cli.Context) {
	// Initialize helpers
	factacular_init(c)

	resp, err := pdb_client.Nodes()
	if err != nil {
		fmt.Println(err)
	}
	for _, element := range resp {
		fmt.Println(element.Name)
	}
}
