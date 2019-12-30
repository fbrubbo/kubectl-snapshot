package main

import (
	"bufio"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Hpa struct
type Hpa struct {
	Namespace     string
	Name          string
	ReferenceKind string
	ReferenceName string
	UsageCPU      int
	Target        int
	MinPods       int
	MaxPods       int
	Replicas      int
	Age           string
	Pods          []Pod
	Pdb           Pdb
}

// GetDeploymentKey should work for most of the cases
func (h Hpa) GetDeploymentKey() string {
	return h.Namespace + "|" + h.GetReference()
}

// GetReference ..
func (h Hpa) GetReference() string {
	return h.ReferenceKind + "/" + h.ReferenceName
}

// GetUsageAndTarget ..
func (h Hpa) GetUsageAndTarget() string {
	if h.UsageCPU == -1 {
		return "<unknown>/" + strconv.Itoa(h.Target) + "%"
	}
	return strconv.Itoa(h.UsageCPU) + "%/" + strconv.Itoa(h.Target) + "%"
}

// IsDeployment ..
func (h Hpa) IsDeployment() bool {
	return h.ReferenceKind == "Deployment"
}

// RefToDeployment ..
func (h Hpa) RefToDeployment(deployment string) bool {
	return h.IsDeployment() && h.ReferenceName == deployment
}

// ContainsPod ..
func (h Hpa) ContainsPod(pod string) bool {
	set := make(map[string]bool, len(h.Pods))
	for _, pod := range h.Pods {
		set[pod.Metadata.Name] = true
	}
	_, ok := set[pod]
	return ok
}

// CountLivenessProbes ..
func (h Hpa) CountLivenessProbes() string {
	if len(h.Pods) > 0 {
		return h.Pods[0].CountLivenessProbes()
	}
	return "N/A"
}

// CountReadinessProbes ..
func (h Hpa) CountReadinessProbes() string {
	if len(h.Pods) > 0 {
		return h.Pods[0].CountReadinessProbes()
	}
	return "N/A"
}

// CountLifecyclePreStop ..
func (h Hpa) CountLifecyclePreStop() string {
	if len(h.Pods) > 0 {
		return h.Pods[0].CountLifecyclePreStop()
	}
	return "N/A"
}

// GetLivenessProbes ..
func (h Hpa) GetLivenessProbes() string {
	if len(h.Pods) > 0 {
		return h.Pods[0].GetLivenessProbes()
	}
	return "N/A"
}

// GetReadinessProbes ..
func (h Hpa) GetReadinessProbes() string {
	if len(h.Pods) > 0 {
		return h.Pods[0].GetReadinessProbes()
	}
	return "N/A"
}

// GetLifecyclePreStop ..
func (h Hpa) GetLifecyclePreStop() string {
	if len(h.Pods) > 0 {
		return h.Pods[0].GetLifecyclePreStop()
	}
	return "N/A"
}

// RetrieveHpas executes kubectl get hpas command
// if ns is empty, then all namespaces are used
func RetrieveHpas(nsFilter string, podList []Pod) []Hpa {
	cmd := "kubectl get hpa --all-namespaces --no-headers"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	data := string(out)

	hpas := buildHpaList(data, nsFilter, podList)
	hpas = enrichHpaWithPdb(hpas)
	return hpas
}

func enrichHpaWithPdb(hpas []Hpa) (ret []Hpa) {
	//TODO: improve performance in this func
	pdbs := RetrievePdbs()
	for _, hpa := range hpas {
		if len(hpa.Pods) > 0 {
			for _, pdb := range pdbs {
				if pdb.match(hpa.Pods[0].Metadata.Labels) {
					hpa.Pdb = pdb
					break
				}
			}
		}
		ret = append(ret, hpa)
	}
	return ret
}

func buildHpaList(data string, nsFilter string, podList []Pod) (hpas []Hpa) {
	deploymentMap, replicaSetMap := buildPodMaps(podList)
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		reg, _ := regexp.Compile(`(\S*)\s*(\S*)\s*(\S*)\/(\S*)\s*((\S*)%|(<unknown>))\/(\S*)%\s*(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*`)
		txt := scanner.Text()
		groups := reg.FindStringSubmatch(txt)
		if len(groups) > 0 {
			mamespace := groups[1]
			if nsFilter == "" || nsFilter == mamespace {
				usageCPU := -1
				if groups[6] != "" {
					usageCPU, _ = strconv.Atoi(groups[6])
				}

				target, _ := strconv.Atoi(groups[8])
				minPods, _ := strconv.Atoi(groups[9])
				maxPods, _ := strconv.Atoi(groups[10])
				replicas, _ := strconv.Atoi(groups[11])
				hpa := Hpa{
					Namespace:     mamespace,
					Name:          groups[2],
					ReferenceKind: groups[3],
					ReferenceName: groups[4],
					UsageCPU:      usageCPU,
					Target:        target,
					MinPods:       minPods,
					MaxPods:       maxPods,
					Replicas:      replicas,
					Age:           groups[12],
				}
				// enrich hpa with pods
				key := hpa.Namespace + "|" + hpa.ReferenceName
				if hpa.ReferenceKind == "Deployment" {
					hpa.Pods = deploymentMap[key]
				} else if hpa.ReferenceKind == "ReplicaSet" {
					hpa.Pods = replicaSetMap[key]
				} else {
					// not implemented - return empty pods
				}
				hpas = append(hpas, hpa)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return hpas
}

func buildPodMaps(podList []Pod) (map[string][]Pod, map[string][]Pod) {
	deploymentMap := make(map[string][]Pod)
	replicaSetMap := make(map[string][]Pod)
	for _, pod := range podList {
		deployment := pod.GetDeploymentdKey()
		if deploymentMap[deployment] == nil {
			deploymentMap[deployment] = []Pod{pod}
		} else {
			deploymentMap[deployment] = append(deploymentMap[deployment], pod)
		}

		replicaset := pod.GetReplicaSetKey()
		if replicaSetMap[replicaset] == nil {
			replicaSetMap[replicaset] = []Pod{pod}
		} else {
			replicaSetMap[replicaset] = append(replicaSetMap[replicaset], pod)
		}
	}
	return deploymentMap, replicaSetMap
}
