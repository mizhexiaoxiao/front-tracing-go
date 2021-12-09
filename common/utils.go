package common

import (
	"strconv"
	"time"
)

func StringToTime(s string) (t time.Time) {
	data, _ := strconv.ParseInt(s, 10, 64)
	// 13 timestamp ms to ns
	t = time.Unix(0, data*1e6)
	return
}
