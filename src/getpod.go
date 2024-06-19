package main

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var softwareFilters = []string{"argocd", "jenkins", "prometheus", "alertmanager", "clustereye", "logstash"}

func GetPodInfo() (map[string]*Software, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	softwares := make(map[string]*Software)
	for _, filter := range softwareFilters {
		softwares[filter] = &Software{Name: filter, Repositories: make(map[string]string)}
	}

	for _, pod := range pods.Items {
		for _, filter := range softwareFilters {
			if strings.Contains(pod.Name, filter) {
				for _, container := range pod.Spec.Containers {
					if strings.Contains(container.Image, filter) {
						parts := strings.Split(container.Image, ":")
						if len(parts) < 2 {
							return nil, fmt.Errorf("error: image '%s' is missing a version tag", container.Image)
						}
						repository := parts[0]
						version := parts[1]
						softwares[filter].Repositories[repository] = version
					}
				}
			}
		}
	}

	return softwares, nil
}
