package new

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	utils "github.com/antonio-leitao/nau/lib/utils"
	"github.com/sahilm/fuzzy"
)

type Data struct {
	Author      string
	Email       string
	Repo        string
	Git         bool
	Name        string
	Description string
}

func loggit(msg string) error {
	// Open the log file for writing, creating it if it doesn't exist
	f, err := os.OpenFile("logger.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// Write the message to the log file
	_, err = fmt.Fprintln(f, msg)
	if err != nil {
		return err
	}
	return nil
}

func createNewProject(sub Submission, config *utils.Config, template string) {
	//if the project is empty just start an empty one
	if template == "Empty" {
		newEmptyProject(sub, config)
	} else {
		createTemplateProject(sub, config, template)
	}
}

func createTemplateProject(sub Submission, config *utils.Config, template string) {
	//new_data_to_colapse
	data := Data{
		Author:      config.Author,
		Email:       config.Email,
		Repo:        config.Remote + "/" + sub.repo_name,
		Git:         sub.git,
		Name:        sub.project_name,
		Description: sub.description,
	}
	//convert source path
	source_path, _ := utils.ConvertPath(config.Templates_path)
	source_path = config.Templates_path + "/" + template + "_" + config.Templates[template]
	//convert target_path
	target_path, _ := utils.ConvertPath(config.Projects_path)
	target_path = target_path + "/" + template + "/" + sub.folder_name
	//create new direcotry
	_ = os.MkdirAll(target_path, 0755)
	//place everything there
	err := CopyDirectory(source_path, target_path, &data)
	if err != nil {
		loggit(err.Error())
		fmt.Printf("NAU ERROR: Failed to create Template: %v", err)
		return
	}
	//colapse template
	err = CollapseDirectory(target_path, &data)
	if err != nil {
		loggit(err.Error())
		fmt.Println("NaError processing directory:", err)
		return
	}
}

func newEmptyProject(sub Submission, config *utils.Config) {
	target_path, _ := utils.ConvertPath(config.Projects_path)
	target_path = target_path + "/"
	err := createEmptyFolder(target_path, sub.folder_name)
	if err != nil {
		fmt.Printf("NAU ERROR: Failed to create folder: %v", err)
		return
	}
}

func createEmptyFolder(target, name string) error {
	// Create the full target path by joining the target directory and the name
	targetPath := filepath.Join(target, name)

	// Create the target directory
	err := os.MkdirAll(targetPath, 0755)
	if err != nil {
		return fmt.Errorf("NAU ERROR: failed to create empty directory: %w", err)
	}

	return nil
}

func CopyDirectory(scrDir, dest string, data *Data) error {
	entries, err := os.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())

		//check any file names
		folder_name, err := collapseString(entry.Name(), data)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, folder_name)

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath, data); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		fInfo, err := entry.Info()
		if err != nil {
			return err
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, fInfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func CollapseDirectory(directory string, data *Data) error {
	nauPath := filepath.Join(directory, ".nau")
	_, err := os.Stat(nauPath)
	if err != nil {
		if os.IsNotExist(err) {
			// The ".nau" file doesn't exist, so we're done
			return nil
		}
		return err
	}
	nauBytes, err := ioutil.ReadFile(nauPath)
	if err != nil {
		return err
	}

	nauLines := strings.Split(string(nauBytes), "\n")
	for _, line := range nauLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		//colapse filename
		line, err = collapseString(line, data)
		if err != nil {
			return err
		}
		matches, err := filepath.Glob(filepath.Join(directory, line))
		if err != nil {
			return err
		}
		for _, match := range matches {
			if fi, err := os.Stat(match); err == nil && fi.Mode().IsRegular() {
				// The match is a regular file
				newContent, err := ProcessFile(match, data)
				if err != nil {
					return err
				}
				err = ioutil.WriteFile(match, []byte(newContent), 0644)
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
		}
	}
	// Delete the .nau file
	err = os.Remove(nauPath)
	if err != nil {
		return err
	}
	return nil
}

func ProcessFile(filename string, data *Data) (string, error) {
	contentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	content := string(contentBytes)
	tmpl, err := template.New(filename).Parse(content)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func collapseString(input string, data *Data) (string, error) {
	// Parse the input string using the data object
	tmpl, err := template.New("input").Parse(input)
	if err != nil {
		return "", err
	}
	// Execute the template with the provided data object
	var output bytes.Buffer
	err = tmpl.Execute(&output, data)
	if err != nil {
		return "", err
	}
	// Convert the output buffer to a string and print it
	return output.String(), nil
}

func HandleArgs(config utils.Config, query string) (string, string, string) {
	//if no no template name was supplied start the with choose
	if len(query) == 0 {
		return "choose", "", config.Base_color
	}
	//Match supplied parameter with known templates
	templates := make([]string, 0, len(config.Templates))
	for key := range config.Templates {
		templates = append(templates, key)
	}
	candidates := fuzzy.Find(query, templates)
	//exit it nothing is found
	if len(candidates) == 0 {
		fmt.Println("ERROR: No matching template found")
		os.Exit(1)
	}
	choice := candidates[0].Str
	return "info", choice, config.Templates[choice]

}
