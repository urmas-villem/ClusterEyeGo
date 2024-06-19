package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

func FetchLatestVersionGithub(repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var releases []map[string]interface{}
	if err := json.Unmarshal(body, &releases); err != nil {
		return "", err
	}

	for _, release := range releases {
		tagName, ok := release["tag_name"].(string)
		if !ok || strings.Contains(tagName, "beta") || strings.Contains(tagName, "rc") || strings.Contains(tagName, "alpha") {
			continue
		}
		return tagName, nil
	}

	return "", fmt.Errorf("no suitable release found for %s", repo)
}

func FetchLatestVersionElastic(repo string) (string, error) {
	url := fmt.Sprintf("https://www.docker.elastic.co/r/%s", repo)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`:([0-9]+\.[0-9]+\.[0-9]+)`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var versions []string
	for _, match := range matches {
		if len(match) > 1 {
			versions = append(versions, match[1])
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found at %s", url)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i] > versions[j]
	})

	return versions[0], nil
}

func UpdateSoftwareVersions(softwares map[string]*Software, repositoryMapGithub, repositoryMapElastic map[string]string) {
	for name, software := range softwares {
		var version string
		var err error

		if repo, exists := repositoryMapGithub[name]; exists {
			version, err = FetchLatestVersionGithub(repo)
		} else if repo, exists := repositoryMapElastic[name]; exists {
			version, err = FetchLatestVersionElastic(repo)
		} else {
			fmt.Printf("Repository not found for software: %s\n", name)
			continue
		}

		if err != nil {
			fmt.Println("Error fetching latest version for", name, ":", err)
			continue
		}

		if software.LatestVersions == nil {
			software.LatestVersions = []string{version}
		} else {
			var found bool
			for _, v := range software.LatestVersions {
				if v == version {
					found = true
					break
				}
			}
			if !found {
				software.LatestVersions = append(software.LatestVersions, version)
			}
		}
	}
}
