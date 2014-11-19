/*
	Functions for accumulating multiple facts for nodes.
*/

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/temal-/go-puppetdb"
	"log"
	"strings"
	"time"
)

type singleFact struct {
	key   string
	value string
}

type FactsContainer struct {
	node  puppetdb.NodeJson
	facts []singleFact
}

type FactsContainerList []FactsContainer

func (slice FactsContainerList) positionOf(nodeName string) int {
	for k, v := range slice {
		if v.node.Name == nodeName {
			return k
		}
	}
	return -1
}

func (slice FactsContainerList) addFactToNode(factList []puppetdb.FactJson) {
	loc := 0
	for _, v := range factList {
		loc = slice.positionOf(v.CertName)
		if loc < 0 {
			log.Fatal("Weired situation. Got a fact for a node that doesn't exist. Check you PuppetDB!")
		}
		slice[loc].facts = append(slice[loc].facts, singleFact{v.Name, v.Value})
	}
}

func facts(c *cli.Context) {
	if c.Args().First() == "" {
		log.Fatal("Please provide at least one fact.")
	}

	// Initialize helpers.
	factacular_init(c)

	// 'Parse' input and check availability.
	facts := strings.Split(c.Args().First(), ",")
	err := checkFactsAvailability(facts)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all nodes.
	nodes, err := pdb_client.Nodes()
	if err != nil {
		log.Fatal(err)
	}
	// Make some space for the output.
	output := make(FactsContainerList, len(nodes))
	// Put all nodes into the output.
	for k, v := range nodes {
		output[k].node = v
	}

	// Get all facts for all nodes.
	nodeChan := getFactList(facts)
	for {
		select {
		case s := <-nodeChan:
			output.addFactToNode(s)
		case <-time.After(500 * time.Millisecond):
			if debug {
				fmt.Println("Timeout! Printing output.")
			}
			printOutput(output)
			return
		}
	}

}

func printOutput(output FactsContainerList) {
	for foo := range output {
		fmt.Printf("%s | ", output[foo].node.Name)
		for _, v := range output[foo].facts {
			fmt.Printf("%s | ", v.value)
		}
		fmt.Printf("\n")
	}
}

func getFactList(factName []string) <-chan []puppetdb.FactJson {
	c := make(chan []puppetdb.FactJson)
	for _, value := range factName {
		go func(value string) {
			allFacts, err := pdb_client.FactPerNode(value)
			if err != nil {
				log.Fatal(err)
			}
			c <- allFacts
		}(value)
	}
	return c
}
