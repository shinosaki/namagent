package utils

import (
	"encoding/json"
	"log"
)

func Unmarshaller[T any](data []byte) *T {
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		log.Println("Unmarshaller Error:", err)
		return nil
	}
	return &result
}
