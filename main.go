package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	archive "github.com/antonio-leitao/nau/lib/archive"
	configure "github.com/antonio-leitao/nau/lib/configure"
	home "github.com/antonio-leitao/nau/lib/home"
	new "github.com/antonio-leitao/nau/lib/new"
	open "github.com/antonio-leitao/nau/lib/open"
	utils "github.com/antonio-leitao/nau/lib/utils"
)

// read naurc file and ouput configuration, default if it doesn exist
func readConfig() (utils.Config, error) {
	defaultConfig := utils.Config{
		Name:           "NAU",
		Version:        "0.1.0",
		Author:         "Antonio Leitao",
		Url:            "https://github.com/antonio-leitao/nau",
		Website:        "https://antonio-leitao.github.io",
		Email:          "aleitao@novaims.unl.pt",
		Remote:         "https://github.com/antonio-leitao",
		Base_color:     "#814584",
		Editor:         "nvim",
		Projects_path:  "/Documents/Projects",
		Templates_path: "/Documents/Templates",
		Archives_path:  "/Documents/Archives",
	}

	configFile, err := utils.ConvertPath(".naurc")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return defaultConfig, nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		return utils.Config{}, err
	}
	defer file.Close()

	config := defaultConfig
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return utils.Config{}, fmt.Errorf("Invalid config format")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "AUTHOR":
			config.Author = value
		case "EMAIL":
			config.Email = value
		case "WEBSITE":
			config.Website = value
		case "REMOTE":
			config.Remote = value
		case "EDITOR":
			config.Editor = value
		case "BASE_COLOR":
			config.Base_color = value
		case "PROJECTS_PATH":
			config.Projects_path = value
		case "TEMPLATES_PATH":
			config.Templates_path = value
		case "ARCHIVES_PATH":
			config.Archives_path = value
		default:
			return utils.Config{}, fmt.Errorf("Unknown config key: %s", key)
		}
	}

	if err := scanner.Err(); err != nil {
		return utils.Config{}, err
	}

	return config, nil
}

// function to load the config stuff
func loadConfig() (utils.Config, error) {
	config, err := readConfig()
	if err != nil {
		return utils.Config{}, err
	}
	//get templates
	templatesPath, err := utils.ConvertPath(config.Templates_path)
	if err != nil {
		return utils.Config{}, err
	}
	color_map, err := utils.LoadTemplatesColorMap(templatesPath)
	if err != nil {
		return utils.Config{}, err
	}
	config.Templates = color_map
	return config, nil
}

func main() {
	//read config
	config, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Parse the command-line arguments
	if len(os.Args) < 2 {
		config.Print()
		os.Exit(0)
	}

	command := strings.ToLower(os.Args[1])

	// Launch the appropriate command
	switch command {
	case "open":
		if len(os.Args) < 3 {
			home.Home(config)
		}
		open.Open(config, os.Args[2])
		os.Exit(0)
	case "new":
		if len(os.Args) < 3 {
			new.New(config, "")
			os.Exit(0)
		}
		new.New(config, os.Args[2])
		os.Exit(0)

	case "archive":
		if len(os.Args) < 3 {
			fmt.Printf("TODO: list and choose all projects")
			os.Exit(1)
		}
		archive.Archive(config, os.Args[2])
		os.Exit(0)
	case "config":
		switch len(os.Args) {
		case 1:
			fmt.Println("Error please supply more arguments:", err)
			os.Exit(1)
		case 2:
			configure.Init(config)
			return
		default:
			err = configure.UpdateConfigField(os.Args[2], strings.Join(os.Args[3:], " "))
			if err != nil {
				fmt.Println("Error updating config file file:", err)
				os.Exit(1)
			}
			return
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
