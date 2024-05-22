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

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	date, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = date
	return
}
