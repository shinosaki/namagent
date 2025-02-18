package nicoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/shinosaki/namagent/internal/utils"
)

func FetchRecentPrograms(isBulkFetch bool, client *http.Client) ([]RecentProgram, error) {
	log.Println("fetch recent programs")

	if client == nil {
		client = utils.NewHttp2Client()
	}

	if !isBulkFetch {
		return recentPrograms(0, client)
	}

	var (
		offset          = 0
		result          []RecentProgram
		MAX_DATA_LENGTH = 70 // Max length of recent programs endpoint
	)
	for {
		if data, err := recentPrograms(offset, client); err != nil {
			return nil, err
		} else {
			result = append(result, data...)

			if len(data) < MAX_DATA_LENGTH {
				return result, nil
			}
		}

		offset++
	}
}

func recentPrograms(offset int, client *http.Client) ([]RecentProgram, error) {
	var (
		endpoint = "https://live.nicovideo.jp/front/api/pages/recent/v1/programs"
		params   = fmt.Sprintf("?offset=%d&sortOrder=recentDesc", offset)
	)

	// HTTP2 Client
	if client == nil {
		client = utils.NewHttp2Client()
	}

	res, err := client.Get(endpoint + params)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data APIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if data.Meta.Status != 200 {
		return nil, fmt.Errorf("%s", data.Meta.ErrorCode)
	}

	var result []RecentProgram
	if err := json.Unmarshal(data.Data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
