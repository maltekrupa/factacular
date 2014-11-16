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
	pdb_client *puppetdb.Client
	debug      bool
)

func startPdbClient(nodeName string) {
	fmt.Println(nodeName)
	pdb_client = puppetdb.NewClient(nodeName)
	fmt.Printf("%+v", pdb_client)
	checkPuppetAvailability()
}

func checkFactAvailability(c *cli.Context, factName string) (err error) {
	facts, err := pdb_client.FactNames()
	for _, fact := range facts {
		if fact == factName {
			return
		}
	}
	return errors.New("\"" + factName + "\" is no valid fact.")
}

func checkPuppetAvailability() {
	pdb_version, err := pdb_client.PuppetdbVersion()
	if err != nil {
		os.Exit(1)
	}
	if debug {
		fmt.Printf("Using PuppetDB (%s) at: %s", pdb_version, pdb_client.BaseURL)
	}
}

func setDebug(state bool) {
	debug = state
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
