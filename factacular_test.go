package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"net/http/httptest"
	"net/url"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	app    *cli.App
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	serverURL, _ := url.Parse(server.URL)

	mux.HandleFunc("/v3/version",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"version":"2.2.1"}`)
		})

	app = cli.NewApp()
	app.Name = "factacular"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "puppetdb, p",
			Value:  serverURL.String(),
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
	}
}

func teardown() {
	server.Close()
}

func ExampleNodeFacts() {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/nodes/foobar/facts",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{
								"value"    : "amd64",
								"name"     : "architecture",
								"certname" : "foobar"
						 	},
							{ 
								"value"    : "3.7.1",
								"name"     : "puppetversion",
								"certname" : "foobar"
						 	}]`)
		})

	app.Action = func(c *cli.Context) {
		nodeFacts(c)
		// Output:
		// foobar - architecture
		// amd64
		// foobar - puppetversion
		// 3.7.1
	}
	app.Run([]string{"factacular", "--debug", "node-facts", "foobar"})
}

func ExampleListFacts() {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/fact-names",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[ "architecture", "os", "puppetversion" ]`)
		})

	app.Action = func(c *cli.Context) {
		listFacts(c)
		// Output:
		// architecture
		// os
		// puppetversion
	}
	app.Run([]string{"factacular", "list-facts"})
}

func ExampleListNodes() {
	setup()
	defer teardown()

	mux.HandleFunc("/v3/nodes",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{ "name" : "fqdn1.example.com",
								"deactivated" : null,
								"catalog_timestamp" : "2014-11-07T23:32:09.998Z",
								"facts_timestamp" : "2014-11-07T23:32:04.723Z",
								"report_timestamp" : "2014-11-07T23:32:10.372Z"
							}, {
								"name" : "fqdn2.example.com",
								"deactivated" : null,
								"catalog_timestamp" : "2014-11-08T08:09:12.544Z",
								"facts_timestamp" : "2014-11-08T08:09:06.224Z",
								"report_timestamp" : "2014-11-08T08:09:16.779Z"
							}, {
								"name" : "fqdn1.example.org",
								"deactivated" : null,
								"catalog_timestamp" : "2014-11-08T06:07:10.296Z",
								"facts_timestamp" : "2014-11-08T06:07:04.789Z",
								"report_timestamp" : "2014-11-08T06:07:10.627Z"
							}]`)
		})

	app.Action = func(c *cli.Context) {
		listNodes(c)
		// Output:
		// fqdn1.example.com
		// fqdn2.example.com
		// fqdn1.example.org
	}
	app.Run([]string{"factacular", "list-nodes"})
}
