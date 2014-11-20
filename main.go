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
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "factacular"
	app.Version = "0.3.2"
	app.Usage = "Get facts and informations from PuppetDB."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "puppetdb, p",
			Value:  "http://localhost:8080",
			Usage:  "PuppetDB host.",
			EnvVar: "PUPPETDB_HOST",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable debug output.",
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
					Name:  "without-data",
					Usage: "Outputs only the nodes which have a value for this fact.",
				},
				cli.BoolFlag{
					Name:  "nofact",
					Usage: "Outputs only the nodes which have NO value for this fact.",
				},
			},
		},
		{
			Name:      "facts",
			ShortName: "fs",
			Usage:     "Accumulate multiple facts for all nodes.",
			Action:    facts,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "nodes",
					Usage: "Show facts only for the provided nodes.",
				},
				cli.BoolFlag{
					Name:  "inflate-facts",
					Usage: "Inflate empty facts fields on all nodes.",
				},
			},
		},
	}
	app.Run(os.Args)
}
