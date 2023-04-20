package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// type project
type Project struct {
	Name         string
	Folder_name  string
	Repo_name    string
	Display_Name string
	Code         string
	Lang         string
	Color        string
	Path         string
}

func ToHyphenName(s string) string {
	s = regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(s, "${1}-${2}")
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(s, "-")

	return s
}

func ToDunderName(s string) string {
	s = regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(s, "${1}_${2}")
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(s, "_")
	return s
}

func ToFolderName(s string) string {
	var folderName string
	var capitalize bool

	for i, c := range s {
		if c == '-' || c == '_' {
			capitalize = true
			continue
		}
		if i == 0 || capitalize {
			folderName += strings.ToUpper(string(c))
			capitalize = false
		} else {
			folderName += string(c)
		}
	}

	return folderName
}

func ToDisplayName(s string) string {
	// Replace PascalCase with separate and capitalized
	s = regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(s, "${1} ${2}")

	// Replace hyphens and underscores with spaces
	s = regexp.MustCompile(`[-_]+`).ReplaceAllString(s, " ")

	// Capitalize each word
	words := strings.Fields(s)
	for i, w := range words {
		words[i] = strings.Title(strings.ToLower(w))
	}

	return strings.Join(words, " ")
}

func discombobulate(s string) (string, string, string, string, string) {
	code := s[:3]
	name := s[4:]
	return code, ToDunderName(name), ToFolderName(name), ToHyphenName(name), ToDisplayName(name)
}

// utils get all projects
func getThemedProjects(path string, lang string, color string) []Project {
	var themedProjects []Project
	subentries, _ := os.ReadDir(path + "/" + lang)
	for _, subentry := range subentries {
		if !validEntry(subentry) {
			continue
		}
		code, name, folder_name, repo_name, display_name := discombobulate(subentry.Name())
		project := Project{
			Name:         name,
			Folder_name:  folder_name,
			Repo_name:    repo_name,
			Display_Name: display_name,
			Code:         code,
			Lang:         lang,
			Color:        color,
			Path:         path + "/" + lang + "/" + subentry.Name(),
		}
		themedProjects = append(themedProjects, project)
	}
	return themedProjects

}

func GetProjects(config Config) ([]Project, error) {
	var projectNames []Project
	entries, err := os.ReadDir(config.Projects_path)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !validEntry(entry) {
			continue
		}
		if contains(config.Templates, entry.Name()) {
			themedProjects := getThemedProjects(config.Projects_path, entry.Name(), config.Templates[entry.Name()])
			projectNames = append(projectNames, themedProjects...)
		} else {
			code, name, folder_name, repo_name, display_name := discombobulate(entry.Name())
			project := Project{
				Name:         name,
				Folder_name:  folder_name,
				Repo_name:    repo_name,
				Display_Name: display_name,
				Code:         code,
				Lang:         "Mixed",
				Color:        config.Base_color,
				Path:         config.Projects_path + "/" + entry.Name(),
			}
			projectNames = append(projectNames, project)
		}

	}

	return projectNames, nil
}

// func DimColor(hexColor string, factor float64) (string, error) {
//     // Parse the hexadecimal color string
//     c, err := colorful.Hex(hexColor)
//     if err != nil {
//         return "", err
//     }

//     // Convert the color to the HSL color space
//     h, s, l := c.Hsl()

//     // Dim the color by reducing the saturation and increasing the lightness
//     s *= factor
//     l += (1 - l) * (1 - factor)

//     // Convert the dimmed color back to the RGB color space
//     dimmed := colorful.Hsl(h, s, l)
//     r, g, b := dimmed.RGB255()

//     // Format the RGB values as a hexadecimal color string
//     dimmedHex := fmt.Sprintf("#%02x%02x%02x", r, g, b)

//     return dimmedHex, nil
// }

func contains(color_map map[string]string, key string) bool {
	_, ok := color_map[key]
	if ok {
		return true
	} else {
		return false
	}
}
func validEntry(entry os.DirEntry) bool {
	if !entry.IsDir() {
		return false
	}
	if strings.HasPrefix(entry.Name(), ".") {
		return false
	}
	return true
}

// get all themes
func LoadTemplatesColorMap(dirPath string) (map[string]string, error) {
	// Compile a regular expression to match the folder names
	re := regexp.MustCompile(`^(\w+)_(#\w{6})$`)

	// Initialize a map to store the Name and color
	nameColorMap := make(map[string]string)

	// Walk the directory and process the subdirectories
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a directory and matches the regular expression
		if info.IsDir() && re.MatchString(info.Name()) {
			// Extract the Name and color from the folder name
			match := re.FindStringSubmatch(info.Name())
			if match != nil {
				name := match[1]
				color := match[2]
				// Add the Name and color to the map
				nameColorMap[name] = color
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nameColorMap, nil
}

type Config struct {
	Name           string `toml:"name"`
	Version        int    `toml:"version"`
	Author         string `toml:"author"`
	Url            string `toml:"url"`
	Base_color     string `toml:"base_color"`
	Projects_path  string `toml:"projects_path"`
	Templates_path string `toml:"templates_path"`
	Archives_path  string `toml:"archives_path"`
	Templates      map[string]string
}

func (c Config) Print() {
	fmt.Println("Name:", c.Name)
	fmt.Println("Version:", c.Version)
	fmt.Println("Author:", c.Author)
	fmt.Println("URL:", c.Url)

	keys := make([]string, 0, len(c.Templates))
	for key := range c.Templates {
		keys = append(keys, key)
	}

	fmt.Println("Templates:", c.Templates)
	fmt.Println("Projects_path:", c.Projects_path)
	fmt.Println("Templates_path:", c.Templates_path)
	fmt.Println("Archives_path:", c.Archives_path)

}
