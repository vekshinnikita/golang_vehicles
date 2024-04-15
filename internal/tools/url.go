package tools

import (
	"net/url"
	"strconv"

	"github.com/vekshinnikita/golang_vehicles/internal/validators"
)

func ReadIntUrlQuery(qs url.Values, v *validators.Validator, validates []validators.ValidateFunc[int], key string, defaultValue int) int {
	value := qs.Get(key)
	if value != "" {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			v.AddError(key, err.Error())
			return 0
		}

		for _, validate := range validates {
			validate(key, valueInt)
		}

		return valueInt
	}
	return defaultValue
}

func ReadStrUrlQuery(qs url.Values, key string, defaultValue string) string {
	value := qs.Get(key)
	if value != "" {
		return value
	}
	return defaultValue
}
