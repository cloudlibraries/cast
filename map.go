package cast

import (
	"fmt"
	"reflect"
)

func flattenMap(m map[string]any) map[string]any {
	return flattenMapWithDelimiter(m, ".")
}

func flattenMapWithDelimiter(m map[string]any, delimiter string) map[string]any {
	flatMap := make(map[string]any)
	for k, v := range m {
		flattenMapWithDelimiterRecursive(flatMap, k, v, delimiter)
	}
	return flatMap
}

func flattenMapWithDelimiterRecursive(flatMap map[string]any, key string, value any, delimiter string) {
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Map:
		for k, v := range value.(map[string]any) {
			flattenMapWithDelimiterRecursive(flatMap, key+delimiter+k, v, delimiter)
		}
	case reflect.Slice:
		for i, v := range value.([]any) {
			flattenMapWithDelimiterRecursive(flatMap, key+delimiter+fmt.Sprintf("%d", i), v, delimiter)
		}
	default:
		flatMap[key] = value
	}
}
