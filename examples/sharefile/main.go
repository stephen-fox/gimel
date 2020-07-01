package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/stephen-fox/gimel"
)

const (
	usage = `Usage: app [options] [path-to-a-file]
This is an example that demonstrates loading a file into memory, and then
sharing it by file path.
Example: app /tmp/some-file

Options:`
)

func main() {
	help := flag.Bool("h", false, "Display this help page")

	flag.Parse()

	if *help {
		fmt.Println(usage)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(flag.Args()) != 1 {
		log.Fatalf("please specify a file to load into memory")
	}

	inMemory, err := gimel.MemfdCreateFromFile("", 0, flag.Arg(0))
	if err != nil {
		log.Fatalf("failed to load file into memory - %s", err.Error())
	}

	log.Printf("file can be accessed as '%s' (press Control+C to exit)", inMemory.Name())

	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	<-interrupts
	inMemory.Close()
}
