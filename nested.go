package cast

import (
	"fmt"
	"math/big"
	"reflect"
)

// ToMap casts an interface to a map[K]V type.
// V must use types after indirect.
func ToMap[K comparable, V any](a any) map[K]V {
	a = indirect(a)
	v := reflect.ValueOf(a)

	var m = map[K]V{}
	for _, key := range v.MapKeys() {
		k := toGeneric[K](indirect(key.Interface()))
		v := toGeneric[V](indirect(v.MapIndex(key).Interface()))
		m[k] = v
	}
	return m
}

// ToMapE casts an interface to a map[K]V type.
// V must use types after indirect.
func ToMapE[K comparable, V any](a any) (map[K]V, error) {
	a = indirect(a)
	v := reflect.ValueOf(a)

	var m = map[K]V{}
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
}

// ToSlice casts an interface to a []V type.
func ToSlice[V comparable](a any) []V {
	a = indirect(a)
	v := reflect.ValueOf(a)

	var s = []V{}
	for i := 0; i < v.Len(); i++ {
		v := toGeneric[V](indirect(v.Index(i).Interface()))
		s = append(s, v)
	}
	return s
}

// ToSliceE casts an interface to a []V type.
func ToSliceE[V comparable](a any) ([]V, error) {
	a = indirect(a)
	v := reflect.ValueOf(a)

	var s = []V{}
	for i := 0; i < v.Len(); i++ {
		v, err := toGenericE[V](indirect(v.Index(i).Interface()))
		if err != nil {
			return s, err
		}
		s = append(s, v)
	}
	return s, nil
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

// toGeneric casts an interface to a T type.
// T must use types after indirect.
func toGeneric[T any](a any) T {
	t, _ := toGenericE[T](a)
	return t
}

// toGenericE casts an interface to a T type.
// T must use types after indirect.
func toGenericE[T any](a any) (t T, err error) {
	switch any(t).(type) {
	case int:
		v, err := ToIntE(a)
		return any(v).(T), err
	case int8:
		v, err := ToInt8E(a)
		return any(v).(T), err
	case int16:
		v, err := ToInt16E(a)
		return any(v).(T), err
	case int32:
		v, err := ToInt32E(a)
		return any(v).(T), err
	case int64:
		v, err := ToInt64E(a)
		return any(v).(T), err
	case uint:
		v, err := ToUintE(a)
		return any(v).(T), err
	case uint8:
		v, err := ToUint8E(a)
		return any(v).(T), err
	case uint16:
		v, err := ToUint16E(a)
		return any(v).(T), err
	case uint32:
		v, err := ToUint32E(a)
		return any(v).(T), err
	case uint64:
		v, err := ToUint64E(a)
		return any(v).(T), err
	case float32:
		v, err := ToFloat32E(a)
		return any(v).(T), err
	case float64:
		v, err := ToFloat64E(a)
		return any(v).(T), err
	case *big.Int:
		v, err := ToBigIntE(a)
		return any(v).(T), err
	case *big.Float:
		v, err := ToBigFloatE(a)
		return any(v).(T), err
	case *big.Rat:
		v, err := ToBigRatE(a)
		return any(v).(T), err
	case complex64:
		v, err := ToComplex64E(a)
		return any(v).(T), err
	case complex128:
		v, err := ToComplex128E(a)
		return any(v).(T), err
	case bool:
		v, err := ToBoolE(a)
		return any(v).(T), err
	case string:
		v, err := ToStringE(a)
		return any(v).(T), err
	case []byte:
		v, err := ToBytesE(a)
		return any(v).(T), err
	case fmt.Stringer:
		v, err := ToStringE(a)
		return any(v).(T), err
	case error:
		v, err := ToErrorE(a)
		return any(v).(T), err
	default:
		return t, fmt.Errorf("unable to cast %#v of type %T to %T", a, a, t)
	}
}
