package main

import (
	"io/ioutil"
	"testing"
)

func TestBuildPdb(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/pdb.json")
	if err != nil {
		t.Fatal(err)
	}
	str := string(b)
	pdbs := buildPdbItems(str)

	pdb := pdbs.Items[0]
	if pdb.Spec.Selector.MatchLabels["app"] != "adservice" {
		t.Fatalf("Test failed! %+v", pdb.Spec.Selector.MatchLabels)
	}

	pdb = pdbs.Items[1]
	if pdb.Spec.Selector.MatchLabels["app"] != "adservice2" || pdb.Spec.Selector.MatchLabels["xyz"] != "abc2" {
		t.Fatalf("Test failed! %+v", pdb.Spec.Selector.MatchLabels)
	}
}

func TestPdbMatch(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/pdb.json")
	if err != nil {
		t.Fatal(err)
	}
	str := string(b)
	pdbs := buildPdbItems(str)

	pdb := pdbs.Items[0]
	labels := make(map[string]string)
	labels["app"] = "adservice"
	if !pdb.match(labels) {
		t.Fatalf("Test failed to match! %+v", pdb)
	}

	pdb = pdbs.Items[1]
	labels = make(map[string]string)
	labels["app"] = "adservice"
	if pdb.match(labels) {
		t.Fatalf("Test failed to match! %+v", pdb)
	}

	labels["app"] = "adservice2"
	labels["xyz"] = "abc2"
	if !pdb.match(labels) {
		t.Fatalf("Test failed to match! %+v", pdb)
	}
}
