/*
	Get a list of all facts.
*/
package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/temal-/go-puppetdb"
)

func listFacts(c *cli.Context) {
	fmt.Println("PuppetDB host: " + c.GlobalString("puppetdb"))
	client := puppetdb.NewClient(c.GlobalString("puppetdb"))
	resp, err := client.FactNames()
	if err != nil {
		fmt.Println(err)
	}
	for _, element := range resp {
		fmt.Println(element)
	}
}
