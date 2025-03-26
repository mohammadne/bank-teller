package entities

import "strconv"

type Sheba string

func (s Sheba) Validate() bool {
	if len(s) != 24 {
		return false
	}
	if s[:2] != "IR" {
		return false
	}
	for _, v := range s[2:] {
		_, err := strconv.Atoi(string(v))
		if err != nil {
			return false
		}
	}
	return true
}
