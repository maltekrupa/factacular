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

	app = cli.NewApp()
	app.Name = "factacular"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "puppetdb, p",
			Value:  serverURL.String(),
			Usage:  "PuppetDB host.",
			EnvVar: "PUPPETDB_HOST",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "node-facts",
			ShortName: "nf",
			Usage:     "List all facts for a specific node.",
			Action:    nodeFacts,
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
	mux.HandleFunc("/v3/version",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"version":"2.2.1"}`)
		})

	app.Action = func(c *cli.Context) {
		nodeFacts(c)
		// Output:
		// foobar - architecture
		// amd64
		// foobar - puppetversion
		// 3.7.1
	}
	app.Run([]string{"factacular", "node-facts", "foobar"})
}
