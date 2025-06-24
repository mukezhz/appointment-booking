package utils

import "strconv"

func AnyToUint(value interface{}) uint {
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case string:
		var result uint
		val, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0
		}
		result = uint(val)
		return result
	case uint:
		return v
	case int:
		return uint(v)
	case int8:
		return uint(v)
	case int16:
		return uint(v)
	case int32:
		return uint(v)
	case int64:
		return uint(v)
	default:
		return 0
	}
}
