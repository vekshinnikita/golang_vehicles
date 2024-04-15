package validators

import (
	"fmt"
	"math"
	"regexp"
)

type ValidatorError struct {
	Key     string
	Message string
}

type Validator struct {
	errors []ValidatorError
}

type ValidateFunc[T any] func(key string, value T)

type Number interface {
	int | float64
}

func NewValidator() *Validator {
	return &Validator{
		errors: make([]ValidatorError, 0),
	}
}

func (v *Validator) AddError(key, message string) {
	v.errors = append(v.errors, ValidatorError{
		Key:     key,
		Message: message,
	})
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *Validator) GetErrorsDict() map[string]interface{} {
	errorsDict := map[string]interface{}{}

	for _, err := range v.errors {
		errorsDict[err.Key] = err.Message
	}

	return errorsDict
}

func ValidateAll[T any](validates []ValidateFunc[T], key string, value T) {
	for _, validate := range validates {
		validate(key, value)
	}
}

func ValidateRegNum(validator *Validator) ValidateFunc[string] {

	return func(key, regNum string) {
		characterSet := "АВЕКМНОРСТУХабекмнорстухABEKMHOPCTYXabekmhopctyx"
		pattern := fmt.Sprintf(`^[%s]\d{3}[%s]{2}\d{1,3}$`, characterSet, characterSet)

		matched, _ := regexp.MatchString(pattern, regNum)
		if !matched {
			validator.AddError(key, "this field must has 'X000XX00' format")
		}
	}

}

func ValidateNumber[T Number](validator *Validator, max, min T) ValidateFunc[T] {
	maxValue := float64(max)
	minValue := float64(min)
	if max == 0 {
		maxValue = math.Inf(1)
	}

	return func(key string, v T) {
		value := float64(v)

		if value < minValue {
			validator.AddError(
				key,
				fmt.Sprintf("value must be grater than or equal to %f", minValue),
			)
		}

		if value > maxValue {
			validator.AddError(
				key,
				fmt.Sprintf("value must be less than or equal to %f", maxValue),
			)
		}
	}
}
