package main

import (
	"fmt"
	"github.com/Saturn/saturn-go/cmd"
	"github.com/jessevdk/go-flags"
	"github.com/op/go-logging"
	"os"
	"os/signal"
)

var log = logging.MustGetLogger("main")

type Opts struct {
	Version bool `short:"v" long:"version" description:"Print the version number and exit"`
}

var opts Opts

var VERSION = "0.1.0"

var parser = flags.NewParser(&opts, flags.Default)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Noticef("Received %s\n", sig)
			log.Info("Saturn Server shutting down...")

			os.Exit(1)
		}
	}()

	parser.AddCommand("start",
		"start the Saturn-Server",
		"The start command starts the Saturn-Server",
		&cmd.Start{})
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println(VERSION)
		return
	}
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
