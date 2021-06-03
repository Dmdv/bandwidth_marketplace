package slices

import (
	"strings"
)

// ContainsStr returns true if slice contains string and false if not.
// Linear search is used.
func ContainsStr(sl []string, str string) bool {
	for _, slStr := range sl {
		if strings.Compare(slStr, str) == 0 {
			return true
		}
	}

	return false
}
