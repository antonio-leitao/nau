package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mergestat/timediff"
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
	Timestamp    time.Time //Time the proejct was last modified
}

func (p Project) Title() string { return p.Display_Name }
func (p Project) Description() string {
	timestr := timediff.TimeDiff(p.Timestamp, timediff.WithStartTime(time.Now()))
	formatted := fmt.Sprintf("Updated %s", timestr)
	return formatted
}
func (p Project) FilterValue() string { return p.Name + p.Lang }
func (p Project) SubmitValue() string { return p.Path }
func (p Project) GetColor() string    { return p.Color }
func (p Project) GetSubduedColor() string {
	subduedColor, _ := DimColor(p.Color, 0.4)
	return subduedColor
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

		timestamp, _ := getDirectoryTimestamp(path + "/" + lang + "/" + subentry.Name())
		project := Project{
			Name:         name,
			Folder_name:  folder_name,
			Repo_name:    repo_name,
			Display_Name: display_name,
			Code:         code,
			Lang:         lang,
			Color:        color,
			Path:         path + "/" + lang + "/" + subentry.Name(),
			Timestamp:    timestamp,
		}
		themedProjects = append(themedProjects, project)
	}
	return themedProjects

}

func getDirectoryTimestamp(dirPath string) (time.Time, error) {
	// Check if the directory exists
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		return time.Time{}, err
	}

	// If the directory contains a .git file, return the modification time of that file
	gitFilePath := dirPath + "/.git"
	gitFileInfo, err := os.Stat(gitFilePath)
	if err == nil && !gitFileInfo.IsDir() {
		return gitFileInfo.ModTime(), nil
	}

	// If the directory does not contain a .git file, return the creation time of the directory
	if fileInfo.IsDir() {
		return fileInfo.ModTime(), nil
	}

	// If the directory exists but is not a directory, return an error
	return time.Time{}, fmt.Errorf("%s is not a directory", dirPath)
}

func GetProjects(config Config) ([]Project, error) {
	var projectNames []Project

	//convert the raw path
	projectPath, err := ConvertPath(config.Projects_path)
	if err != nil {
		return nil, err
	}
	//read projects directory
	entries, err := os.ReadDir(projectPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !validEntry(entry) {
			continue
		}
		if contains(config.Templates, entry.Name()) {
			themedProjects := getThemedProjects(projectPath, entry.Name(), config.Templates[entry.Name()])
			projectNames = append(projectNames, themedProjects...)
		} else {
			code, name, folder_name, repo_name, display_name := discombobulate(entry.Name())

			timestamp, err := getDirectoryTimestamp(projectPath + "/" + entry.Name())
			if err != nil {
				return nil, err
			}
			project := Project{
				Name:         name,
				Folder_name:  folder_name,
				Repo_name:    repo_name,
				Display_Name: display_name,
				Code:         code,
				Lang:         "Mixed",
				Color:        config.Base_color,
				Path:         projectPath + "/" + entry.Name(),
				Timestamp:    timestamp,
			}
			projectNames = append(projectNames, project)
		}

	}

	return projectNames, nil
}

func DimColor(hexColor string, factor float64) (string, error) {
	// Parse the hexadecimal color string
	c, err := colorful.Hex(hexColor)
	if err != nil {
		return "", err
	}

	// Convert the color to the HSL color space
	h, s, l := c.Hsl()

	// Dim the color by reducing the saturation and increasing the lightness
	s *= (1 - factor)
	l *= (1 - factor)

	// Convert the dimmed color back to the RGB color space
	dimmed := colorful.Hsl(h, s, l)
	r, g, b := dimmed.RGB255()

	// Format the RGB values as a hexadecimal color string
	dimmedHex := fmt.Sprintf("#%02x%02x%02x", r, g, b)

	return dimmedHex, nil
}

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
	Name           string
	Version        string
	Url            string
	Author         string
	Email          string
	Website        string
	Remote         string
	Base_color     string
	Projects_path  string
	Templates_path string
	Archives_path  string
	Editor         string
	Templates      map[string]string
}

func (c Config) Print() {
	fmt.Println("NAME:", c.Name)
	fmt.Println("VERSION:", c.Version)
	fmt.Println("URL:", c.Url)
	fmt.Println("AUTHOR:", c.Author)
	fmt.Println("EMAIL:", c.Email)
	fmt.Println("REMOTE:", c.Remote)
	fmt.Println("EDITOR:", c.Editor)

	keys := make([]string, 0, len(c.Templates))
	for key := range c.Templates {
		keys = append(keys, key)
	}

	fmt.Println("TEMPLATES:", keys)
	fmt.Println("PROJECTS_PATH:", c.Projects_path)
	fmt.Println("TEMPLATES_PATH:", c.Templates_path)
	fmt.Println("ARCHIVES_PATH:", c.Archives_path)

}
func ConvertPath(path string) (string, error) {
	// Get the current user's home directory
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	// Expand the tilde symbol to the full path of the home directory
	expandedPath := filepath.Join(usr.HomeDir, path)
	// Return the absolute path
	return filepath.Abs(expandedPath)
}
