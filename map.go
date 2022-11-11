package cast

import (
	"fmt"
	"reflect"
)

// ToStringMapStringSliceE casts an interface to a map[string][]string type.
func ToStringMapStringSliceE(i any) (map[string][]string, error) {
	var m = map[string][]string{}

	switch v := i.(type) {
	case map[string][]string:
		return v, nil
	case map[string][]any:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[string]string:
		for k, val := range v {
			m[ToString(k)] = []string{val}
		}
	case map[string]any:
		for k, val := range v {
			switch vt := val.(type) {
			case []any:
				m[ToString(k)] = ToStringSlice(vt)
			case []string:
				m[ToString(k)] = vt
			default:
				m[ToString(k)] = []string{ToString(val)}
			}
		}
		return m, nil
	case map[any][]string:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[any]string:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[any][]any:
		for k, val := range v {
			m[ToString(k)] = ToStringSlice(val)
		}
		return m, nil
	case map[any]any:
		for k, val := range v {
			key, err := ToStringE(k)
			if err != nil {
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			}
			value, err := ToStringSliceE(val)
			if err != nil {
				return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
			}
			m[key] = value
		}
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string][]string", i, i)
	}
	return m, nil
}

// ToStringMapBoolE casts an interface to a map[string]bool type.
func ToStringMapBoolE(i any) (map[string]bool, error) {
	var m = map[string]bool{}

	switch v := i.(type) {
	case map[any]any:
		for k, val := range v {
			m[ToString(k)] = ToBool(val)
		}
		return m, nil
	case map[string]any:
		for k, val := range v {
			m[ToString(k)] = ToBool(val)
		}
		return m, nil
	case map[string]bool:
		return v, nil
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]bool", i, i)
	}
}

// ToStringMapE casts an interface to a map[string]any type.
func ToStringMapE(i any) (map[string]any, error) {
	var m = map[string]any{}

	switch v := i.(type) {
	case map[any]any:
		for k, val := range v {
			m[ToString(k)] = val
		}
		return m, nil
	case map[string]any:
		return v, nil
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]any", i, i)
	}
}

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

// ToStringMapIntE casts an interface to a map[string]int{} type.
func ToStringMapIntE(i any) (map[string]int, error) {
	var m = map[string]int{}
	if i == nil {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int", i, i)
	}

	switch v := i.(type) {
	case map[any]any:
		for k, val := range v {
			m[ToString(k)] = ToInt(val)
		}
		return m, nil
	case map[string]any:
		for k, val := range v {
			m[k] = ToInt(val)
		}
		return m, nil
	case map[string]int:
		return v, nil
	}

	if reflect.TypeOf(i).Kind() != reflect.Map {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int", i, i)
	}

	mVal := reflect.ValueOf(m)
	v := reflect.ValueOf(i)
	for _, keyVal := range v.MapKeys() {
		val, err := ToIntE(v.MapIndex(keyVal).Interface())
		if err != nil {
			return nil, fmt.Errorf("unable to cast %#v of type %T to map[string]int", i, i)
		}
		mVal.SetMapIndex(keyVal, reflect.ValueOf(val))
	}
	return m, nil
}

// ToStringMapInt64E casts an interface to a map[string]int64{} type.
func ToStringMapInt64E(i any) (map[string]int64, error) {
	var m = map[string]int64{}
	if i == nil {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", i, i)
	}

	switch v := i.(type) {
	case map[any]any:
		for k, val := range v {
			m[ToString(k)] = ToInt64(val)
		}
		return m, nil
	case map[string]any:
		for k, val := range v {
			m[k] = ToInt64(val)
		}
		return m, nil
	case map[string]int64:
		return v, nil
	}

	if reflect.TypeOf(i).Kind() != reflect.Map {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", i, i)
	}
	mVal := reflect.ValueOf(m)
	v := reflect.ValueOf(i)
	for _, keyVal := range v.MapKeys() {
		val, err := ToInt64E(v.MapIndex(keyVal).Interface())
		if err != nil {
			return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", i, i)
		}
		mVal.SetMapIndex(keyVal, reflect.ValueOf(val))
	}
	return m, nil
}
