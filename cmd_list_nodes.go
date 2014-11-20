/*
	Get a list of all nodes.
*/

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func listNodes(c *cli.Context) {
	// Initialize helpers.
	factacularInit(c)

	resp, err := pdbClient.Nodes()
	if err != nil {
		fmt.Println(err)
	}
	for _, element := range resp {
		fmt.Println(element.Name)
	}
}
