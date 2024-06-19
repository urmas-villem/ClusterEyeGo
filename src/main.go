package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Software struct {
	Name          string
	Repositories  map[string]string
	LatestVersion string
}

var (
	reg = prometheus.NewRegistry()

	softwareInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "software_info",
			Help: "Information about software found in the cluster.",
		},
		[]string{"software_name", "image_repository", "image_version", "newest_image"},
	)
)

func init() {
	reg.MustRegister(softwareInfo)
}

func PrintResults(softwares map[string]*Software) {
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

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":9191", nil))
	}()

	repoConfig, err := getConfigMap(clientset, "clustereye-config", "default")
	if err != nil {
		fmt.Println("Failed to get config map:", err)
		return
	}

	for {
		softwares, err := GetPodInfo(repoConfig["github_search"], repoConfig["elastic_search"])
		if err != nil {
			fmt.Printf("Error fetching pod information: %v\n", err)
			continue
		}

		UpdateSoftwareVersions(softwares, repoConfig["github_search"], repoConfig["elastic_search"])
		PrintResults(softwares)
		time.Sleep(3600 * time.Second)
	}
}
