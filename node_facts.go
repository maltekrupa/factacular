/*
	Get all facts for a specific node.
*/

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/temal-/go-puppetdb"
)

func nodeFacts(c *cli.Context) {
	if c.Args().First() == "" {
		fmt.Println("Please provide the FQDN of a node.")
		return
	}
	fmt.Println("PuppetDB host: " + c.GlobalString("puppetdb"))
	client := puppetdb.NewClient(c.GlobalString("puppetdb"))
	resp, err := client.NodeFacts(c.Args().First())
	if err != nil {
		fmt.Println(err)
	}
	for _, element := range resp {
		fmt.Printf("%v - %v\n", c.Args().First(), element.Name)
		fmt.Println(element.Value)
	}
}
