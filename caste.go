// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ToStringE casts an interface to a string type.
func ToStringE(i any) (string, error) {
	i = indirectToStringerOrError(i)

	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case json.Number:
		return s.String(), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

// ToBoolE casts an interface to a bool type.
func ToBoolE(i any) (bool, error) {
	i = indirect(i)

	switch b := i.(type) {
	case bool:
		return b, nil
	case nil:
		return false, nil
	case int:
		if i.(int) != 0 {
			return true, nil
		}
		return false, nil
	case string:
		return strconv.ParseBool(i.(string))
	case json.Number:
		v, err := ToInt64E(b)
		if err == nil {
			return v != 0, nil
		}
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	default:
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", i, i)
	}
}

// ToStringMapStringE casts an interface to a map[string]string type.
func ToStringMapStringE(i any) (map[string]string, error) {
	var m = map[string]string{}

	switch v := i.(type) {
	case map[string]string:
		return v, nil
	case map[string]any:
		for k, val := range v {
			m[ToString(k)] = ToString(val)
		}
		return m, nil
	case map[any]string:
		for k, val := range v {
			m[ToString(k)] = ToString(val)
		}
		return m, nil
	case map[any]any:
		for k, val := range v {
			m[ToString(k)] = ToString(val)
		}
		return m, nil
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]string", i, i)
	}
}


// ToSliceE casts an interface to a []any type.
func ToSliceE(i any) ([]any, error) {
	var s []any

	switch v := i.(type) {
	case []any:
		return append(s, v...), nil
	case []map[string]any:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []any", i, i)
	}
}

// ToBoolSliceE casts an interface to a []bool type.
func ToBoolSliceE(i any) ([]bool, error) {
	if i == nil {
		return nil, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
	}

	switch v := i.(type) {
	case []bool:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]bool, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ToBoolE(s.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []bool", i, i)
	}
}

// ToStringSliceE casts an interface to a []string type.
func ToStringSliceE(i any) ([]string, error) {
	var s []string

	switch v := i.(type) {
	case []any:
		for _, u := range v {
			s = append(s, ToString(u))
		}
		return s, nil
	case []string:
		return v, nil
	case []int8:
		for _, u := range v {
			s = append(s, ToString(u))
		}
		return s, nil
	case []int:
		for _, u := range v {
			s = append(s, ToString(u))
		}
		return s, nil
	case []int32:
		for _, u := range v {
			s = append(s, ToString(u))
		}
		return s, nil
	case []int64:
		for _, u := range v {
			s = append(s, ToString(u))
		}
		return s, nil
	case []float32:
		for _, u := range v {
			s = append(s, ToString(u))
		}
		return s, nil
	case []float64:
		for _, u := range v {
			s = append(s, ToString(u))
		}
		return s, nil
	case string:
		return strings.Fields(v), nil
	case []error:
		for _, err := range i.([]error) {
			s = append(s, err.Error())
		}
		return s, nil
	case any:
		str, err := ToStringE(v)
		if err != nil {
			return s, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
		}
		return []string{str}, nil
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []string", i, i)
	}
}

// ToIntSliceE casts an interface to a []int type.
func ToIntSliceE(i any) ([]int, error) {
	if i == nil {
		return nil, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	}

	switch v := i.(type) {
	case []int:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ToIntE(s.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []int", i, i)
	}
}

// ToDurationSliceE casts an interface to a []time.Duration type.
func ToDurationSliceE(i any) ([]time.Duration, error) {
	if i == nil {
		return nil, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
	}

	switch v := i.(type) {
	case []time.Duration:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]time.Duration, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ToDurationE(s.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []time.Duration", i, i)
	}
}

func ToErrorE(i any) (error, error) {
	switch v := i.(type) {
	case error:
		return v, nil
	case string:
		return fmt.Errorf("%s", v), nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to error", i, i)
	}
}
