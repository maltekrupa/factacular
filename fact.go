package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/temal-/go-puppetdb"
	"log"
)

func fact(c *cli.Context) {
	if c.Args().First() == "" {
		fmt.Println("Please provide a fact.")
		return
	}

	fmt.Println("PuppetDB host: " + c.GlobalString("puppetdb"))
	client := puppetdb.NewClient(c.GlobalString("puppetdb"))

	// Check if fact is a valid fact.
	err := checkFactAvailability(c, c.Args().First())
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.FactPerNode(c.Args().First())
	if err != nil {
		fmt.Println(err)
	}
	if c.Bool("stats") {
		fmt.Printf("Nodes with fact %s: %d\n", c.Args().First(), len(resp))

		wordCounts := make(map[string]int)
		for _, element := range resp {
			wordCounts[element.Value]++
		}
		vs := NewValSorter(wordCounts)
		vs.Sort()
		for k := range vs.Keys {
			fmt.Printf("%s (%d)\n", vs.Keys[k], vs.Vals[k])
		}
	} else if c.Bool("nodata") {
		for _, element := range resp {
			fmt.Println(element.CertName)
		}
	} else if c.Bool("nofact") {
		// Get a list of all nodes.
		allNodes, _ := client.Nodes()
		// If resp and allNodes have the same length: done.
		if len(allNodes) == len(resp) {
			fmt.Println("All nodes have this fact.")
		} else {
			// Put all nodes in a map (for easy deleting).
			nodesWithoutFact := make(map[string]bool)
			for _, element := range allNodes {
				nodesWithoutFact[element.Name] = true
			}
			// Remove all nodes which provide a valid fact from the map.
			for _, element := range resp {
				if nodesWithoutFact[element.CertName] {
					delete(nodesWithoutFact, element.CertName)
				}
			}
			for name := range nodesWithoutFact {
				fmt.Println(name)
			}
			fmt.Printf("Amount of nodes without this fact: %d\n", len(nodesWithoutFact))
		}
	} else {
		for _, element := range resp {
			fmt.Printf("%v - %v\n", element.CertName, element.Value)
		}
	}
}
