package cast

import (
	"fmt"
	"reflect"
)

const (
	flattenDelimiter = "."
)

func flatten(v any) map[string]any {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Map:
		shadowMap := make(map[string]any)
		for k, v := range v.(map[string]any) {
			flattenWithDelimiterRecursive(shadowMap, k, v, flattenDelimiter)
		}
		return shadowMap
	case reflect.Slice:
		shadowMap := make(map[string]any)
		for i, v := range v.([]any) {
			flattenWithDelimiterRecursive(shadowMap, fmt.Sprintf("%d", i), v, flattenDelimiter)
		}
		return shadowMap
	default:
		return nil
	}
}

func flattenWithDelimiterRecursive(shadowMap map[string]any, key string, value any, delimiter string) {
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Map:
		for k, v := range value.(map[string]any) {
			flattenWithDelimiterRecursive(shadowMap, key+delimiter+k, v, delimiter)
		}
	case reflect.Slice:
		for i, v := range value.([]any) {
			flattenWithDelimiterRecursive(shadowMap, key+delimiter+fmt.Sprintf("%d", i), v, delimiter)
		}
	default:
		shadowMap[key] = value
	}
}
