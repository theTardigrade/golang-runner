package time

import (
	"math"
	"time"
)

func UnixMilli() int64 {
	return int64(math.Round(float64(time.Now().UTC().UnixNano()) / 1e6))
}
