package main

import (
	"time"
)

func main() {
	for {
		GetPodInfo()
		time.Sleep(3600 * time.Second)
	}
}
