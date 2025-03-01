package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

var funcs = template.FuncMap{
	"formatCookies": formatCookies,
}

func OutputTemplate(id string, tmpl string, params any) (string, error) {
	t, err := template.New(id).Funcs(funcs).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func BulkOutputTemplate(id string, tmpls []string, params any) (results []string, err error) {
	for _, tmpl := range tmpls {
		res, err := OutputTemplate(id, tmpl, params)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}

func formatCookies(cookies []*http.Cookie, separator string) string {
	var results []string
	for _, c := range cookies {
		results = append(results, fmt.Sprintf("%s=%s; domain=%s; path=%s",
			c.Name, c.Value, c.Domain, c.Path,
		))
	}
	return strings.Join(results, separator)
}
