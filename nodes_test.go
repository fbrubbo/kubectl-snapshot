package main

import (
	"io/ioutil"
	"testing"
)

func TestBuildNodeList(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/node1.json")
	if err != nil {
		t.Fatal(err)
	}
	str := string(b)
	node := buildNodeList(str).Items[0]

	if node.GetInstanceType() != "n1-highmem-8" ||
		node.GetNodepool() != "pool-1" ||
		node.GetZone() != "us-central1-b" ||
		node.GetAllocatableMilliCPU() != 7910 ||
		node.GetAllocatableMiMemory() != 47399 ||
		node.GetAllocatablePods() != 110 {
		t.Fatalf("Test failed! %+v", node)
	}
}

func TestBuildNodeListSize(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/nodes.json")
	if err != nil {
		t.Fatal(err)
	}
	str := string(b)
	nodes := buildNodeList(str).Items

	if len(nodes) != 4 {
		t.Fatalf("Test failed! %+v", nodes)
	}
}
