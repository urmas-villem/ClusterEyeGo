package main

import (
	"fmt"
	"time"
)

type Software struct {
	Name   string
	Images map[string]bool
}

func main() {
	for {
		softwares := GetPodInfo()
		PrintResults(softwares)
		time.Sleep(3600 * time.Second)
	}
}

func PrintResults(softwares map[string]*Software) {
	for _, software := range softwares {
		fmt.Printf("%s:\n", software.Name)
		for image := range software.Images {
			fmt.Printf("  %s\n", image)
		}
	}
}
