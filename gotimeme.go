package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "Name of the timer to collect (required)")
	flag.Parse()

	if name == "" {
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "No arguments found, gotimeme expects to run a command.")
		os.Exit(1)
	}

	args := flag.Args()
	head := args[0]
	tail := args[1:]

	cmd := exec.Command(head, tail...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	tags := []string{"client=gotimeme"}
	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithTags(tags))
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	cmd.Run()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Fprintf(os.Stderr, "Took %s", elapsed)
	statsd.Timing(name, elapsed, []string{}, 1)
}
