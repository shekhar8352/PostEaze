package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var envObj map[string]string

func InitEnv() {
	_ = godotenv.Load()
	env := os.Environ()
	envObj = make(map[string]string)
	for _, e := range env {
		s := strings.Split(e, "=")
		if len(s) >= 2 {
			key := strings.TrimSpace(s[0])
			value := strings.Join(s[1:], "=") // Don't trim whitespace from values
			envObj[key] = value
		}
	}
}

func ApplyEnvironmentToString(value string) string {
	// Find all ${VAR} patterns in the ORIGINAL string only
	// and replace them one by one to avoid nested replacement
	result := value
	originalValue := value
	
	// Find all variables in the original string
	var varsToReplace []struct {
		placeholder string
		varName     string
		start       int
		end         int
	}
	
	searchPos := 0
	for {
		start := strings.Index(originalValue[searchPos:], "${")
		if start == -1 {
			break
		}
		start += searchPos
		
		end := strings.Index(originalValue[start:], "}")
		if end == -1 {
			break
		}
		end += start
		
		varName := originalValue[start+2 : end]
		placeholder := originalValue[start : end+1]
		
		varsToReplace = append(varsToReplace, struct {
			placeholder string
			varName     string
			start       int
			end         int
		}{placeholder, varName, start, end})
		
		searchPos = end + 1
	}
	
	// Replace variables from right to left to maintain positions
	for i := len(varsToReplace) - 1; i >= 0; i-- {
		v := varsToReplace[i]
		
		// Get the value from environment, use empty string if not found
		envValue := ""
		if val, exists := envObj[v.varName]; exists {
			envValue = val
		}
		
		// Replace in the result string
		result = result[:v.start] + envValue + result[v.end+1:]
	}
	
	return result
}
