package main

import (
	"bufio"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Top struct
type Top struct {
	Namespace  string
	Pod        string
	Containers []Container
}

// Container struct
type Container struct {
	Name   string
	CPU    string
	Memory string
}

// GetDeploymentName should work for most of the cases
func (t Top) GetDeploymentName() string {
	reg, _ := regexp.Compile(`(.*)-([^-]*)-([^-]*)`)
	result := reg.FindStringSubmatch(t.Pod)
	return result[1]
}

// GetMilliCPU total pod cpu
func (t Top) GetMilliCPU() int {
	total := 0
	for _, c := range t.Containers {
		str := strings.ReplaceAll(c.CPU, "m", "")
		milli, _ := strconv.Atoi(str)
		total += milli
	}
	return total
}

// GetMiMemory returns the memory in Mi
func (t Top) GetMiMemory() int {
	total := 0
	for _, c := range t.Containers {
		reg, _ := regexp.Compile(`(\d*)(.*)`)
		groups := reg.FindStringSubmatch(c.Memory)
		memory, _ := strconv.Atoi(groups[1])
		total += memory
	}
	return total
}

// RetrieveTopMap executes kubectl get pods command
// if ns is empty, then all namespaces are used
// returns key = namespace + pod name
func RetrieveTopMap(ns string) map[string]Top {
	cmd := "kubectl top pods --all-namespaces --containers"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	data := string(out)
	return buildTopMap(data, ns)
}

func buildTopList(data string, nsFilter string) []Top {
	topMap := buildTopMap(data, nsFilter)
	var tops []Top
	for _, v := range topMap {
		tops = append(tops, v)
	}
	return tops
}

func buildTopMap(data string, nsFilter string) map[string]Top {
	r := strings.NewReader(data)
	scanner := bufio.NewScanner(r)
	top := make(map[string]Top)
	for scanner.Scan() {
		reg, _ := regexp.Compile(`(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*`)
		groups := reg.FindStringSubmatch(scanner.Text())
		mamespace := groups[1]
		if nsFilter == "" || nsFilter == mamespace {
			key := mamespace + "|" + groups[2]
			val, ok := top[key]
			if !ok {
				val = Top{
					Namespace:  mamespace,
					Pod:        groups[2],
					Containers: []Container{},
				}
			}
			val.Containers = append(val.Containers, Container{
				Name:   groups[3],
				CPU:    groups[4],
				Memory: groups[5],
			})
			top[key] = val
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return top
}
