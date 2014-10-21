/* factacular.go */
package main

import (
    "fmt"
    "os"
    "sort"
    "github.com/temal-/go-puppetdb"
    "github.com/codegangsta/cli"
)

type ValSorter struct {
    Keys []string
    Vals []int
}

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
    app.Version = "0.2"
    app.Usage = "Get facts and informations from PuppetDB."
    app.Flags = []cli.Flag {
      cli.StringFlag{
        Name: "puppetdb, p",
        Value: "http://localhost:8080",
        Usage: "PuppetDB host.",
        EnvVar: "PUPPETDB_HOST",
      },
    }
    app.Commands = []cli.Command{
        {
            Name:      "list-facts",
            ShortName: "lf",
            Usage:     "List all available facts",
            Action: func(c *cli.Context) {
                fmt.Println("PuppetDB host: " + c.GlobalString("puppetdb"))
                client := puppetdb.NewClient(c.GlobalString("puppetdb"))
                resp, err := client.FactNames()
                if err != nil {
                    fmt.Println(err)
                }
                fmt.Println("Facts: ")
                for _, element := range resp {
                    fmt.Println(element)
                }
            },
        },
        {
            Name:      "list-nodes",
            ShortName: "ln",
            Usage:     "List all available nodes",
            Action: func(c *cli.Context) {
                fmt.Println("PuppetDB host: " + c.GlobalString("puppetdb"))
                client := puppetdb.NewClient(c.GlobalString("puppetdb"))
                resp, err := client.Nodes()
                if err != nil {
                    fmt.Println(err)
                }
                fmt.Println("Nodes: ")
                for _, element := range resp {
                    fmt.Println(element.Name)
                }
            },
        },
        {
            Name:      "node-facts",
            ShortName: "nf",
            Usage:     "List all facts for a specific node.",
            Action: func(c *cli.Context) {
                if(c.Args().First() == "") {
                    fmt.Println("Please provide the FQDN of a node.")
                    return
                }
                fmt.Println("PuppetDB host: " + c.GlobalString("puppetdb"))
                client := puppetdb.NewClient(c.GlobalString("puppetdb"))
                resp, err := client.NodeFacts(c.Args().First())
                if err != nil {
                    fmt.Println(err)
                }
                fmt.Println("Node-facts: ")
                for _, element := range resp {
                    fmt.Printf("%v - %v\n", c.Args().First(), element.Name)
                    fmt.Println(element.Value)
                }
            },
        },
        {
            Name:      "fact",
            ShortName: "f",
            Usage:     "List fact for all nodes.",
            Flags:     []cli.Flag {
                cli.BoolFlag{
                    Name:  "stats",
                    Usage: "Accumulate some stats.",
                },
            },
            Action: func(c *cli.Context) {
                if(c.Args().First() == "") {
                    fmt.Println("Please provide a fact.")
                    return
                }
                fmt.Println("PuppetDB host: " + c.GlobalString("puppetdb"))
                client := puppetdb.NewClient(c.GlobalString("puppetdb"))
                resp, err := client.FactPerNode(c.Args().First())
                if err != nil {
                    fmt.Println(err)
                }
                if(c.Bool("stats")) {
                    fmt.Println("Node-facts")
                    fmt.Printf("Nodes with fact %s: %d\n", c.Args().First(), len(resp))

                    wordCounts := make(map[string]int)
                    for _, element := range resp {
                        wordCounts[element.Value]++
                    }
                    vs := NewValSorter(wordCounts)
                    vs.Sort()
                    for k, _ := range vs.Keys {
                        fmt.Printf("%s (%d)\n", vs.Keys[k], vs.Vals[k])
                    }
                } else {
                    fmt.Println("Fact per node: ")
                    for _, element := range resp {
                        fmt.Printf("%v - %v - %v\n", element.CertName, element.Name, element.Value)
                    }
                }
            },
        },
    }
    app.Action = func(c *cli.Context) {
        fmt.Println("Please provide a command to do stuff. 'h' brings up the help.")
    }
    app.Run(os.Args)
}
