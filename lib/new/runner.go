package new

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func contains(slice []string, str string) bool {
    for _, item := range slice {
        if item == str {
            return true
        }
    }
    return false
}
func validEntry(entry os.DirEntry)bool{
	if !entry.IsDir() {
		return false
	}
	if strings.HasPrefix(entry.Name(), ".") {
		return false
	}
	return true
}

func getThemedProjects(path string, lang string)[]string{
	var themedProjects []string
	subentries, _ := os.ReadDir(path+"/"+lang)
	for _,subentry := range subentries{
		if !validEntry(subentry){continue}
		project := subentry.Name()
		themedProjects = append(themedProjects, project)
	}
	return themedProjects
}

func getProjectNames(path string, projectTypes []string) ([]string, error) {
	var projectNames []string
	entries, err := os.ReadDir(path)
    if err != nil {
        return nil, err
    }
    for _, entry := range entries {
		if !validEntry(entry){continue}
		if contains(projectTypes,entry.Name()){
			themedProjects:=getThemedProjects(path,entry.Name())
			projectNames = append(projectNames, themedProjects...)
		} else {
			project := entry.Name()
			projectNames = append(projectNames, project)
		}
    }
    return projectNames, nil
}

func handleProjectNames(projectNames []string)([]string, []string){
	var codes []string
	var folderNames []string
	for _, projectName := range projectNames{
		codes = append(codes,projectName[:3])
		folderNames = append(folderNames,ToHyphenName(projectName[:3]))  
	}
	return codes, folderNames
}


//stuff that gets exported
func GetProjects(path string, projectTypes []string)([]string, []string){
	projectNames, err :=getProjectNames(path,projectTypes)
	if err !=nil{
		fmt.Println(err)
	}
	codes, folderNames := handleProjectNames(projectNames)
	return codes, folderNames
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