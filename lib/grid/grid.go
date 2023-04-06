package grid

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	structs "github.com/antonio-leitao/nau/lib/structs"
)

func getThemedProjects(path string, lang string)[]structs.Project{
	var themedProjects []structs.Project
	subentries, _ := os.ReadDir(path+"/"+lang)
	for _,subentry := range subentries{
		if !validEntry(subentry){continue}
		project := structs.Project{
			Name:subentry.Name(),
			Language: lang,
			Desc: "Lorem Ipsum",
		}
		themedProjects = append(themedProjects, project)
	}
	return themedProjects

}

func getProjectNames(path string, projectTypes []string) ([]structs.Project, error) {
	var projectNames []structs.Project
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
			project:=structs.Project{
				Name:entry.Name(),
				Language: "Mixed",
				Desc: "Lorem Ipsum",
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

//BUBBLE TEA MODEL

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Grid(config structs.Config) {
	//get project names
	projectNames,_ := getProjectNames(config.Projects_path,config.Projects_themes)
	items := make([]list.Item, len(projectNames))
	for i, project := range projectNames {
		items[i] = list.Item(project)
	}

	//start list with those names
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Projects"


	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}