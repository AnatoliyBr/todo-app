package entity

import (
	"bytes"
	"time"
)

type TimeISO struct {
	*time.Time
}

func (t *TimeISO) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(time.DateTime)), nil
}

func (t *TimeISO) UnMarshalJSON(data []byte) error {

	// from json doc: by convention, unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op
	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	time, err := time.Parse(time.DateTime, string(data))
	if err != nil {
		return err
	}

	t = &TimeISO{&time}
	return nil
}
