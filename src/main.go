package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		results := GetPodInfo()
		for software, images := range results {
			fmt.Printf("%s:\n", software)
			for _, image := range images {
				fmt.Printf("  %s\n", image)
			}
		}
		time.Sleep(3600 * time.Second)
	}
}
