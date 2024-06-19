package main

import (
	"fmt"
	"time"
)

type Software struct {
	Name          string
	Repositories  map[string]string
	LatestVersion string
}

func PrintResults(softwares map[string]*Software) {
	for _, software := range softwares {
		fmt.Printf("%s:\n", software.Name)
		for repo, currentVersion := range software.Repositories {
			fmt.Printf("  repository: %s\n", repo)
			fmt.Printf("  current-version: %s\n", currentVersion)
			fmt.Printf("  latest-version: %s\n", software.LatestVersion)
		}
	}
}

func main() {
	for {
		softwares, err := GetPodInfo()
		if err != nil {
			fmt.Printf("Error fetching pod information: %v\n", err)
			continue
		}

		UpdateSoftwareVersions(softwares)
		PrintResults(softwares)
		time.Sleep(3600 * time.Second)
	}
}
