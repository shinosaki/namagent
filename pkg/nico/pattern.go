package nico

import "regexp"

func ExtractId(v string) (id string) {
	re := regexp.MustCompile(`(lv\d+)`)
	match := re.FindStringSubmatch(v)

	if len(match) == 2 {
		id = match[1]
	}

	return id
}
