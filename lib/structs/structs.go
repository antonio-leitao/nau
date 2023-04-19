package structs

import "fmt"

type Config struct {
	Name            string   `toml:"name"`
	Version         int      `toml:"version"`
	Author          string   `toml:"author"`
	Url             string   `toml:"url"`
	Projects_themes []string `toml:"projects"`
	Projects_path   string   `toml:"projects_path"`
}

func (c Config) Print() {
	fmt.Println("Name:", c.Name)
	fmt.Println("Version:", c.Version)
	fmt.Println("Author:", c.Author)
	fmt.Println("URL:", c.Url)

	fmt.Println("Projects_themes:", c.Projects_themes)
	fmt.Println("Projects_path:", c.Projects_path)
}
