/*
	Get a list of all facts.
*/

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func listFacts(c *cli.Context) {
	// Set debug level.
	setDebug(c.GlobalBool("debug"))
	// Start PuppetDB connector.
	startPdbClient(c.GlobalString("puppetdb"))

	resp, err := pdb_client.FactNames()
	if err != nil {
		fmt.Println(err)
	}
	for _, element := range resp {
		fmt.Println(element)
	}
}
