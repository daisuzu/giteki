package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "giteki"
	app.Version = Version
	app.Usage = "A tool for Japan certified radio equipment"
	app.Author = "daisuzu"
	app.Email = "daisuzu@gmail.com"
	app.Commands = Commands

	app.Run(os.Args)
}
