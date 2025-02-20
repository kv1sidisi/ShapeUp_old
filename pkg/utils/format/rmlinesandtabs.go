package format

import "strings"

// RemoveLinesAndTabs removes \n and \t from string.
func RemoveLinesAndTabs(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\t", "")
	return input
}
