package recorder

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/shinosaki/namagent/internal/recorder/types"
	"github.com/shinosaki/namagent/utils"
)

func ExtractProgramId(input string) string {
	pattern := `(lv\d+)`
	match := regexp.MustCompile(pattern).FindStringSubmatch(input)
	if len(match) == 0 {
		return ""
	}
	return match[1]
}

func FetchProgramData(programId string, client *http.Client) (*types.ProgramData, error) {
	if client == nil {
		client = utils.NewHttp2Client()
	}

	res, err := client.Get("https://live.nicovideo.jp/watch/" + programId)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	pattern := `id="embedded-data" data-props="([^"]+)"`
	match := regexp.MustCompile(pattern).FindStringSubmatch(string(body))
	if len(match) == 0 {
		return nil, fmt.Errorf("does not contain embedded-data in watch page")
	}

	var result types.ProgramData
	data := strings.ReplaceAll(match[1], "&quot;", `"`)
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
