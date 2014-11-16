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
	name  string
	value string
}

func facts(c *cli.Context) {
	if c.Args().First() == "" {
		log.Fatal("Please provide at least one fact.")
	}

	// Initialize helpers.
	factacular_init(c)

	// 'Parse' input and check availability.
	facts := strings.Split(c.Args().First(), ",")
	factChan := make(chan error)
	refCountChan := make(chan int)
	refCount := 0
	for _, fact := range facts {
		go func(fact string) {
			factChan <- checkFactAvailability(fact)
			refCountChan <- 1
		}(fact)
	}

L:
	for {
		select {
		case e := <-factChan:
			if e != nil {
				log.Fatal(e)
			}
		case r := <-refCountChan:
			refCount += r
			if refCount == len(facts) {
				break L
			}
		}
	}

	output := make(map[string][]singleFact)
	// Get all facts for all nodes.
	nodeChan := getFactList(facts)
	for {
		select {
		case s := <-nodeChan:
			addToOutput(output, s)
		case <-time.After(1 * time.Second):
			if debug {
				fmt.Println("Timeout!")
			}
			printOutput(output)
			return
		}
	}

}

func addToOutput(result map[string][]singleFact, factList []puppetdb.FactJson) {
	for _, value := range factList {
		result[value.CertName] = append(result[value.CertName], singleFact{value.Name, value.Value})
	}
}

func printOutput(output map[string][]singleFact) {
	for key, val := range output {
		fmt.Printf("%s | ", key)
		for _, v := range val {
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
