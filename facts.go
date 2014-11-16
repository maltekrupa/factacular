/*
	Functions for accumulating multiple facts for nodes.
*/

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"strings"
)

func facts(c *cli.Context) {
	if c.Args().First() == "" {
		log.Fatal("Please provide at least one fact.")
	}

	// Initialize helpers
	factacular_init(c)

	facts := strings.Split(c.Args().First(), ",")
	if debug {
		fmt.Printf("Gettings values ")
	}

	// Get all facts for all nodes.
	counter := make(chan int)
	for _, value := range facts {
		// TODO: Put this in a function.
		// https://talks.golang.org/2012/concurrency.slide#39
		//go getNodeFacts(value, )
	}
	rets := 0
	for {
		rets += <-counter
		if rets == len(facts) {
			break
		}
	}
	if debug {
		fmt.Println(" done.")
	}

	// Shake the data till it's done.

}

//func getNodeFacts(factName string) puppetdb.NodeJson {
//	_, err := pdb_client.FactPerNode(value)
//	if debug {
//		fmt.Printf(".")
//	}
//	if err != nil {
//		log.Fatal(err)
//	}
//	counter <- 1
//}
