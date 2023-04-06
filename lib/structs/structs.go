package structs

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
}

// implement the list.Item interface
func (p Project) FilterValue() string {
	return p.Name + p.Desc
}

func (p Project) Title() string {
	return p.Name
}

func (p Project) Description() string {
	return p.Desc
}