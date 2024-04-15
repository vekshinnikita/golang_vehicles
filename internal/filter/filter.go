package filter

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vekshinnikita/golang_vehicles/internal/validators"
)

const (
	OperatorGT  = "gt"
	OperatorLT  = "lt"
	OperatorEq  = "eq"
	OperatorBtw = "btw"
	OperatorLK  = "lk"
)

const (
	DataTypeInt = "int"
	DataTypeStr = "str"
)

type FilterField struct {
	Name     string
	Value    string
	Operator string
	Type     string
}

type Filters struct {
	Fields []FilterField
}

func NewFilter() *Filters {
	return &Filters{
		Fields: make([]FilterField, 0),
	}
}

func (f *Filters) GetSQLFilterQuery() string {
	var setFilter []string
	fmt.Println(f.Fields)
	for _, field := range f.Fields {
		value := field.Value
		if field.Type == DataTypeStr {
			value = fmt.Sprintf("'%s'", field.Value)
		}

		switch field.Operator {
		case OperatorEq:
			setFilter = append(setFilter, fmt.Sprintf("%s=%s", field.Name, value))
		case OperatorBtw:
			value := strings.Split(value, ":")
			setFilter = append(setFilter, fmt.Sprintf("%s BETWEEN %s AND %s", field.Name, value[0], value[1]))
		case OperatorGT:
			setFilter = append(setFilter, fmt.Sprintf("%s>%s", field.Name, value))
		case OperatorLT:
			setFilter = append(setFilter, fmt.Sprintf("%s<%s", field.Name, value))
		case OperatorLK:
			setFilter = append(setFilter, fmt.Sprintf("LOWER(%s) LIKE '%%%s%%'", field.Name, strings.ToLower(field.Value)))
		}
	}

	if len(setFilter) > 0 {
		return "WHERE " + strings.Join(setFilter, " AND ")
	}
	return ""
}

func (f *Filters) add(name, operator, dtype string, value string) {
	f.Fields = append(f.Fields, FilterField{
		Name:     name,
		Value:    value,
		Operator: operator,
		Type:     dtype,
	})
}

func (f *Filters) HandleAddString(qs url.Values, validates []validators.ValidateFunc[string], key, DBKey string) {

	operators := [...]string{OperatorEq, OperatorLK}
	for _, operator := range operators {
		var value string
		switch operator {
		case OperatorEq:
			value = qs.Get(key)
			if value != "" {
				validators.ValidateAll(validates, key, value)
			}
		default:
			value = qs.Get(fmt.Sprintf("%s_%s", key, operator))
		}

		if value != "" {
			f.add(
				DBKey,
				operator,
				DataTypeStr,
				value,
			)
		}
	}
}

func getBtwIntValue(value string) (string, error) {
	if !strings.Contains(value, ":") {
		return "", errors.New("this field must have the format 'int:int'")
	}
	rangeValue := strings.Split(value, ":")
	start, err := strconv.Atoi(rangeValue[0])
	if err != nil {
		return "", errors.New("first value must be an integer")
	}

	end, err := strconv.Atoi(rangeValue[1])
	if err != nil {
		return "", errors.New("second value must be an integer")
	}

	if start >= end {
		return "", errors.New("first value must be grater than second value")
	}

	return value, nil
}

func getParseIntValue[V any](value string, parseFunc func(string) (V, error)) (V, error) {
	if value != "" {
		return parseFunc(value)
	}
	var r V
	return r, errors.New("")
}

func (f *Filters) HandleAddInt(qs url.Values, validator *validators.Validator, validates []validators.ValidateFunc[int], key, DBKey string) error {

	operators := [...]string{OperatorGT, OperatorLT, OperatorEq, OperatorBtw}

	for _, operator := range operators {
		var queryKey string
		var parsedValue string
		var validateValue int
		var err error

		switch operator {
		case OperatorEq:
			queryKey = key
			parsedValue = qs.Get(queryKey)

			validateValue, err = getParseIntValue(parsedValue, strconv.Atoi)

		case OperatorBtw:
			queryKey = fmt.Sprintf("%s_%s", key, operator)
			parsedValue = qs.Get(queryKey)

			_, err = getParseIntValue(parsedValue, getBtwIntValue)
		default:
			queryKey = fmt.Sprintf("%s_%s", key, operator)
			parsedValue = qs.Get(queryKey)

			validateValue, err = getParseIntValue(parsedValue, strconv.Atoi)
		}

		fmt.Println(err)

		if err != nil {
			if parsedValue != "" {
				validator.AddError(
					queryKey,
					err.Error(),
				)
			}

			continue
		}

		if operator != OperatorBtw {
			validators.ValidateAll(validates, key, validateValue)
		}

		f.add(
			DBKey,
			operator,
			DataTypeInt,
			parsedValue,
		)
	}

	return nil
}
