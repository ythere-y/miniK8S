package utils

import (
	"fmt"
	"time"
)

func GetAge(lasttime time.Time) time.Duration {
	return time.Now().Sub(lasttime)
}

func GetProbability(num int) float64 {
	if num <= 0 {
		fmt.Printf("num error !\n")
		return float64(0)
	}
	return float64(1) / float64(num)
}
