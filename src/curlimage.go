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

var repositoryMapGithub = map[string]string{
	"argocd":                  "argoproj/argo-cd",
	"harbor":                  "goharbor/harbor",
	"istio":                   "istio/istio",
	"kiali":                   "kiali/kiali",
	"kube-state-metrics":      "kubernetes/kube-state-metrics",
	"sealed-secrets":          "bitnami-labs/sealed-secrets",
	"alertmanager":            "prometheus/alertmanager",
	"blackbox_exporter":       "prometheus/blackbox_exporter",
	"kafka_exporter":          "danielqsj/kafka_exporter",
	"loki":                    "grafana/loki",
	"mimir":                   "grafana/mimir",
	"node_exporter":           "prometheus/node_exporter",
	"opensearch-dashboards":   "opensearch-project/OpenSearch-Dashboards",
	"opentelemetry-collector": "open-telemetry/opentelemetry-collector-contrib",
	"prometheus":              "prometheus/prometheus",
	"tempo":                   "grafana/tempo",
	"jenkins":                 "jenkinsci/jenkins",
}

var repositoryMapElastic = map[string]string{
	"filebeat": "beats/filebeat",
	"logstash": "logstash/logstash-oss",
}

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

func UpdateSoftwareVersions(softwares map[string]*Software) {
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

		software.LatestVersion = version
	}
}
