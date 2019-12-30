package main

import (
	"encoding/json"
	"log"
	"os/exec"
)

// PdbItems a list of Pod Disruption Budget
type PdbItems struct {
	Items []Pdb
}

// Pdb Pod Disruption Budget
type Pdb struct {
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		MinAvailable   int `json:"minAvailable"`
		MaxUnavailable int `json:"maxUnavailable"`
		Selector       struct {
			MatchLabels map[string]string `json:"matchLabels"`
		} `json:"selector"`
	} `json:"spec"`
	Status struct {
		CurrentHealthy     int `json:"currentHealthy"`
		DesiredHealthy     int `json:"desiredHealthy"`
		DisruptionsAllowed int `json:"disruptionsAllowed"`
		ExpectedPods       int `json:"expectedPods"`
	} `json:"status"`
}

func (p Pdb) match(labels map[string]string) bool {
	for k, v := range p.Spec.Selector.MatchLabels {
		if labels[k] != v {
			return false
		}
	}
	return true
}

// RetrievePdbs executes kubectl get pdb command
func RetrievePdbs() []Pdb {
	cmd := "kubectl get pdb --all-namespaces -o json"
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %s", cmd)
	}
	json := string(out)
	return buildPdbItems(json).Items
}

func buildPdbItems(str string) (pdbs PdbItems) {
	err := json.Unmarshal([]byte(str), &pdbs)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return pdbs
}
