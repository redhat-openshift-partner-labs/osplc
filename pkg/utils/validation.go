package utils

import (
	"regexp"
	"time"
)

func IsValidTimezone(tz string) bool {
	_, err := time.LoadLocation(tz)
	return err == nil
}

func IsValidUptime(uptime string) bool {
	pattern := `^([0-9]+(\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$`
	matched, _ := regexp.MatchString(pattern, uptime)
	return matched
}
