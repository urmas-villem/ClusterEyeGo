package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

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

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go func() {
		log.Fatal(http.ListenAndServe(":9191", nil))
	}()

	for sanityCheckGithub() {
		log.Println("Retrying in 10 minutes...")
		time.Sleep(10 * time.Minute)
	}

	repoConfig, err := getConfigMap(clientset, "clustereye-config", "default")
	if err != nil {
		log.Println("Failed to get config map:", err)
		return
	}

	for {
		softwares, err := GetPodInfo(repoConfig["github_search"], repoConfig["elastic_search"])
		if err != nil {
			log.Printf("Error fetching pod information: %v\n", err)
			continue
		}

		UpdateSoftwareVersions(softwares, repoConfig["github_search"], repoConfig["elastic_search"])
		PrintResults(softwares, softwareInfo)
		time.Sleep(3600 * time.Second)
	}
}
