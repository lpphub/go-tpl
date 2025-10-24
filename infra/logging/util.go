package logging

import (
	"strconv"
	"time"
)

func GenerateLogID() string {
	return strconv.FormatUint(uint64(time.Now().UnixNano())&0x7FFFFFFF|0x80000000, 10)
}
