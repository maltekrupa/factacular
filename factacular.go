/*
This CLI implements some handy features to query facts from PuppetDB
very conveniently.

Examples:
- Get the fact "os" from all nodes.
`factacular fact os`

- Get a list of all facts.
`factacular list-facts`

- Get all facts from a specific host.
`factacular node-facts fqdn.example.com`
*/
package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/temal-/go-puppetdb"
	"log"
	"os"
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

func main() {
	app := cli.NewApp()
	app.Name = "factacular"
	app.Version = "0.3.1"
	app.Usage = "Get facts and informations from PuppetDB."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "puppetdb, p",
			Value:  "http://localhost:8080",
			Usage:  "PuppetDB host.",
			EnvVar: "PUPPETDB_HOST",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "list-facts",
			ShortName: "lf",
			Usage:     "List all available facts.",
			Action:    listFacts,
		},
		{
			Name:      "list-nodes",
			ShortName: "ln",
			Usage:     "List all available nodes.",
			Action:    listNodes,
		},
		{
			Name:      "node-facts",
			ShortName: "nf",
			Usage:     "List all facts for a specific node.",
			Action:    nodeFacts,
		},
		{
			Name:      "fact",
			ShortName: "f",
			Usage:     "List fact for all nodes (which have this fact).",
			Action:    fact,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "stats",
					Usage: "Accumulate some stats over all nodes based on this fact.",
				},
				cli.BoolFlag{
					Name:  "nodata",
					Usage: "Outputs only the nodes which have a value for this fact.",
				},
				cli.BoolFlag{
					Name:  "nofact",
					Usage: "Outputs only the nodes which have NO value for this fact.",
				},
			},
		},
	}
	app.Run(os.Args)
}

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

func listNodes(c *cli.Context) {
	fmt.Println("PuppetDB host: " + c.GlobalString("puppetdb"))
	client := puppetdb.NewClient(c.GlobalString("puppetdb"))
	resp, err := client.Nodes()
	if err != nil {
		fmt.Println(err)
	}
	for _, element := range resp {
		fmt.Println(element.Name)
	}
}

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
