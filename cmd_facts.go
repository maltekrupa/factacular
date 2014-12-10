/*
	Functions for accumulating multiple facts for nodes.
*/

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/temal-/go-puppetdb"
	"log"
	"sort"
	"strings"
	"time"
)

type singleFact struct {
	key   string
	value string
}

type multipleFacts []singleFact

/*
multipleFacts is a slice of facts for every node that must be sorted.
TODO: Sort it according to the provided facts. Currently we sort based on the alphabet.
*/
func (mF multipleFacts) Len() int {
	return len(mF)
}

func (mF multipleFacts) Less(i, j int) bool {
	return mF[i].key < mF[j].key
}

func (mF multipleFacts) Swap(i, j int) {
	mF[i].key, mF[j].key = mF[j].key, mF[i].key
	mF[i].value, mF[j].value = mF[j].value, mF[i].value
}

type factsContainer struct {
	node  puppetdb.NodeJson
	facts multipleFacts
}

type factsContainerList []factsContainer

func (slice factsContainerList) positionOf(nodeName string) int {
	for k, v := range slice {
		if v.node.Name == nodeName {
			return k
		}
	}
	return -1
}

func (slice factsContainerList) factAvailableForAllNodes(factName string) bool {
	cnt := 0
E:
	for entry := range slice {
		for fact := range slice[entry].facts {
			if slice[entry].facts[fact].key == factName {
				cnt++
				continue E
			}
		}
	}
	if cnt == len(slice) {
		return true
	}
	return false
}

func (slice factsContainerList) inflateFact(factName string) {
E:
	for entry := range slice {
		for fact := range slice[entry].facts {
			if slice[entry].facts[fact].key == factName {
				continue E
			}
		}
		slice[entry].facts = append(slice[entry].facts, singleFact{factName, "N/A"})
	}
}

func (slice factsContainerList) addFactToNode(factList []puppetdb.FactJson) {
	loc := -1
	for _, v := range factList {
		loc = slice.positionOf(v.CertName)
		if loc < 0 {
			log.Fatal("Weired situation. Got a fact for a node that doesn't exist. Check you PuppetDB!")
		}
		slice[loc].facts = append(slice[loc].facts, singleFact{v.Name, v.Value})
	}
}

func (slice factsContainerList) print() {
	for foo := range slice {
		fmt.Printf("%s | ", slice[foo].node.Name)
		sort.Sort(slice[foo].facts)
		for _, v := range slice[foo].facts {
			fmt.Printf("%s | ", v.value)
		}
		fmt.Printf("\n")
	}
}

func facts(c *cli.Context) {
	if c.Args().First() == "" {
		log.Fatal("Please provide at least one fact.")
	}

	// Initialize helpers.
	factacularInit(c)

	// 'Parse' input and check availability.
	facts := strings.Split(c.Args().First(), ",")
	err := checkFactsAvailability(facts)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all nodes.
	nodes, err := pdbClient.Nodes()
	if err != nil {
		log.Fatal(err)
	}
	// Make some space for the output.
	output := make(factsContainerList, len(nodes))
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
		// Please change this to something else ...
		case <-time.After(500 * time.Millisecond):
			if debug {
				fmt.Println("Timeout! Printing output.")
			}
			if c.Bool("inflate-facts") {
				for _, v := range facts {
					output.inflateFact(v)
				}
			}
			output.print()
			return
		}
	}

}

func getFactList(factName []string) <-chan []puppetdb.FactJson {
	c := make(chan []puppetdb.FactJson)
	for _, value := range factName {
		go func(value string) {
			allFacts, err := pdbClient.FactPerNode(value)
			if err != nil {
				log.Fatal(err)
			}
			c <- allFacts
		}(value)
	}
	return c
}
