/*
	Helper functions and structs for better code handling.
*/

package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/temal-/go-puppetdb"
	"os"
	"sort"
)

var (
	pdbClient *puppetdb.Client
	debug     bool
)

func factacularInit(c *cli.Context) {
	// Set debug level.
	setDebug(c.GlobalBool("debug"))
	// Start PuppetDB connector.
	startPdbClient(c.GlobalString("puppetdb"))
}

func setDebug(state bool) {
	debug = state
}

func startPdbClient(nodeName string) {
	pdbClient = puppetdb.NewClient(nodeName)
	checkPuppetAvailability()
}

func checkFactsAvailability(facts []string) (err error) {
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
				return e
			}
		case r := <-refCountChan:
			refCount += r
			if refCount == len(facts) {
				break L
			}
		}
	}
	return

}

func checkFactAvailability(factName string) (err error) {
	facts, err := pdbClient.FactNames()
	for _, fact := range facts {
		if fact == factName {
			return
		}
	}
	return errors.New("\"" + factName + "\" is no valid fact.")
}

func checkPuppetAvailability() {
	pdbVersion, err := pdbClient.PuppetdbVersion()
	if err != nil {
		os.Exit(1)
	}
	if debug {
		fmt.Printf("Using PuppetDB (%s) at: %s\n", pdbVersion, pdbClient.BaseURL)
	}
}

// ValSorter is a helper struct to make the sort of PuppetDB data more easy.
type ValSorter struct {
	Keys []string
	Vals []int
}

// NewValSorter maps a map[string]int to a ValSorter struct.
func NewValSorter(m map[string]int) *ValSorter {
	vs := &ValSorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]int, 0, len(m)),
	}
	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		vs.Vals = append(vs.Vals, v)
	}
	return vs
}

// Sort sorts ValSorter descending.
func (vs *ValSorter) Sort() {
	sort.Sort(sort.Reverse(vs))
}

func (vs *ValSorter) Len() int {
	return len(vs.Vals)
}
func (vs *ValSorter) Less(i, j int) bool {
	return vs.Vals[i] < vs.Vals[j]
}
func (vs *ValSorter) Swap(i, j int) {
	vs.Vals[i], vs.Vals[j] = vs.Vals[j], vs.Vals[i]
	vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}
