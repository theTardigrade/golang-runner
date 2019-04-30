package main

import (
	"math"
	"time"
)

func unixMilli() int64 {
	return int64(math.Round(float64(time.Now().UTC().UnixNano()) / 1e6))
}
