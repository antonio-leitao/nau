package configure

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	utils "github.com/antonio-leitao/nau/lib/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// these have to be lowercase for better matching
var customizableFields = []string{"AUTHOR", "EMAIL", "REMOTE", "BASE_COLOR", "EDITOR", "PROJECTS_PATH", "TEMPLATES_PATH", "ARCHIVES_PATH"}

func isCustomizableField(field string) bool {
	for _, f := range customizableFields {
		if f == field {
			return true
		}
	}
	return false
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func isEmailValid(email string) bool {
	regex := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}
func isValidUrl(url string) bool {
	regex := `[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	match, _ := regexp.MatchString(regex, url)
	return match
}
func isValidHexColor(color string) bool {
	regex := `^#(?:[0-9a-fA-F]{3}){1,2}$`
	match, _ := regexp.MatchString(regex, color)
	return match
}
func validateValue(field string, value string) string {
	switch field {
	case "EMAIL":
		if !isEmailValid(value) {
			return "• Email is not valid"
		}
	case "REMOTE", "WEBSITE":
		if !isValidUrl(value) {
			return "• Url is not valid"
		}
	case "PROJECTS_PATH", "TEMPLATES_PATH", "ARCHIVES_PATH":
		path, err := utils.ConvertPath(value)
		if err != nil {
			return "• " + err.Error()
		}
		if !dirExists(path) {
			return "• Directory does not exist"
		}
	case "BASE_COLOR":
		if !isValidHexColor(value) {
			return "• Not a valid hex color"
		}
	}
	return ""
}
func UpdateConfigField(field string, value string) error {
	//make it lowercase so we can match. maybe upper case?
	field = strings.ToUpper(field)
	//check if user can customize it
	if !isCustomizableField(field) {
		return fmt.Errorf("Invalid field: %s", field)
	}
	//check if the values are correct
	err_string := validateValue(field, value)
	if err_string != "" {
		return fmt.Errorf(err_string)
	}
	//migh tnot me able to get user
	configFile, err := utils.ConvertPath(".naurc")
	if err != nil {
		return err
	}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create an empty config file if it does not exist
		if _, err := os.Create(configFile); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(configFile, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Invalid config format")
		}

		key := strings.TrimSpace(parts[0])

		if key == field {
			lines = append(lines, fmt.Sprintf("%s=%s", field, value))
			found = true
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", field, value))
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}
func Init(config utils.Config) {
	model := initialModel(config.Base_color)
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
