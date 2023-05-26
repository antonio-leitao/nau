package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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

func contains(color_map map[string]string, key string) bool {
	_, ok := color_map[key]
	if ok {
		return true
	} else {
		return false
	}
}
func ExpandPath(path string) (string, error) {
	//Converts a path starting with ~ into a full path
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}
	return path, nil
}
func validEntry(entry os.DirEntry) bool {
	//Chekcs if a path is not hidden nor a directory
	if !entry.IsDir() {
		return false
	}
	if strings.HasPrefix(entry.Name(), ".") {
		return false
	}
	return true
}
func countProjects(config Config) (int, error) {
	projectsPath, err := ExpandPath(config.Projects_path)
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

// get all themes
func loadTemplatesColorMap(dirPath string) (map[string]string, error) {
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

// utils get all projects
func getThemedProjects(path string, lang string, color string) ([]Project, error) {
	var themedProjects []Project
	subentries, err := os.ReadDir(path + "/" + lang)
	if err != nil {
		return themedProjects, err
	}
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
	return themedProjects, nil

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
    return time.Time{}, fmt.Errorf("Could not read timestamp from %s", dirPath)
}

func GetProjects(config Config) ([]Project, error) {
	var projectNames []Project

	//convert the raw path
	projectPath, err := ExpandPath(config.Projects_path)
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
			themedProjects, err := getThemedProjects(projectPath, entry.Name(), config.Templates[entry.Name()])
            if err != nil{
                return projectNames, err
            }
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

// hidden function
func readProjects()  {}
func readTemplates() {}

// public methods
func GetProject() {
	//gets query. attempts to
}
