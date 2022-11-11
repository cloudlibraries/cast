package cast

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/shopspring/decimal"
)

// ToString casts an interface to a string type.
func ToString(a any) string {
	v, _ := ToStringE(a)
	return v
}

// ToStringE casts an interface to a string type.
func ToStringE(a any) (string, error) {
	a = indirectToStringerOrError(a)

	switch v := a.(type) {
	case int:
		return strconv.Itoa(v), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.Itoa(int(v)), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(v), 10), nil
	case float32:
		// Use decimal to fix precision issue, FormatFloat is unstable.
		// Optional:
		//		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
		return decimal.NewFromFloat32(v).String(), nil
	case float64:
		// Use decimal to fix precision issue, FormatFloat is unstable.
		// Optional:
		// 		return strconv.FormatFloat(s, 'f', -1, 64), nil
		return decimal.NewFromFloat(v).String(), nil
	case *big.Int:
		return v.String(), nil
	case *big.Float:
		return v.String(), nil
	case *big.Rat:
		return v.String(), nil
	case bool:
		return strconv.FormatBool(v), nil
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case fmt.Stringer:
		return v.String(), nil
	case error:
		return v.Error(), nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", a, a)
	}
}

// ToBytes casts an interface to a []byte type.
func ToBytes(a any) []byte {
	v, _ := ToBytesE(a)
	return v
}

// ToBytesE casts an interface to a []byte type.
func ToBytesE(a any) ([]byte, error) {
	a = indirectToStringerOrError(a)

	switch v := a.(type) {
	case int:
		return []byte(strconv.Itoa(v)), nil
	case int8:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int16:
		return []byte(strconv.FormatInt(int64(v), 10)), nil
	case int32:
		return []byte(strconv.Itoa(int(v))), nil
	case int64:
		return []byte(strconv.FormatInt(v, 10)), nil
	case uint:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint8:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint16:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint32:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case uint64:
		return []byte(strconv.FormatUint(uint64(v), 10)), nil
	case float32:
		// Use decimal to fix precision issue, FormatFloat is unstable.
		// Optional:
		//		return []byte(strconv.FormatFloat(float64(s), 'f', -1, 32)), nil
		return []byte(decimal.NewFromFloat32(v).String()), nil
	case float64:
		// Use decimal to fix precision issue, FormatFloat is unstable.
		// Optional:
		// 		return []byte(strconv.FormatFloat(s, 'f', -1, 64)), nil
		return []byte(decimal.NewFromFloat(v).String()), nil
	case *big.Int:
		return []byte(v.String()), nil
	case *big.Float:
		return []byte(v.String()), nil
	case *big.Rat:
		return []byte(v.String()), nil
	case bool:
		return []byte(strconv.FormatBool(v)), nil
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case fmt.Stringer:
		return []byte(v.String()), nil
	case error:
		return []byte(v.Error()), nil
	case nil:
		return []byte{}, nil
	default:
		return []byte{}, fmt.Errorf("unable to cast %#v of type %T to []byte", a, a)
	}
}

// ToStringer casts an interface to a fmt.Stringer type.
func ToStringer(a any) fmt.Stringer {
	v, _ := ToStringerE(a)
	return v
}

// ToStringerE casts an interface to a fmt.Stringer type.
func ToStringerE(a any) (fmt.Stringer, error) {
	a = indirectToStringerOrError(a)

	switch v := a.(type) {
	case fmt.Stringer:
		return v, nil
	default:
		if _, err := ToStringE(a); err == nil {
			return stringer{a}, nil
		}
		return nil, fmt.Errorf("unable to cast %#v of type %T to fmt.Stringer", a, a)
	}
}

// ToError casts an interface to an error type.
func ToError(a any) error {
	v, _ := ToErrorE(a)
	return v
}

// ToErrorE casts an interface to an error type.
func ToErrorE(a any) (error, error) {
	a = indirectToStringerOrError(a)

	switch v := a.(type) {
	case error:
		return v, nil
	default:
		if _, err := ToStringE(a); err == nil {
			return fmt.Errorf("%v", a), nil
		}
		return nil, fmt.Errorf("unable to cast %#v of type %T to error", a, a)
	}
}

type stringer struct{ any }

func (s stringer) String() string {
	return ToString(any(s))
}

func autoParseInt(s string) (int64, error) {
	var foundZero bool
loop:
	for i := len(s); i > 0; i-- {
		switch s[i-1] {
		case '.':
			if foundZero {
				s = s[:i-1]
				break loop
			}
		case '0':
			foundZero = true
		default:
			break loop
		}
	}
	return strconv.ParseInt(s, 0, 0)
}
