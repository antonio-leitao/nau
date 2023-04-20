package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

//type project
type Project struct {
Name string
Folder_name string
Repo_name string
Code string
Path string
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

//utils to DisplayName

//
func discombobulate(s string) (string, string, string, string) {
    code := s[:3]
    name := s[4:]
    return code, ToDunderName(name), ToFolderName(name), ToHyphenName(name)
}

//utils get all projects
func getThemedProjects(path string, lang string) []Project {
	var themedProjects []Project
	subentries, _ := os.ReadDir(path + "/" + lang)
	for _, subentry := range subentries {
		if !validEntry(subentry) {
			continue
		}
		code, name, folder_name, repo_name := discombobulate(subentry.Name())
		project := Project{
			Name: name,
			Folder_name: folder_name,
			Repo_name: repo_name,
			Code: code,
			Path: path + "/" + lang + "/" + subentry.Name(),
		}
		themedProjects = append(themedProjects, project)
	}
	return themedProjects

}

func GetProjects(path string, projectTypes []string) ([]Project, error) {
	var projectNames []Project
	entries, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !validEntry(entry) {
			continue
		}
		if contains(projectTypes, entry.Name()) {
			themedProjects := getThemedProjects(path, entry.Name())
			projectNames = append(projectNames, themedProjects...)
		} else {
			code, name, folder_name, repo_name := discombobulate(entry.Name())
			project := Project{
				Name: name,
				Folder_name: folder_name,
				Repo_name: repo_name,
				Code: code,
				Path: path + "/" + entry.Name(),
			}
			projectNames = append(projectNames, project)
		}

	}

	return projectNames, nil
}

func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
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




type Config struct {
	Name            string   `toml:"name"`
	Version         int      `toml:"version"`
	Author          string   `toml:"author"`
	Url             string   `toml:"url"`
	Projects_themes []string `toml:"projects"`
	Projects_path   string   `toml:"projects_path"`
	Templates_path  string   `toml:"templates_path"`
}

func (c Config) Print() {
	fmt.Println("Name:", c.Name)
	fmt.Println("Version:", c.Version)
	fmt.Println("Author:", c.Author)
	fmt.Println("URL:", c.Url)

	fmt.Println("Projects_themes:", c.Projects_themes)
	fmt.Println("Projects_path:", c.Projects_path)
	fmt.Println("Tamplates_path:", c.Templates_path)
}
