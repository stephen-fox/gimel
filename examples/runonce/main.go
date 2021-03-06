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
This is an example that demonstrates running another application from memory.
Example: app cat /tmp/some-file

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

	if len(flag.CommandLine.Args()) == 0 {
		log.Fatalf("please specify an executable to use and any optional arguments")
	}

	lookedUpApp, err := exec.LookPath(flag.CommandLine.Args()[0])
	if err != nil {
		log.Fatalf("failed to lookup '%s'", flag.CommandLine.Args()[0])
	}

	cmd, inMemoryFile, err := gimel.MemfdCreateFromExe(
		"",
		lookedUpApp,
		flag.CommandLine.Args()[1:]...)
	if err != nil {
		log.Fatalf("failed to copy file and setup - %s", err.Error())
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("executing: '%v'...", cmd.Args)

	err = cmd.Run()
	inMemoryFile.Close()
	if err != nil {
		log.Fatalf("failed to run process from memory - %s", err.Error())
	}
}
