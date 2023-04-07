package open

import (
	"fmt"
	"os"
	"strings"

	structs "github.com/antonio-leitao/nau/lib/structs"
	"github.com/sahilm/fuzzy"
)

type Project struct {
	name          string
	path          string
}

type Projects []Project

func (p Projects) String(i int) string {
	return p[i].name
}

func (p Projects) Len() int {
	return len(p)
}

func getThemedProjects(path string, lang string)[]Project{
	var themedProjects []Project
	subentries, _ := os.ReadDir(path+"/"+lang)
	for _,subentry := range subentries{
		if !validEntry(subentry){continue}
		project := Project{
			name:subentry.Name(),
			path: path+"/"+lang+"/"+subentry.Name(),
		}
		themedProjects = append(themedProjects, project)
	}
	return themedProjects

}

func getProjectNames(path string, projectTypes []string) ([]Project, error) {
	var projectNames []Project
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
			project := Project{
				name:entry.Name(),
				path: path+"/"+entry.Name(),
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
func validEntry(entry os.DirEntry)bool{
	if !entry.IsDir() {
		return false
	}
	if strings.HasPrefix(entry.Name(), ".") {
		return false
	}
	return true
}

func Open(config structs.Config, query string) {
	projectList,_ := getProjectNames(config.Projects_path,config.Projects_themes)
	projects := Projects(projectList)
	candidates := fuzzy.FindFrom(query, projects)
	if len(candidates)==0{
		fmt.Println("ERROR: No project found")
		os.Exit(1)
	}else{
		fmt.Println(projects[candidates[0].Index])
	}
	
}