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

func GetPodInfo() {
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
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	for _, pod := range pods.Items {
		if containsAny(pod.Name, softwareFilters) {
			fmt.Printf("Pod Name: %s\n", pod.Name)
			for _, container := range pod.Spec.Containers {
				fmt.Printf("  Container Name: %s, Image: %s\n", container.Name, container.Image)
			}
		}
	}
}

func containsAny(podName string, filters []string) bool {
	for _, filter := range filters {
		if strings.Contains(podName, filter) {
			return true
		}
	}
	return false
}
