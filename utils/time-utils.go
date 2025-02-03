package utils

import (
	"fmt"
	"time"
)

func TimesMatch(a *time.Time, b *time.Time) bool {
	match :=
		a == nil &&
			b == nil ||
			(a != nil && b != nil && a.Unix() == b.Unix())

	return match
}

func TextToTime(textTime string) (*time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, textTime)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse time: %v", err)
	}

	return &parsedTime, nil
}
