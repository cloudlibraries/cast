package cast

import (
	"fmt"
	"reflect"
)

func ToMap[K comparable, V any](a any) map[K]V {
	v, _ := ToMapE[K, V](a)
	return v
}

// ToMapE casts an interface to a map[K]V type.
func ToMapE[K comparable, V any](a any) (map[K]V, error) {
	a = indirect(a)
	v := reflect.ValueOf(a)

	var m = map[K]V{}
	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			k, err := toGenericE[K](indirect(key.Interface()))
			if err != nil {
				return m, err
			}
			v, err := toGenericE[V](indirect(v.MapIndex(key).Interface()))
			if err != nil {
				return m, err
			}
			m[k] = v
		}
		return m, nil
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map", a, a)
	}
}

// ToSlice casts an interface to a []V type.
func ToSlice[V comparable](a any) []V {
	v, _ := ToSliceE[V](a)
	return v
}

// ToSliceE casts an interface to a []V type.
func ToSliceE[V comparable](a any) ([]V, error) {
	a = indirect(a)
	v := reflect.ValueOf(a)

	var s = []V{}
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			v, err := toGenericE[V](indirect(v.Index(i).Interface()))
			if err != nil {
				return s, err
			}
			s = append(s, v)
		}
		return s, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to slice", a, a)
	}
}

// ToFlatMapE
func ToFlatMapE(i any) (map[string]any, error) {
	switch v := i.(type) {
	case map[any]any:
		var m = map[string]any{}
		for k, val := range v {
			m[ToString(k)] = val
		}
		return flatten(m), nil
	case map[string]any:
		return flatten(v), nil
	case []any:
		var s = make([]any, 0)
		return flatten(append(s, v...)), nil
	case []map[string]any:
		var s = make([]any, 0)
		for _, u := range v {
			s = append(s, u)
		}
		return flatten(s), nil
	default:
		return map[string]any{}, fmt.Errorf("unable to cast %#v of type %T to flat map", i, i)
	}
}

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
