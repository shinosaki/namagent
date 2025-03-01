package utils

import (
	"encoding/json"
	"strconv"
	"time"
)

type UnixTime struct {
	time.Time
}

func (t *UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

func (t *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return err
	}
	t.Time = time.Unix(timestamp, 0)
	return nil
}

func (t *UnixTime) ToTime() time.Time {
	return t.Time
}

type UnixTimeMilli struct {
	time.Time
}

func (t *UnixTimeMilli) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Time.UnixMilli(), 10)), nil
}

func (t *UnixTimeMilli) UnmarshalJSON(b []byte) error {
	var timestamp int64
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return err
	}
	t.Time = time.UnixMilli(timestamp)
	return nil
}

func (t *UnixTimeMilli) ToTime() time.Time {
	return t.Time
}
