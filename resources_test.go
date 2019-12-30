package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
)

type testResource struct {
	res      Resource
	expected int
}

func TestGetMilliCPU(t *testing.T) {
	tests := []testResource{
		testResource{res: Resource{CPU: "130m", Memory: "350"}, expected: 130},
		testResource{res: Resource{CPU: "1", Memory: "450"}, expected: 1000},
		testResource{res: Resource{CPU: "0.5", Memory: "500"}, expected: 500},
		testResource{res: Resource{CPU: "1.64", Memory: "640"}, expected: 1640},
	}

	log.Infof("%+v", tests)

	for i, test := range tests {
		log.Infof("test info %d -> %+v", i, test)
		if result := test.res.GetMilliCPU(); result != test.expected {
			t.Fatalf("Test failed! %d but expected %d", result, test.expected)
		}
	}
}

func TestGetMiMemory(t *testing.T) {
	tests := []testResource{
		testResource{res: Resource{CPU: "130m", Memory: "123Mi"}, expected: 123},
		testResource{res: Resource{CPU: "1", Memory: "129M"}, expected: 123},
		testResource{res: Resource{CPU: "0.5", Memory: "128974848"}, expected: 123},
	}

	log.Infof("%+v", tests)

	for i, test := range tests {
		log.Infof("test info %d -> %+v", i, test)
		if result := test.res.GetMiMemory(); result != test.expected {
			t.Fatalf("Test failed! %d but expected %d", result, test.expected)
		}
	}
}

func TestGetDeploymentName(t *testing.T) {
	type testPodResource struct {
		res      Pod
		expected string
	}
	tests := []testPodResource{
		testPodResource{res: Pod{Metadata: Metadata{Name: "shippingservice-545f46fb7f-f4c5b"}}, expected: "shippingservice"},
		testPodResource{res: Pod{Metadata: Metadata{Name: "shipping-service-545f46fb7f-f4c5b"}}, expected: "shipping-service"},
	}

	log.Infof("%+v", tests)

	for i, test := range tests {
		log.Infof("test info %d -> %+v", i, test)
		if result := test.res.GetDeploymentName(); result != test.expected {
			t.Fatalf("Test failed! %s but expected %s", result, test.expected)
		}
	}
}

func TestStartupDuration(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/one-pod.json")
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	pr := buildPodList(str).Items[0]

	expected := 42.
	if result := pr.GetStartupDuration().Seconds(); result != expected {
		t.Fatalf("Test failed! %f but expected %f", result, expected)
	}
}

func TestStartupMissingDurationInfo(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/one-pod-missing-duration-info.json")
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	pr := buildPodList(str).Items[0]

	expected := 0.
	if result := pr.GetStartupDuration().Seconds(); result != expected {
		t.Fatalf("Test failed! %f but expected %f", result, expected)
	}
}

func TestBuildOnePod(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/one-pod.json")
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	pr := buildPodList(str).Items[0]

	ex := "shippingservice-545f46fb7f-f4c5b"
	if re := pr.Metadata.Name; re != ex {
		t.Fatalf("Test failed! %s but expected %s", re, ex)
	}
	ex = "shippingservice"
	if re := pr.GetDeploymentName(); re != ex {
		t.Fatalf("Test failed! %s but expected %s", re, ex)
	}
	ex = "shippingservice-545f46fb7f"
	if re := pr.GetReplicaSetName(); re != ex {
		t.Fatalf("Test failed! %s but expected %s", re, ex)
	}
	ex = "gke-central-pool-1-47d730e3-sh01"
	if re := pr.Spec.NodeName; re != ex {
		t.Fatalf("Test failed! %s but expected %s", re, ex)
	}
	expected := 200
	if result := pr.GetRequestsMilliCPU(); result != expected {
		t.Fatalf("Test failed! %d but expected %d", result, expected)
	}
	expected = 192
	if result := pr.GetRequestsMiMemory(); result != expected {
		t.Fatalf("Test failed! %d but expected %d", result, expected)
	}
	expected = 2200
	if result := pr.GetLimitsMilliCPU(); result != expected {
		t.Fatalf("Test failed! %d but expected %d", result, expected)
	}
	expected = 256
	if result := pr.GetLimitsMiMemory(); result != expected {
		t.Fatalf("Test failed! %d but expected %d", result, expected)
	}
	ex = "shippingservice"
	if result := pr.Metadata.Labels["app"]; result != ex {
		t.Fatalf("Test failed! %s but expected %s", result, ex)
	}
}

func TestBuildManyPods(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/many-pods.json")
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	items := buildPodList(str).Items

	ex := 23
	if l := len(items); l != ex {
		t.Fatalf("Test failed! found %d expected %d", l, ex)
	}
	for _, pod := range items {
		if &pod.Metadata.Name == nil || pod.Metadata.Name == "" {
			t.Fatalf("Test failed! Name should not be empty")
		}
	}
}

func TestBuildKubectlCmd(t *testing.T) {
	type testCmd struct {
		ns       string
		expected string
	}
	tests := []testCmd{
		testCmd{ns: "test", expected: "kubectl get pods -n test -o json"},
		testCmd{ns: "", expected: "kubectl get pods --all-namespaces -o json"},
	}

	log.Infof("%+v", tests)

	for i, test := range tests {
		log.Infof("test info %d -> %+v", i, test)
		if result := buildKubectlCmd(test.ns); result != test.expected {
			t.Fatalf("Test failed! %s but expected %s", result, test.expected)
		}
	}
}
