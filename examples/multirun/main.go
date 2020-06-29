package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/stephen-fox/gimel"
)

const (
	usage = `Usage: app [options] [path-to-another-program] [program-arguments]
This is an example that demonstrates running another application from memory
one or more times in a row.
Example: app cat /tmp/some-file

Options:`
)

func main() {
	numIterations := flag.Int("i", 2, "The number of times to run the command")
	help := flag.Bool("h", false, "Display this help page")

	flag.Parse()

	if *help {
		fmt.Println(usage)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(flag.CommandLine.Args()) == 0 {
		log.Fatalf("please specify an executable to use and any optional arguments")
	}

	lookedUpApp, err := exec.LookPath(flag.CommandLine.Args()[0])
	if err != nil {
		log.Fatalf("failed to lookup '%s'", flag.CommandLine.Args()[0])
	}

	inMemory, err := gimel.FileFromMemfdCreate("", gimel.MfdCloExec, lookedUpApp)
	if err != nil {
		log.Fatalf("failed to create in memory file - %s", err.Error())
	}
	defer inMemory.Close()

	for i := 1; i < *numIterations+1; i++ {
		cmd := gimel.InMemoryFileToCmd(inMemory, flag.CommandLine.Args()[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		log.Printf("[run %d] executing '%v'...", i, cmd.Args)

		err := cmd.Run()
		if err != nil {
			log.Fatalf("[run %d] failed to run executable from memory - %s",
				i, err.Error())
		}
	}
}
