package utils

import "time"

func TimesMatch(a *time.Time, b *time.Time) bool {
	match :=
		a == nil &&
			b == nil ||
			(a != nil && b != nil && a.Unix() == b.Unix())

	return match
}
