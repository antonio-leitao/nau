package open

import (
	"fmt"
	"os"
	"os/exec"

	utils "github.com/antonio-leitao/nau/lib/utils"
	"github.com/sahilm/fuzzy"
)

type Projects []utils.Project

func (p Projects) String(i int) string {
	return p[i].Name
}

func (p Projects) Len() int {
	return len(p)
}

func Open(config utils.Config, query string) {
	projectList, _ := utils.GetProjects(config)
	projects := Projects(projectList)
	candidates := fuzzy.FindFrom(query, projects)

	//exit it nothing is found
	if len(candidates) == 0 {
		fmt.Println("ERROR: No project found")
		os.Exit(1)
	}

	//open vscode if something is found
	path := projects[candidates[0].Index].Path
	cmd := exec.Command(config.Editor, path)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
