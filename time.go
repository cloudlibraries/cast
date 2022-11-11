package cast

import (
	"fmt"
	"math/big"
	"time"
)

// ToTime casts an interface to a time.Time type.
func ToTime(a any) time.Time {
	t, _ := ToTimeE(a)
	return t
}

// ToTimeE casts an interface to a time.Time type.
func ToTimeE(a any) (time.Time, error) {
	a = indirect(a)

	switch v := a.(type) {
	case int:
		return time.Unix(int64(v), 0), nil
	case int8:
		return time.Unix(int64(v), 0), nil
	case int16:
		return time.Unix(int64(v), 0), nil
	case int32:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(v, 0), nil
	case uint:
		return time.Unix(int64(v), 0), nil
	case uint8:
		return time.Unix(int64(v), 0), nil
	case uint16:
		return time.Unix(int64(v), 0), nil
	case uint32:
		return time.Unix(int64(v), 0), nil
	case uint64:
		return time.Unix(int64(v), 0), nil
	case float32:
		return time.Unix(int64(v), 0), nil
	case float64:
		return time.Unix(int64(v), 0), nil
	case *big.Int:
		return time.Unix(v.Int64(), 0), nil
	case *big.Float:
		n, _ := v.Int64()
		return time.Unix(n, 0), nil
	case *big.Rat:
		n, _ := v.Float64()
		return time.Unix(int64(n), 0), nil
	case complex64:
		return time.Unix(int64(real(v)), 0), nil
	case complex128:
		return time.Unix(int64(real(v)), 0), nil
	case bool:
		if v {
			return time.Unix(1, 0), nil
		}
		return time.Unix(0, 0), nil
	case string:
		return parseTime(v)
	case []byte:
		return parseTime(string(v))
	case fmt.Stringer:
		return parseTime(v.String())
	case time.Time:
		return v, nil
	case time.Duration:
		return time.Unix(0, int64(v)), nil
	case nil:
		return time.Time{}, nil
	default:
		return time.Time{}, fmt.Errorf("unable to cast %#v of type %T to Time", a, a)
	}
}

// ToDuration casts an interface to a time.Duration type.
func ToDuration(i any) time.Duration {
	d, _ := ToDurationE(i)
	return d
}

// ToDurationE casts an interface to a time.Duration type.
func ToDurationE(a any) (d time.Duration, err error) {
	a = indirect(a)

	switch v := a.(type) {
	case int:
		return time.Duration(v), nil
	case int8:
		return time.Duration(v), nil
	case int16:
		return time.Duration(v), nil
	case int32:
		return time.Duration(v), nil
	case int64:
		return time.Duration(v), nil
	case uint:
		return time.Duration(v), nil
	case uint8:
		return time.Duration(v), nil
	case uint16:
		return time.Duration(v), nil
	case uint32:
		return time.Duration(v), nil
	case uint64:
		return time.Duration(v), nil
	case float32:
		return time.Duration(v), nil
	case float64:
		return time.Duration(v), nil
	case *big.Int:
		return time.Duration(v.Int64()), nil
	case *big.Float:
		n, _ := v.Int64()
		return time.Duration(n), nil
	case *big.Rat:
		n, _ := v.Float64()
		return time.Duration(n), nil
	case complex64:
		return time.Duration(real(v)), nil
	case complex128:
		return time.Duration(real(v)), nil
	case bool:
		if v {
			return time.Duration(1), nil
		}
		return time.Duration(0), nil
	case string:
		return parseDuration(v)
	case []byte:
		return parseDuration(string(v))
	case fmt.Stringer:
		return parseDuration(v.String())
	case error:
		return parseDuration(v.Error())
	case time.Time:
		return v.Sub(time.Time{}), nil
	case time.Duration:
		return v, nil
	default:
		return time.Duration(0), fmt.Errorf("unable to cast %#v of type %T to Duration", a, a)
	}
}
