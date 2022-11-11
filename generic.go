package cast

import (
	"fmt"
	"math/big"
)

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
