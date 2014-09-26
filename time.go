package podio

import (
	"fmt"
	"time"
)

// Marshals and Unmarshals time in the common podio format.
// Podio time is always UTC.
type Time struct {
	time.Time
}

const podioLayout = "2006-01-02 15:04:05"

func (t *Time) UnmarshalJSON(from []byte) error {
	// apparently we need to trim "
	if from[0] == '"' {
		from = from[1:]
	}
	if from[len(from)-1] == '"' {
		from = from[:len(from)-1]
	}

	if string(from) == "null" {
		// on null value we set the time to the time.Time zero value
		t.Time = time.Time{}
		return nil
	}

	tm, err := time.ParseInLocation(podioLayout, string(from), time.UTC)
	if err == nil {
		t.Time = tm
	}
	return err
}

func (t *Time) MarshalJSON() ([]byte, error) {
	s := t.Format(podioLayout)
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}
