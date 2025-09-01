package tools

import "time"

func MsToTime(ms *int64) time.Time {
	if ms == nil || *ms == 0 {
		return time.Time{}
	}
	return time.Unix(0, *ms*int64(time.Millisecond)).UTC()
}