package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// PodList struct
type PodList struct {
	Items []Pod
}

// Pod struct
type Pod struct {
	Metadata Metadata
	Spec     Spec
	Status   struct {
		Conditions        []Condition `json:"conditions"`
		ContainerStatuses []struct {
			RestartCount int `json:"restartCount"`
		} `json:"containerStatuses"`
		Phase string
	}
	Top Top
}

// Condition struct
type Condition struct {
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Status             string    `json:"status"`
	Type               string    `json:"type"`
}

// Metadata struct
type Metadata struct {
	Name            string
	Namespace       string
	Labels          map[string]string
	OwnerReferences []struct {
		Kind string
		Name string
	}
}

// Spec struct
type Spec struct {
	NodeName   string
	Containers []struct {
		Name      string
		Lifecycle struct {
			PreStop struct {
				Exec struct {
					Command []string
				}
				HTTPGet struct {
					Path string
				}
			}
		}
		LivenessProbe struct {
			HTTPGet struct {
				Path string `json:"path"`
			} `json:"httpGet,omitempty"`
			Exec struct {
				Command []string `json:"command"`
			} `json:"exec,omitempty"`
			FailureThreshold    int `json:"failureThreshold"`
			InitialDelaySeconds int `json:"initialDelaySeconds"`
			PeriodSeconds       int `json:"periodSeconds"`
			SuccessThreshold    int `json:"successThreshold"`
			TimeoutSeconds      int `json:"timeoutSeconds"`
		} `json:"livenessProbe,omitempty"`
		ReadinessProbe struct {
			HTTPGet struct {
				Path string `json:"path"`
			} `json:"httpGet,omitempty"`
			Exec struct {
				Command []string `json:"command"`
			} `json:"exec,omitempty"`
			FailureThreshold    int `json:"failureThreshold"`
			InitialDelaySeconds int `json:"initialDelaySeconds"`
			PeriodSeconds       int `json:"periodSeconds"`
			SuccessThreshold    int `json:"successThreshold"`
			TimeoutSeconds      int `json:"timeoutSeconds"`
		} `json:"readinessProbe,omitempty"`
		Resources struct {
			Requests Resource
			Limits   Resource
		}
	}
}

// Resource struct
type Resource struct {
	CPU    string
	Memory string
}

// GetMilliCPU returns the CPU in MilliCPU
func (r Resource) GetMilliCPU() int {
	return String2MilliCPU(r.CPU)
}

// GetMiMemory returns the memory in Mi
func (r Resource) GetMiMemory() int {
	return String2MiMemory(r.Memory)
}

// GetPodKey returns <namespace>-<pod name>
func (p Pod) GetPodKey() string {
	return p.Metadata.Namespace + "|" + p.Metadata.Name
}

// GetDeploymentdKey returns <namespace>-<pod name>
func (p Pod) GetDeploymentdKey() string {
	return p.Metadata.Namespace + "|" + p.GetDeploymentName()
}

// GetReplicaSetKey returns <namespace>-<pod name>
func (p Pod) GetReplicaSetKey() string {
	return p.Metadata.Namespace + "|" + p.GetReplicaSetName()
}

// GetStartupDuration returns the best effort for geting startup time (ready - schedule), 0 otherwise
func (p Pod) GetStartupDuration() time.Duration {
	restartCount := 0
	for _, cs := range p.Status.ContainerStatuses {
		restartCount = restartCount + cs.RestartCount
	}

	ready := p.findStatusCondition(func(c Condition) bool { return c.Status == "True" && c.Type == "Ready" })
	podScheduled := p.findStatusCondition(func(c Condition) bool { return c.Status == "True" && c.Type == "PodScheduled" })
	if restartCount > 0 || ready.Status == "NA" || podScheduled.Status == "NA" {
		// if pod has restarts or no info in any of the containers statuses, that means the timestamps in the statuses do not represent startup duration
		return time.Duration(0)
	}
	diff := ready.LastTransitionTime.Sub(podScheduled.LastTransitionTime)
	if diff > time.Hour {
		// if diff is too big (guessing 1+ hour), that probably means the timestamps in the statuses do not represent startup duration. Eg, the pod may became unhealty and then healty again
		return time.Duration(0)
	}
	return diff
}

func (p Pod) findStatusCondition(test func(Condition) bool) Condition {
	for _, c := range p.Status.Conditions {
		if test(c) {
			return c
		}
	}
	return Condition{Status: "NA"}
}

//senninha-quotation-redis-slave-0
// zoidberg-pentaho-report-1572104400-rklgx
const stafulsetPattern = `(.*)-(\d*)`
const deploymentPattern = `(.*)-([^-]*)-([^-]*)`
const jobPattern = `(.*)-([^-]*)`

// GetDeploymentName should work for most of the cases
func (p Pod) GetDeploymentName() string {
	name := p.Metadata.Name
	var reg *regexp.Regexp
	if match, _ := regexp.MatchString(deploymentPattern, name); match {
		reg, _ = regexp.Compile(deploymentPattern)
	} else if match, _ := regexp.MatchString(stafulsetPattern, name); match {
		reg, _ = regexp.Compile(stafulsetPattern)
	} else if p.Metadata.OwnerReferences != nil && p.Metadata.OwnerReferences[0].Kind == "Job" {
		reg, _ = regexp.Compile(jobPattern)
	}
	result := reg.FindStringSubmatch(name)
	return result[1]
}

// GetReplicaSetName should work for most of the cases
func (p Pod) GetReplicaSetName() string {
	if p.Metadata.OwnerReferences == nil {
		return "<no-references>"
	}
	return p.Metadata.OwnerReferences[0].Name
}

// GetRequestsMilliCPU total
func (p Pod) GetRequestsMilliCPU() int {
	total := 0
	for _, c := range p.Spec.Containers {
		total += c.Resources.Requests.GetMilliCPU()
	}
	return total
}

// GetTopMilliCPU total
func (p Pod) GetTopMilliCPU() int {
	return p.Top.GetMilliCPU()
}

// GetUsageCPU %
func (p Pod) GetUsageCPU() float32 {
	top := float32(p.GetTopMilliCPU())
	requests := float32(p.GetRequestsMilliCPU())
	if top == 0 && requests != 0 {
		return 0
	} else if requests == 0 {
		return 100
	}
	return top / requests * 100
}

// GetRequestsMiMemory total
func (p Pod) GetRequestsMiMemory() int {
	total := 0
	for _, c := range p.Spec.Containers {
		total += c.Resources.Requests.GetMiMemory()
	}
	return total
}

// GetTopMiMemory total
func (p Pod) GetTopMiMemory() int {
	return p.Top.GetMiMemory()
}

// GetUsageMemory %
func (p Pod) GetUsageMemory() float32 {
	top := float32(p.GetTopMiMemory())
	requests := float32(p.GetRequestsMiMemory())
	if top == 0 && requests != 0 {
		return 0
	} else if requests == 0 {
		return 100
	}
	return top / requests * 100
}

// GetLimitsMilliCPU total
func (p Pod) GetLimitsMilliCPU() int {
	total := 0
	for _, c := range p.Spec.Containers {
		total += c.Resources.Limits.GetMilliCPU()
	}
	return total
}

// GetLimitsMiMemory total
func (p Pod) GetLimitsMiMemory() int {
	total := 0
	for _, c := range p.Spec.Containers {
		total += c.Resources.Limits.GetMiMemory()
	}
	return total
}

// CountLivenessProbes ..
func (p Pod) CountLivenessProbes() string {
	numContainers := len(p.Spec.Containers)
	numLiveness := 0
	for _, c := range p.Spec.Containers {
		if c.LivenessProbe.HTTPGet.Path != "" || c.LivenessProbe.Exec.Command != nil {
			numLiveness = numLiveness + 1
		}
	}
	return fmt.Sprintf("%d/%d", numLiveness, numContainers)
}

// CountReadinessProbes ..
func (p Pod) CountReadinessProbes() string {
	numContainers := len(p.Spec.Containers)
	numReadiness := 0
	for _, c := range p.Spec.Containers {
		if c.ReadinessProbe.HTTPGet.Path != "" || c.ReadinessProbe.Exec.Command != nil {
			numReadiness = numReadiness + 1
		}
	}
	return fmt.Sprintf("%d/%d", numReadiness, numContainers)
}

// CountLifecyclePreStop ..
func (p Pod) CountLifecyclePreStop() string {
	numContainers := len(p.Spec.Containers)
	preStop := 0
	for _, c := range p.Spec.Containers {
		if c.Lifecycle.PreStop.HTTPGet.Path != "" || c.Lifecycle.PreStop.Exec.Command != nil {
			preStop = preStop + 1
		}
	}
	return fmt.Sprintf("%d/%d", preStop, numContainers)
}

// GetLivenessProbes ..
func (p Pod) GetLivenessProbes() string {
	str := ""
	for _, c := range p.Spec.Containers {
		if len(str) > 0 {
			str += "\n"
		}
		str += c.Name + " {"
		if c.LivenessProbe.HTTPGet.Path != "" {
			str += "HttpGet: " + c.LivenessProbe.HTTPGet.Path
		} else if c.LivenessProbe.Exec.Command != nil {
			str += "Exec: " + strings.Join(c.LivenessProbe.Exec.Command, " ")
		}
		str += "}"
	}
	return str
}

// GetReadinessProbes ..
func (p Pod) GetReadinessProbes() string {
	str := ""
	for _, c := range p.Spec.Containers {
		if len(str) > 0 {
			str += "\n"
		}
		str += c.Name + " {"
		if c.ReadinessProbe.HTTPGet.Path != "" {
			str += "HttpGet: " + c.ReadinessProbe.HTTPGet.Path
		} else if c.ReadinessProbe.Exec.Command != nil {
			str += "Exec: " + strings.Join(c.ReadinessProbe.Exec.Command, " ")
		}
		str += "}"
	}
	return str
}

// GetLifecyclePreStop ..
func (p Pod) GetLifecyclePreStop() string {

	str := ""
	for _, c := range p.Spec.Containers {
		if len(str) > 0 {
			str += "\n"
		}
		str += c.Name + " {"
		if c.Lifecycle.PreStop.HTTPGet.Path != "" {
			str += "HttpGet: " + c.Lifecycle.PreStop.HTTPGet.Path
		} else if c.Lifecycle.PreStop.Exec.Command != nil {
			str += "Exec: " + strings.Join(c.Lifecycle.PreStop.Exec.Command, " ")
		}
		str += "}"
	}
	return str
}

// RetrievePods executes kubectl get pods command and return only status.phase == "Running" pods
// if ns is empty, then all namespaces are used
func RetrievePods(ns string) []Pod {
	cmd := buildKubectlCmd(ns)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	json := string(out)
	pods := buildPodList(json).Items
	return enrichPodsWithTopInfoAndFilterRunning(pods, ns)
}

func enrichPodsWithTopInfoAndFilterRunning(pods []Pod, ns string) []Pod {
	var podList []Pod
	topMap := RetrieveTopMap(ns)
	for _, pod := range pods {
		if pod.Status.Phase == "Running" {
			if top, ok := topMap[pod.GetPodKey()]; ok {
				pod.Top = top
			}
			podList = append(podList, pod)
		}
	}
	return podList
}

func buildKubectlCmd(ns string) string {
	cmd := fmt.Sprintf("kubectl get pods --all-namespaces -o json")
	if ns != "" {
		cmd = fmt.Sprintf("kubectl get pods -n %s -o json", ns)
	}
	return cmd
}

func buildPodList(str string) PodList {
	pods := PodList{}
	err := json.Unmarshal([]byte(str), &pods)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return pods
}
