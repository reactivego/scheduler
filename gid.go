package scheduler

import (
	"bytes"
	"runtime"
	"strconv"
)

type Error string

const UnrecognizedGID = Error("unrecognized gid")

func (e Error) Error() string { return string(e) }

// Gid returns the numerical part of the goroutine id as a uint64. So, for:
// "goroutine 18446744073709551615" it will return uint64:18446744073709551615.
//  If for some reason the id cannot be determined, an error UnrecognizedGID
// is returned. Calling this function takes in the order of 10 microseconds.
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
