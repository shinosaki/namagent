package nico

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func IsLogin(client *http.Client) bool {
	res, err := client.Get("https://www.nicovideo.jp")
	if err != nil {
		return false
	}

	authFlag := res.Header.Get("x-niconico-authflag")
	return authFlag == "1"
}

func FetchProgramData(programId string, client *http.Client) (result ProgramData, err error) {
	res, err := client.Get("https://live.nicovideo.jp/watch/" + programId)
	if err != nil {
		return result, fmt.Errorf("http error %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("read error %v", err)
	}

	re := regexp.MustCompile(`id="embedded-data" data-props="([^"]+)"`)
	match := re.FindStringSubmatch(string(body))
	if len(match) != 2 {
		return result, fmt.Errorf("not contain embedded-data in watch page")
	}

	embeddedData := strings.ReplaceAll(match[1], "&quot;", `"`)

	if err := json.Unmarshal([]byte(embeddedData), &result); err != nil {
		return result, fmt.Errorf("unmarshal error %v", err)
	}

	return result, nil
}
