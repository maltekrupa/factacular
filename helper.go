/*
	Helper functions and structs for better code handling.
*/

package main

import (
	"errors"
	"github.com/codegangsta/cli"
	"github.com/temal-/go-puppetdb"
	"log"
	"sort"
)

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

func checkFactAvailability(c *cli.Context, factName string) (err error) {
	client := puppetdb.NewClient(c.GlobalString("puppetdb"))
	facts, err := client.FactNames()
	for _, fact := range facts {
		if fact == factName {
			return
		}
	}
	return errors.New("\"" + factName + "\" is no valid fact.")
}

func checkPuppetAvailability(c *cli.Context) error {
	client := puppetdb.NewClient(c.GlobalString("puppetdb"))
	_, err := client.PuppetdbVersion()
	if err != nil {
		log.Fatal(err)
	}
	return err
}
