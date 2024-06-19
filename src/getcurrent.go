package main

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetPodInfo(githubSearch, elasticSearch map[string]string) (map[string]*Software, error) {
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
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			parts := strings.Split(container.Image, ":")
			if len(parts) < 2 {
				continue
			}
			repository := parts[0]
			currentVersion := parts[1]
			key := pod.Name + ":" + repository

			if software, exists := softwares[key]; exists {
				software.CurrentVersions = append(software.CurrentVersions, currentVersion)
			} else {
				softwares[key] = &Software{
					Name:            pod.Name,
					Repositories:    map[string]string{repository: currentVersion},
					CurrentVersions: []string{currentVersion},
					LatestVersions:  []string{},
				}
			}
		}
	}

	return softwares, nil
}
