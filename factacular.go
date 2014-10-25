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
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/temal-/go-puppetdb"
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
			Action: func(c *cli.Context) {
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
			},
		},
		{
			Name:      "fact",
			ShortName: "f",
			Usage:     "List fact for all nodes (which have this fact).",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "stats",
					Usage: "Accumulate some stats over all nodes based on this fact.",
				},
				cli.BoolFlag{
					Name:  "nodata",
					Usage: "Outputs only the hostnames which have a value for this fact.",
				},
			},
			Action: func(c *cli.Context) {
				if c.Args().First() == "" {
					fmt.Println("Please provide a fact.")
					return
				}
				fmt.Println("PuppetDB host: " + c.GlobalString("puppetdb"))
				client := puppetdb.NewClient(c.GlobalString("puppetdb"))
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
				} else {
					for _, element := range resp {
						fmt.Printf("%v - %v - %v\n", element.CertName, element.Name, element.Value)
					}
				}
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
