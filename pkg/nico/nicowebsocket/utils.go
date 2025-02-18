package nicowebsocket

import (
	"fmt"
	"strings"
)

func CookieBuilder(cookies []Cookie) string {
	var result []string
	for _, c := range cookies {
		result = append(result, fmt.Sprintf("%s=%s; domain=%s; path=%s",
			c.Name, c.Value, c.Domain, c.Path,
		))
	}
	return strings.Join(result, ";\n")
}
