package configure

import (
	"fmt"
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
	structs "github.com/antonio-leitao/nau/lib/structs"
)

func UpdateConfigField(config *structs.Config, field string, value interface{}) error {
	v := reflect.ValueOf(config).Elem()
	fieldValue := v.FieldByName(field)

	if !fieldValue.IsValid() {
		return fmt.Errorf("invalid field name: %s", field)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("cannot set field value: %s", field)
	}

	fieldType := fieldValue.Type()
	val := reflect.ValueOf(value)

	if !val.Type().ConvertibleTo(fieldType) {
		return fmt.Errorf("value is not convertible to field type: %s", fieldType)
	}

	fieldValue.Set(val.Convert(fieldType))

	// Encode the updated struct to TOML and save it to the file
	filePath := "nau.config.toml"
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = toml.NewEncoder(f).Encode(config)
	if err != nil {
		return err
	}

	return nil
}
