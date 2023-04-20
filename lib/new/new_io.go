package new

import (
	"fmt"
	"os"

	utils "github.com/antonio-leitao/nau/lib/utils"
	"github.com/sahilm/fuzzy"
)

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
