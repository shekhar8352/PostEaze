package env

import (
	"fmt"
	"os"
	"strings"
)

var envObj map[string]string

func InitEnv() {
	env := os.Environ()
	envObj = make(map[string]string)
	for _, e := range env {
		s := strings.Split(e, "=")
		if len(s) >= 2 {
			envObj[strings.TrimSpace(s[0])] = strings.TrimSpace(strings.Join(s[1:], "="))
		}
	}
}

func ApplyEnvironmentToString(value string) string {
	for k, v := range envObj {
		value = strings.ReplaceAll(value, fmt.Sprintf("${%s}", k), v)
	}
	return value
}
