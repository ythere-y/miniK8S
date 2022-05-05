package utils

import "time"

func GetAge(lasttime time.Time) time.Duration {
	return time.Now().Sub(lasttime)
}
