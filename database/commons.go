package database

import "strings"

//Contains (string in []string)
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveTabsFromLines(lines []string) []string {
	//Remove the first group of tabs (before the real characters)
	for i := 0; i < len(lines); i++ {
	redo:
		if strings.HasPrefix(lines[i], "\t") {
			lines[i] = strings.TrimPrefix(lines[i], "\t")
			goto redo
		}
	}

	return lines
}
