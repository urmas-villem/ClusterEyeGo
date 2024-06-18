package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var repositoryMap = map[string]string{
	"argocd":       "argoproj/argo-cd",
	"jenkins":      "jenkinsci/jenkins",
	"prometheus":   "prometheus/prometheus",
	"alertmanager": "prometheus/alertmanager",
	"clustereye":   "huxlee/clustereye",
	"logstash":     "elastic/logstash",
}

func FetchLatestVersion(repo string) (string, error) {
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
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return "", err
	}

	for _, release := range releases {
		tagName, ok := release["tag_name"].(string)
		if !ok {
			continue
		}
		if strings.Contains(tagName, "beta") || strings.Contains(tagName, "rc") || strings.Contains(tagName, "alpha") {
			continue
		}
		return tagName, nil
	}

	return "", fmt.Errorf("no suitable release found for %s", repo)
}

func UpdateSoftwareVersions(softwares map[string]*Software) {
	for name, software := range softwares {
		repo, exists := repositoryMap[name]
		if !exists {
			fmt.Printf("Repository not found for software: %s\n", name)
			continue
		}

		version, err := FetchLatestVersion(repo)
		if err != nil {
			fmt.Println("Error fetching latest version for", name, ":", err)
			continue
		}

		software.ImageFound = version
	}
}
