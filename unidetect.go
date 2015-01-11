package unidetect

import (
	"unicode"
)

// DetectScripts detects the Unicode scripts in the input string.
func DetectScripts(s string) (scripts []string) {
	set := make(map[string]bool)
	for _, r := range s {
		for script, rangeTable := range unicode.Scripts {
			if !set[script] && unicode.Is(rangeTable, r) {
				set[script] = true
				scripts = append(scripts, script)
			}
		}
	}
	return
}
