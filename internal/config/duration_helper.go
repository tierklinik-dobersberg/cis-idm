package config

import (
	"encoding/json"
	"time"
)

type JSONDuration time.Duration

func (d *JSONDuration) UnmarshalJSON(blob []byte) error {
	var s string
	if err := json.Unmarshal(blob, &s); err != nil {
		return err
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = JSONDuration(duration)

	return nil
}

func (d *JSONDuration) MarshalJSON() ([]byte, error) {
	str := time.Duration(*d).String()
	return json.Marshal(str)
}

func (d *JSONDuration) AsDuration() time.Duration {
	return time.Duration(*d)
}
