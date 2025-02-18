package utils

import (
	"fmt"
	"strings"
)

func Template(template string, params map[string]string) string {
	result := template
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

func BulkTemplate(templates []string, params map[string]string) []string {
	results := make([]string, len(templates))
	for i, template := range templates {
		results[i] = Template(template, params)
	}
	return results
}
