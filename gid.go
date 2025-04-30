package scheduler

import (
	"bytes"
	"runtime"
	"strconv"
)

const UnrecognizedGID = Error("unrecognized gid")

// Gid returns the numerical part of the goroutine id as a uint64. So, for:
// "goroutine 18446744073709551615" it will return uint64:18446744073709551615.
// If for some reason the id cannot be determined, the function panics with either
// UnrecognizedGID or the parsing error. The function works by getting a stack trace
// of the current goroutine, extracting the goroutine ID prefix, and parsing it into an integer.
// Calling this function takes in the order of 10 microseconds.
func Gid() uint64 {
	b := make([]byte, 32)
	if runtime.Stack(b, false) < 12 {
		panic(UnrecognizedGID)
	}
	b = b[10:]
	idx := bytes.IndexByte(b, ' ')
	if idx < 1 {
		panic(UnrecognizedGID)
	}
	gid, err := strconv.ParseUint(string(b[:idx]), 10, 64)
	if err != nil {
		panic(err)
	}
	return gid
}
