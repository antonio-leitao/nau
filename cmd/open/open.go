package open

import (
	"os"
	"log"
	"os/exec"
	lib "github.com/antonio-leitao/nau/lib"
	"github.com/sahilm/fuzzy"
)

type Projects []lib.Project
func (p Projects) String(i int) string {
	return p[i].Name
}
func (p Projects) Len() int {
	return len(p)
}

func Execute(config lib.Config, query string) {
	projectList, _ := lib.GetProjects(config)
	projects := Projects(projectList)
	candidates := fuzzy.FindFrom(query, projects)

	//exit it nothing is found
	if len(candidates) == 0 {
		log.Println("NAU error: No project with a match found")
		os.Exit(1)
	}
	//get project path
	path := projects[candidates[0].Index].Path
	// Change to the specified directory
	if err := os.Chdir(path); err != nil {
        log.Printf("NAU error: %s",err)
		os.Exit(1)
	}
	// Open Neovim
	cmd := exec.Command(config.Editor)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
