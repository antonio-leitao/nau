package main

import (
	"fmt"
	"os"
	"strings"

	home "github.com/antonio-leitao/nau/lib/home"
	new "github.com/antonio-leitao/nau/lib/new"
)

func main() {
	// Parse the command-line arguments
	if len(os.Args) < 2 {
		home.Home()
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])

	// Launch the appropriate command
	switch command {
	case "new":
		new.New()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
