package lib

import (
	"fmt"
	"os"
    "strings"
    "reflect"
    "regexp"
	"bufio"
)

var version = "v0.1.3"

type Config struct {
	Version        string
	Url            string
	Author         string
	Email          string
	Website        string
	Remote         string
	Base_color     string
	Projects_path  string
	Templates_path string
	Archives_path  string
	Editor         string
	Templates      map[string]string
	Projects       int
}

// these have to be lowercase for better matching
var CustomizableFields = []string{"AUTHOR", "EMAIL", "REMOTE", "BASE_COLOR", "EDITOR", "PROJECTS_PATH", "TEMPLATES_PATH", "ARCHIVES_PATH"}
func ReadConfig() (Config, error) {
    //read CONFIG file!
	defaultConfig := Config{
		Url:            "https://github.com/antonio-leitao/nau",
		Author:         "Antonio Leitao",
		Website:        "https://antonio-leitao.github.io",
		Email:          "aleitao@novaims.unl.pt",
		Remote:         "https://github.com/antonio-leitao",
		Base_color:     "#814584",
		Editor:         "nvim",
		Projects_path:  "~/Projects",
		Templates_path: "~/Templates",
		Archives_path:  "~/Archives",
	}
    //if file not exsits defaults config
	configFile, err := ExpandPath("~/.config/naurc")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return defaultConfig, nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := defaultConfig
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("Invalid config format")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "AUTHOR":
			config.Author = value
		case "EMAIL":
			config.Email = value
		case "WEBSITE":
			config.Website = value
		case "REMOTE":
			config.Remote = value
		case "EDITOR":
			config.Editor = value
		case "BASE_COLOR":
			config.Base_color = value
		case "PROJECTS_PATH":
			config.Projects_path = value
		case "TEMPLATES_PATH":
			config.Templates_path = value
		case "ARCHIVES_PATH":
			config.Archives_path = value
		default:
			return Config{}, fmt.Errorf("Unknown config field: %s", key)
		}
	}
	if err := scanner.Err(); err != nil {
		return Config{}, err
	}
	return config, nil
}
//############## EXPOSED FUNCTIONS #################s
// function to load the config stuff
func LoadConfig() (Config, error) {
	config, err := ReadConfig()
	if err != nil {
		return Config{}, err
	}
	//get templates
	templatesPath, err := ExpandPath(config.Templates_path)
	if err != nil {
		return Config{}, err
	}
	color_map, err := loadTemplatesColorMap(templatesPath)
	if err != nil {
		return Config{}, err
	}
	config.Templates = color_map
	//gget the number of projects
	project_count, err := countProjects(config)
	if err != nil {
		return Config{}, err
	}
	config.Projects = project_count
    //add version
    config.Version = version
	return config, nil
}

func isCustomizableField(field string) bool {
	for _, f := range CustomizableFields {
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
func ValidateValue(field string, value string) string {
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
		path, err := ExpandPath(value)
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
	err_string := ValidateValue(field, value)
	if err_string != "" {
		return fmt.Errorf(err_string)
	}
	//migh tnot me able to get user
	configFile, err := ExpandPath("~/.config/naurc")
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

func OutputField(config interface{}, field string) {
	configValue := reflect.ValueOf(config)
	if configValue.Kind() == reflect.Ptr {
		configValue = configValue.Elem()
	}
	if configValue.Kind() != reflect.Struct {
		fmt.Println("!")
		return
	}
	field = strings.ToUpper(field)
	fieldValue := configValue.FieldByNameFunc(func(fieldName string) bool {
		return strings.ToUpper(fieldName) == field
	})
	if !fieldValue.IsValid() {
		fmt.Println("!")
		return
	}
	fmt.Println(fieldValue.Interface())
}

