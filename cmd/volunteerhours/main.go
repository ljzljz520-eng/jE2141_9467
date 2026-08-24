package main

import (
	"fmt"
	"os"

	"volunteerhours/cli"
	"volunteerhours/storage"
)

func main() {
	path := os.Getenv("VOLUNTEER_HOURS_DB")
	if path == "" {
		path = "volunteer-hours.db"
	}
	store, err := storage.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()
	runner := cli.NewRunner(os.Stdin, os.Stdout, store)
	if err := runner.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
