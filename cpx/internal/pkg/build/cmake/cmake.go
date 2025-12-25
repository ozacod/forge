package cmake

import (
	"os"
	"regexp"
)

func GetProjectNameFromCMakeLists() string {
	cmakeListsPath := "CMakeLists.txt"
	data, err := os.ReadFile(cmakeListsPath)
	if err != nil {
		return ""
	}

	// Look for: project(PROJECT_NAME ...)
	re := regexp.MustCompile(`project\s*\(\s*([^\s\)]+)`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
