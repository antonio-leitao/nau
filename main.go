package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		Version:        "v0.1.2",
		Url:            "https://github.com/antonio-leitao/nau",
		Author:         "Antonio Leitao",
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
	//gget the number of projects
	project_count, err := countProjects(config)
	if err != nil {
		return utils.Config{}, err
	}
	config.Projects = project_count
	return config, nil
}
func countProjects(config utils.Config) (int, error) {
	projectsPath, err := utils.ConvertPath(config.Projects_path)
	if err != nil {
		return 0, err
	}
	return countSubdirectories(projectsPath, config.Templates)
}
func countSubdirectories(path string, templates map[string]string) (int, error) {
	// Check if the path is a directory
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, err
	}
	count := 0

	for _, file := range fileInfo {
		//if it is not a direcotory
		if !file.IsDir() {
			continue
		}
		//if it is named after a template
		if _, ok := templates[file.Name()]; ok {
			// Recursively count subdirectories of template subdirectory
			subdirPath := filepath.Join(path, file.Name())
			subdirCount, err := countSubdirectories(subdirPath, templates)
			if err != nil {
				return 0, err
			}
			count += subdirCount
		} else {

			count++
		}
	}
	return count, nil
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
		home.Home(config)
		os.Exit(0)
	}

	command := strings.ToLower(os.Args[1])

	// Launch the appropriate command
	switch command {
	case "open":
		if len(os.Args) < 3 {
			open.Expand(config)
			os.Exit(0)
		}
		open.Open(config, os.Args[2])
		os.Exit(0)
	case "goto":
		fmt.Println("Goto: needs to be implemented")
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
		case 2:
			configure.Init(config)
			return
		case 3:
			configure.OutputField(config, os.Args[2])
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
