package configure

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	utils "github.com/antonio-leitao/nau/lib/utils"
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
func UpdateConfigField(field string, value string) error {
	//make it lowercase so we can match. maybe upper case?
	field = strings.ToUpper(field)
	//check if user can customize it
	if !isCustomizableField(field) {
		return fmt.Errorf("Invalid field: %s", field)
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
