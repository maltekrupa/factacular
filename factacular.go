/* factacular.go */
package main

import (
    "os"
    //"github.com/temal-/go-puppetdb"
    "github.com/codegangsta/cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "factacular"
    app.Version = "0.1"
    app.Usage = "Get facts and informations from PuppetDB."
    app.Action = func(c *cli.Context) {
        println("PuppetDB")
    }
    app.Run(os.Args)
}
