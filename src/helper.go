package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Software struct {
	Name           string
	Repositories   map[string]string
	LatestVersions []string
}

func PrintResults(softwares map[string]*Software, softwareInfo *prometheus.GaugeVec) {
	softwareInfo.Reset()

	fmt.Println("Software found on the cluster:")
	for _, software := range softwares {
		if len(software.Repositories) > 0 {
			fmt.Printf("%s:\n", software.Name)
			for repo, currentVersion := range software.Repositories {
				fmt.Printf("  repository: %s\n", repo)
				fmt.Printf("  versions: %v\n", strings.Join(software.LatestVersions, ", "))
				for _, ver := range software.LatestVersions {
					softwareInfo.WithLabelValues(software.Name, repo, currentVersion, ver).Set(1)
				}
			}
		}
	}

	fmt.Println("\nSoftware not found on the cluster:")
	for _, software := range softwares {
		if len(software.Repositories) == 0 {
			fmt.Printf("%s:\n", software.Name)
		}
	}
}

func getConfigMap(clientset *kubernetes.Clientset, name, namespace string) (map[string]map[string]string, error) {
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	configData := make(map[string]map[string]string)
	for key, jsonStr := range cm.Data {
		var result map[string]string
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			return nil, fmt.Errorf("error parsing ConfigMap key %s: %v", key, err)
		}
		configData[key] = result
	}
	return configData, nil
}

func sanityCheckGithub() bool {
	url := "https://api.github.com/rate_limit"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to check GitHub rate limits: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		return false
	}

	var rateLimitInfo struct {
		Resources struct {
			Core struct {
				Limit     int   `json:"limit"`
				Remaining int   `json:"remaining"`
				Reset     int64 `json:"reset"`
			} `json:"core"`
		} `json:"resources"`
	}

	if err := json.Unmarshal(body, &rateLimitInfo); err != nil {
		log.Printf("Error unmarshalling GitHub rate limit response: %v\n", err)
		log.Printf("Raw GitHub API response: %s\n", string(body))
		return false
	}

	if rateLimitInfo.Resources.Core.Remaining == 0 {
		resetTime := time.Unix(rateLimitInfo.Resources.Core.Reset, 0)
		remainingTime := time.Until(resetTime).Round(time.Minute)
		log.Printf("GitHub API core rate limit exceeded. Resets in %v", remainingTime)
		return true
	}

	log.Println("GitHub API core rate limit check passed.")
	return false
}
