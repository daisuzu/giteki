package main

import (
	"log"

	"github.com/codegangsta/cli"
)

var (
	downloadDir = "downloads"
	defaultDb   = "giteki.db"
	defaultBind = ":8000"
)

var Commands = []cli.Command{
	cli.Command{
		Name:   "download",
		Usage:  "Downloads the xls of Japan certified radio equipment list",
		Action: doDownload,
		Flags:  downloadFlags,
	},
	cli.Command{
		Name:   "load",
		Usage:  "Loads the downloaded xls to the database",
		Action: doLoad,
		Flags:  loadFlags,
	},
	cli.Command{
		Name:   "server",
		Usage:  "Runs webserver to display the information in the database",
		Action: doServer,
		Flags:  serverFlags,
	},
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var downloadFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "dst, d",
		Value: downloadDir,
		Usage: "destination directory of the file to download",
	},
	cli.BoolFlag{
		Name:  "all, A",
		Usage: "download all of the files, including the past",
	},
	cli.BoolFlag{
		Name:  "update, U",
		Usage: "overwrite if the file to download is exists",
	},
}

func doDownload(c *cli.Context) {
	err := Download(
		c.String("dst"),
		c.Bool("all"),
		c.Bool("update"),
	)
	assert(err)
}

var loadFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "src, s",
		Value: downloadDir,
		Usage: "source directory of the xls to load",
	},
	cli.StringFlag{
		Name:  "dst, d",
		Value: defaultDb,
		Usage: "destination db to store the data of the xls",
	},
}

func doLoad(c *cli.Context) {
	err := Load(
		c.String("src"),
		c.String("dst"),
	)
	assert(err)
}

var serverFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "bind , b",
		Value: defaultBind,
		Usage: "address and port that bind to sockets",
	},
	cli.StringFlag{
		Name:  "src, s",
		Value: defaultDb,
		Usage: "source db of the equipment information",
	},
}

func doServer(c *cli.Context) {
	err := Server(
		c.String("bind"),
		c.String("src"),
	)
	assert(err)
}
