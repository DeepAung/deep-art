package types

import (
	"fmt"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (t CustomTime) MarshalJSON() ([]byte, error) {
	date := t.Time.Format(time.RFC3339)
	date = fmt.Sprintf(`"%s"`, date)
	return []byte(date), nil
}

func (t *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	date, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = date
	return nil
}

func (t *CustomTime) UnmarshalParam(param string) error {
	date, err := time.Parse("2006-01-02T15:04:05", param)
	if err != nil {
		return err
	}
	t.Time = date
	return nil
}
