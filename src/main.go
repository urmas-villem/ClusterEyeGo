package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

func getConfigMap(clientset *kubernetes.Clientset, name, namespace string) (map[string]string, error) {
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cm.Data, nil
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	repoConfig, err := getConfigMap(clientset, "clustereye-config", "default")
	if err != nil {
		fmt.Println("Failed to get config map:", err)
		return
	}

	repositoryMapGithub := make(map[string]string)
	repositoryMapElastic := make(map[string]string)
	for key, value := range repoConfig {
		if strings.HasPrefix(key, "github-") {
			repositoryMapGithub[strings.TrimPrefix(key, "github-")] = value
		} else if strings.HasPrefix(key, "elastic-") {
			repositoryMapElastic[strings.TrimPrefix(key, "elastic-")] = value
		}
	}

	for {
		softwares, err := GetPodInfo(repoConfig)
		if err != nil {
			fmt.Printf("Error fetching pod information: %v\n", err)
			continue
		}

		UpdateSoftwareVersions(softwares, repositoryMapGithub, repositoryMapElastic)
		PrintResults(softwares)
		time.Sleep(3600 * time.Second)
	}
}
