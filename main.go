package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	grid "github.com/antonio-leitao/nau/lib/grid"
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
		grid.Grid(config)
		os.Exit(1)
	}

	//command := strings.ToLower(os.Args[1])

	// Launch the appropriate command
	// switch command {
	// case "new":
	// 	new.New()
	// default:
	// 	fmt.Printf("Unknown command: %s\n", command)
	// 	os.Exit(1)
	// }
}
