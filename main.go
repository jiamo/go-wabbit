package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"wabbit-go/model"
)

func main() {
	dlevel := flag.StringP("dlevel", "l", "debug", "Runs the given string in commandline.")

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

	// P1
	p1 := &model.PrintStatement{&model.Integer{42}}
	log.Println(model.NodeAsSource(p1, model.NewContext()))
	var p2 = &model.Statements{
		[]model.Statement{
			&model.PrintStatement{&model.Add{&model.Integer{2}, &model.Integer{3}}},
			&model.PrintStatement{&model.Neg{&model.Integer{5}}},
			&model.PrintStatement{&model.Add{&model.Integer{2}, &model.Mul{&model.Integer{3}, &model.Integer{4}}}},
			&model.PrintStatement{&model.Add{&model.Mul{&model.Integer{2}, &model.Integer{3}}, &model.Integer{4}}},
			&model.PrintStatement{&model.Mul{&model.Grouping{&model.Add{&model.Integer{2}, &model.Integer{3}}}, &model.Integer{4}}},
		},
	}

	log.Println(model.NodeAsSource(p2, model.NewContext()))
}
