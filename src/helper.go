package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Software struct {
	Name          string
	Repositories  map[string]string
	LatestVersion string
}

func PrintResults(softwares map[string]*Software, softwareInfo *prometheus.GaugeVec) {
	softwareInfo.Reset()

	fmt.Println("Software found on the cluster:")
	for _, software := range softwares {
		if len(software.Repositories) > 0 {
			fmt.Printf("%s:\n", software.Name)
			for repo, currentVersion := range software.Repositories {
				fmt.Printf("  repository: %s\n", repo)
				fmt.Printf("  current-version: %s\n", currentVersion)
				fmt.Printf("  latest-version: %s\n", software.LatestVersion)

				softwareInfo.WithLabelValues(software.Name, repo, currentVersion, software.LatestVersion).Set(1)
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
