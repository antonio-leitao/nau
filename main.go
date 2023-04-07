package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	open "github.com/antonio-leitao/nau/lib/open"
	structs "github.com/antonio-leitao/nau/lib/structs"
)

// function to load the config stuff
func loadConfig(path string) (structs.Config, error) {
	var config structs.Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return structs.Config{}, err
	}
	return config, nil
}

func main() {
	//read config
	config, err := loadConfig("nau.config.toml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse the command-line arguments
	if len(os.Args) < 2 {
		//home.Home()
		fmt.Printf("TODO: Nau's homepage, a TUI.")
	}

	command := strings.ToLower(os.Args[1])

	// Launch the appropriate command
	switch command {
	case "open":
		if len(os.Args) < 3 {
			fmt.Printf("TODO: list and choose all projects")
			os.Exit(1)
		}
		open.Open(config, os.Args[2])
		os.Exit(1)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
