package globaltime

import "time"

// Now returns the current time. This wrapper is useful for unit testing.
func Now() time.Time {
	return time.Now()
}

// UnixNano returns the current time as a Unix time, the number of nanoseconds
// elapsed since January 1, 1970 UTC.
func UnixNano() int64 {
	return time.Now().UnixNano()
}

// Since returns the time elapsed since t.
func Since(t time.Time) time.Duration {
	return time.Since(t)
}
