package scheduler

import (
	"bytes"
	"runtime"
)

// Gid returns the numerical part of the goroutine id as a string. So, for:
// "goroutine 18446744073709551615" it will return "18446744073709551615". If
// for some reason the id cannot be detected, an empty string is returned.
// Calling this function takes in the order of 10 microseconds.
func Gid() string {
	b := make([]byte, 32)
	if runtime.Stack(b, false) < 12 {
		return ""
	}
	b = b[10:]
	idx := bytes.IndexByte(b, ' ')
	if idx < 1 {
		return ""
	}
	return string(b[:idx])
}
