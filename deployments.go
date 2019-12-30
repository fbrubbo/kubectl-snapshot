package main

import (
	"bufio"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Deployment struct
type Deployment struct {
	Namespace        string
	Name             string
	Replicas         int
	ReplicasExpected int
	UpToDate         int
	Avaliable        int
	Age              string
	Pods             []Pod
	Pdb              Pdb
}

// GetDeploymentKey should work for most of the cases
func (d Deployment) GetDeploymentKey() string {
	return d.Namespace + "|" + d.Name
}

// ContainsPod ..
func (d Deployment) ContainsPod(pod string) bool {
	set := make(map[string]bool, len(d.Pods))
	for _, pod := range d.Pods {
		set[pod.Metadata.Name] = true
	}
	_, ok := set[pod]
	return ok
}

// CountLivenessProbes ..
func (d Deployment) CountLivenessProbes() string {
	if len(d.Pods) > 0 {
		return d.Pods[0].CountLivenessProbes()
	}
	return "N/A"
}

// CountReadinessProbes ..
func (d Deployment) CountReadinessProbes() string {
	if len(d.Pods) > 0 {
		return d.Pods[0].CountReadinessProbes()
	}
	return "N/A"
}

// CountLifecyclePreStop ..
func (d Deployment) CountLifecyclePreStop() string {
	if len(d.Pods) > 0 {
		return d.Pods[0].CountLifecyclePreStop()
	}
	return "N/A"
}

// GetLivenessProbes ..
func (d Deployment) GetLivenessProbes() string {
	if len(d.Pods) > 0 {
		return d.Pods[0].GetLivenessProbes()
	}
	return "N/A"
}

// GetReadinessProbes ..
func (d Deployment) GetReadinessProbes() string {
	if len(d.Pods) > 0 {
		return d.Pods[0].GetReadinessProbes()
	}
	return "N/A"
}

// GetLifecyclePreStop ..
func (d Deployment) GetLifecyclePreStop() string {
	if len(d.Pods) > 0 {
		return d.Pods[0].GetLifecyclePreStop()
	}
	return "N/A"
}

// RetrieveDeployments executes kubectl get deployments command
// if ns is empty, then all namespaces are used
func RetrieveDeployments(nsFilter string, podList []Pod) []Deployment {
	cmd := "kubectl get deployments --all-namespaces --no-headers"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	data := string(out)
	deploys := buildDeploymentList(data, nsFilter, podList)
	deploys = enrichDeployWithPdb(deploys)
	return deploys
}

func enrichDeployWithPdb(deploys []Deployment) (ret []Deployment) {
	//TODO: improve performance in this func
	pdbs := RetrievePdbs()
	for _, deploy := range deploys {
		if len(deploy.Pods) > 0 {
			for _, pdb := range pdbs {
				if pdb.match(deploy.Pods[0].Metadata.Labels) {
					deploy.Pdb = pdb
					break
				}
			}
		}
		ret = append(ret, deploy)
	}
	return ret
}

const patternOld = `(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*`
const pattern = `(\S*)\s*(\S*)\s*(\S*)\/(\S*)\s*(\S*)\s*(\S*)\s*(\S*)\s*`

func buildDeploymentList(data string, nsFilter string, podList []Pod) []Deployment {
	podsMap := make(map[string][]Pod)
	for _, pod := range podList {
		deployment := pod.GetDeploymentdKey()
		if podsMap[deployment] == nil {
			podsMap[deployment] = []Pod{pod}
		} else {
			podsMap[deployment] = append(podsMap[deployment], pod)
		}
	}

	var deployments []Deployment
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		txt := scanner.Text()
		var reg *regexp.Regexp
		if match, _ := regexp.MatchString(pattern, txt); match {
			reg, _ = regexp.Compile(pattern)
		} else {
			reg, _ = regexp.Compile(patternOld)
		}
		groups := reg.FindStringSubmatch(txt)
		mamespace := groups[1]
		if nsFilter == "" || nsFilter == mamespace {
			replicasReady, _ := strconv.Atoi(groups[3])
			replicasExpected, _ := strconv.Atoi(groups[4])
			upToDate, _ := strconv.Atoi(groups[5])
			avaliable, _ := strconv.Atoi(groups[6])
			deployment := Deployment{
				Namespace:        mamespace,
				Name:             groups[2],
				Replicas:         replicasReady,
				ReplicasExpected: replicasExpected,
				UpToDate:         upToDate,
				Avaliable:        avaliable,
				Age:              groups[7],
			}
			deployment.Pods = podsMap[deployment.GetDeploymentKey()]
			deployments = append(deployments, deployment)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return deployments
}
