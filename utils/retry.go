package utils

import (
	"fmt"
	"log"
	"time"
)

func Retry(toRetry func() error, numberOfTime ...int) error {
	n := 0
	var err error
	var t int
	if len(numberOfTime) != 0 {
		t = numberOfTime[0]
	} else {
		t = 10
	}
	for n < t {
		err = toRetry()
		if err == nil {
			break
		}
		n++
		log.Printf("[DEBUG] Issue on request retry: %d, -> try again", n)
		sleepDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", 2*n))
		time.Sleep(sleepDuration)
	}
	return err
}
