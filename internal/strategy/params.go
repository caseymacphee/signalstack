package strategy

import (
	"fmt"
	"strconv"
)

func parseIntParam(m map[string]string, key string, def int) (int, error) {
	if m == nil {
		return def, nil
	}
	raw, ok := m[key]
	if !ok {
		return def, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return def, fmt.Errorf("invalid int value for %s: %v", key, err)
	}
	return value, nil
}
