package utils

import (
	"os"
	"regexp"
)

func ReplacePlaceHoldersWithEnv(input string) string {
	// Regex to find ${VAR_NAME} patterns
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_]+)\}`)

	// Replace all matches using a function
	return re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract VAR_NAME
		submatch := re.FindStringSubmatch(match)
		if len(submatch) == 2 {
			envVar := os.Getenv(submatch[1])
			return envVar
		}
		return match // fallback to original if not matched properly
	})
}
