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

	//get project path
	path := projects[candidates[0].Index].Path
	// Change to the specified directory
	if err := os.Chdir(path); err != nil {
		fmt.Println(err)
		return
	}
	// Open Neovim
	cmd := exec.Command(config.Editor)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
