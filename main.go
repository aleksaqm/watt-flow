package main

import (
	"flag"
	"fmt"
	"watt-flow/config"
	"watt-flow/server"
)

func main() {
	environment := flag.String("e", "development", "use development configuration")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
	}
	flag.Parse()

	config.Init(*environment)
	server.Init()
}
