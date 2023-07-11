package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

func main() {
	dlevel := flag.StringP("dlevel", "l", "error", "Runs the given string in commandline.")

	flag.Usage = func() {
		fmt.Print("Usage: monkey [flags] [program file] [arguments]\n\nAvailable flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	switch *dlevel {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Hello World! ", args)

}
