package main

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var softwareFilters = []string{"argocd", "jenkins", "prometheus", "alertmanager", "clustereye", "logstash"}

func GetPodInfo() map[string]*Software {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	softwares := make(map[string]*Software)
	for _, filter := range softwareFilters {
		softwares[filter] = &Software{Name: filter, Images: make(map[string]bool)}
	}

	for _, pod := range pods.Items {
		for _, filter := range softwareFilters {
			if strings.Contains(pod.Name, filter) {
				for _, container := range pod.Spec.Containers {
					if strings.Contains(container.Image, filter) {
						softwares[filter].Images[container.Image] = true
					}
				}
			}
		}
	}

	return softwares
}
