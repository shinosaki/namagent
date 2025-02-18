package nico

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/shinosaki/namagent/internal/utils"
)

func ExtractProgramId(input string) string {
	pattern := `(lv\d+)`
	match := regexp.MustCompile(pattern).FindStringSubmatch(input)
	if len(match) == 0 {
		return ""
	}
	return match[1]
}

func FetchProgramData(programId string, client *http.Client) (*ProgramData, error) {
	if client == nil {
		client = utils.NewHttp2Client()
	}

	url := "https://live.nicovideo.jp/watch/" + programId
	log.Println("requesting url:", url)
	res, err := client.Get(url)
	if err != nil {
		log.Println("http error of fetch program data:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("read body error of fetch program data", err)
		return nil, err
	}

	pattern := `id="embedded-data" data-props="([^"]+)"`
	match := regexp.MustCompile(pattern).FindStringSubmatch(string(body))
	if len(match) == 0 {
		return nil, fmt.Errorf("does not contain embedded-data in watch page")
	}

	var result ProgramData
	data := strings.ReplaceAll(match[1], "&quot;", `"`)
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Println("failed to unmarshal program data:", err)
		return nil, err
	}

	return &result, nil
}
