package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestBuildManyPodTop(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/top-many-pods.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := string(b)

	tops := buildTopList(data, "")
	ex := 69
	if l := len(tops); l != ex {
		t.Fatalf("Test failed! found %d expected %d", l, ex)
	}
	for _, pod := range tops {
		if &pod.Pod == nil || pod.Pod == "" {
			t.Fatalf("Test failed! Pod Name should not be empty")
		}

		if &pod.Containers == nil || len(pod.Containers) == 0 {
			t.Fatalf("Test failed! Pod must have containers")
		}

		for _, c := range pod.Containers {
			if &c.Name == nil || c.Name == "" ||
				&c.CPU == nil || c.CPU == "" ||
				&c.Memory == nil || c.Memory == "" {
				t.Fatalf("Test failed! Containaer data must be set")
			}
		}
	}
}

func TestBuildPodTopDefaultNamespace(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/top-many-pods.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := string(b)

	tops := buildTopList(data, "default")
	ex := 23
	if l := len(tops); l != ex {
		t.Fatalf("Test failed! found %d expected %d", l, ex)
	}
}

func TestBuildOnePodTop(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/top-one-pod.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := string(b)

	top := buildTopList(data, "")[0]
	expectedCPU := 32
	if cpu := top.GetMilliCPU(); cpu != expectedCPU {
		t.Fatalf("Test failed! %d but expected %d", cpu, expectedCPU)
	}
	expectedMemory := 25
	if mem := top.GetMiMemory(); mem != expectedMemory {
		t.Fatalf("Test failed! %d but expected %d", mem, expectedMemory)
	}
}
