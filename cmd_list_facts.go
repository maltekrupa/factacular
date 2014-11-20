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
	factacularInit(c)

	resp, err := pdbClient.FactNames()
	if err != nil {
		fmt.Println(err)
	}
	for _, element := range resp {
		fmt.Println(element)
	}
}
