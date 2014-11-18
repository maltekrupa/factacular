/*
	Get a list of all facts.
*/

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func listFacts(c *cli.Context) {
	// Initialize helpers.
	factacular_init(c)

	resp, err := pdb_client.FactNames()
	if err != nil {
		fmt.Println(err)
	}
	for _, element := range resp {
		fmt.Println(element)
	}
}
