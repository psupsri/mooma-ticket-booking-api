package utils

import "strconv"

func StringToUint(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
