package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
)

// NodeItems struct ..
type NodeItems struct {
	Items []Node
}

// Node struct ..
type Node struct {
	Metadata struct {
		Labels struct {
			InstanceType string `json:"beta.kubernetes.io/instance-type"`
			Nodepool     string `json:"cloud.google.com/gke-nodepool"`
			Zone         string `json:"failure-domain.beta.kubernetes.io/zone"`
		} `json:"labels"`
		Name string `json:"name"`
	} `json:"metadata"`
	Status struct {
		Allocatable struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
			Pods   string `json:"pods"`
		} `json:"allocatable"`
	} `json:"status"`
	Pods []Pod
}

// GetName ..
func (n Node) GetName() string {
	return n.Metadata.Name
}

// GetInstanceType ..
func (n Node) GetInstanceType() string {
	return n.Metadata.Labels.InstanceType
}

// GetNodepool ..
func (n Node) GetNodepool() string {
	return n.Metadata.Labels.Nodepool
}

// GetZone ..
func (n Node) GetZone() string {
	return n.Metadata.Labels.Zone
}

// GetAllocatableMilliCPU ..
func (n Node) GetAllocatableMilliCPU() int {
	return String2MilliCPU(n.Status.Allocatable.CPU)
}

// GetAllocatableMiMemory ..
func (n Node) GetAllocatableMiMemory() int {
	return String2MiMemory(n.Status.Allocatable.Memory)
}

// GetAllocatablePods ..
func (n Node) GetAllocatablePods() int {
	numPods, _ := strconv.Atoi(n.Status.Allocatable.Pods)
	return numPods
}

// RetrieveNodes executes kubectl get pods command
// if ns is empty, then all namespaces are used
func RetrieveNodes(podList []Pod) (ret []Node) {
	cmd := "kubectl get nodes -o json"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	json := string(out)
	nodes := buildNodeList(json).Items
	podMap := make(map[string][]Pod)
	for _, pod := range podList {
		nodeName := pod.Spec.NodeName
		if pods, ok := podMap[nodeName]; ok {
			podMap[nodeName] = append(pods, pod)
		} else {
			podMap[nodeName] = []Pod{pod}
		}
	}
	for _, node := range nodes {
		node.Pods = podMap[node.GetName()]
		ret = append(ret, node)
	}
	return
}

func buildNodeList(str string) NodeItems {
	nodes := NodeItems{}
	err := json.Unmarshal([]byte(str), &nodes)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return nodes
}
