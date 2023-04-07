package structs


type Config struct {
    Name              string   `toml:"name"`
    Version           int      `toml:"version"`
    Author            string   `toml:"author"`
    Url               string   `toml:"url"`
    Projects_themes   []string `toml:"projects"`
    Projects_path     string   `toml:"projects_path"`
}
