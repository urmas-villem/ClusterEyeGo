package main

import (
	"fmt"
	"time"
)

type Software struct {
	Name   string
	Images map[string]string
}

func main() {
	for {
		softwares, err := GetPodInfo()
		if err != nil {
			fmt.Printf("Error fetching pod information: %v\n", err)
			continue
		}
		PrintResults(softwares)
		time.Sleep(3600 * time.Second)
	}
}

func PrintResults(softwares map[string]*Software) {
	for _, software := range softwares {
		fmt.Printf("%s:\n", software.Name)
		for repo, version := range software.Images {
			fmt.Printf("  repository: %s\n", repo)
			fmt.Printf("  image-version: %s\n", version)
		}
	}
}
