package utils

import "strconv"

func AnyToUint(value any) (uint, bool) {
	switch v := value.(type) {
	case string:
		if parsed, err := strconv.ParseUint(v, 10, 32); err == nil {
			return uint(parsed), true
		}
	case uint:
		return v, true
	case int:
		return uint(v), true
	case int64:
		return uint(v), true
	case float64:
		return uint(v), true
	default:
		return 0, false
	}
	return 0, false
}
