package open

import (
	lib "github.com/antonio-leitao/nau/lib"
	"github.com/sahilm/fuzzy"
	"log"
	"os"
	"os/exec"
)

type Projects []lib.Project

func (p Projects) String(i int) string {
	return p[i].Name
}
func (p Projects) Len() int {
	return len(p)
}
func Open(path string, editor string) {
	// Change to the specified directory
	if err := os.Chdir(path); err != nil {
		log.Printf("NAU error: %s", err)
		os.Exit(1)
	}
	// Open Neovim
	cmd := exec.Command(editor)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

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
	Open(path, config.Editor)
}
