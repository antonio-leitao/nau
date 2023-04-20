package new

import (
	"fmt"

	utils "github.com/antonio-leitao/nau/lib/utils"
)

// stuff that gets exported
func GetCodesAndNames(path string, projectTypes []string) ([]string, []string) {
	projects, err := utils.GetProjects(path, projectTypes)
	if err != nil {
		fmt.Println(err)
	}
	var codes []string
	var repoNames []string
	for _, project := range projects {
		codes = append(codes, project.Code)
		repoNames = append(repoNames, project.Repo_name)
	}
	return codes, repoNames
}


