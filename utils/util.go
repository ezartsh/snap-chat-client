package utils

import (
	"regexp"
	"time"
)

func IsAlphanumeric(value string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9-_]*$`).MatchString(value)
}

func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
