package tools

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fatih/structs"
)

func GetUpdateQueryString(dict interface{}, excludeFields []string) ([]string, []interface{}, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsCounter := 1

	inputMap := structs.Map(dict)
	for _, field := range excludeFields {
		delete(inputMap, field)
	}

	for key, value := range inputMap {
		field, ok := reflect.ValueOf(dict).Type().FieldByName(key)
		dbTag := GetStructTag(field, "db")
		if !ok {
			message := fmt.Sprintf("field '%s' not found", key)
			return make([]string, 0), make([]interface{}, 0), errors.New(message)
		}

		if !reflect.ValueOf(value).IsZero() && dbTag != "" {
			setValues = append(setValues, fmt.Sprintf("%s=$%d", dbTag, argsCounter))
			args = append(args, value)
			argsCounter++
		}
	}

	return setValues, args, nil
}
