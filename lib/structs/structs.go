package structs

import (
	"regexp"
	"strings"
	"time"
)

type Config struct {
    Name              string   `toml:"name"`
    Version           int      `toml:"version"`
    Author            string   `toml:"author"`
    Url               string   `toml:"url"`
    Projects_themes   []string `toml:"projects"`
    Projects_path     string   `toml:"projects_path"`
}

type Project struct {
	Name          string
	Language      string
	Desc          string
	LastModified  time.Time
}

// implement the list.Item interface
func (p Project) FilterValue() string {
	return p.Name + p.Desc
}

// func (p Project) Title() string {
// 	return p.Name
// }

func (p Project) Title() string {
	
    parts := strings.Split(p.Name, "_")
	if len(parts)<2{
		return parts[0]
	}
	var matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	name  := matchAllCap.ReplaceAllString(parts[1], "${1} ${2}")
	return name
}

func (p Project) Description() string {
	return p.Desc
}