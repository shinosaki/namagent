package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

// Escape invalid characters
func Escape(input string, replace string) string {
	pattern := `[\\/:;*"<>|&#!?%@+=^~]`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, replace)
}

func SaveToFile[T any](data []T, path string) error {
	if len(data) > 0 {
		var existsData []T

		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("failed to oepn file: %v", err)
		}
		defer file.Close()

		body, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		if len(body) > 0 {
			if err := json.Unmarshal(body, &existsData); err != nil {
				return err
			}
		}

		// marshalできないデータを弾く
		var failedItems []T
		for _, item := range data {
			_, err := json.Marshal(item)
			if err != nil {
				log.Println("failed to marshal item:", err)
				failedItems = append(failedItems, item)
				continue
			}
			existsData = append(existsData, item)
		}

		result, err := json.Marshal(existsData)
		if err != nil {
			return err
		}

		file.Seek(0, 0)  // seek to head
		file.Truncate(0) // remove data

		_, err = file.Write(result)
		log.Printf("Write %d data to %s", len(existsData), path)

		data = nil
		return err
	}

	return nil
}
