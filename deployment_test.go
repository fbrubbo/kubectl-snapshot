package main

import (
	"io/ioutil"
	"testing"
)

func TestBuildDeploymentList(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/deployment.txt")
	if err != nil {
		t.Fatal(err)
	}
	data := string(b)

	deployments := buildDeploymentList(data, "", []Pod{})
	ex := 116
	if l := len(deployments); l != ex {
		t.Fatalf("Test failed! found %d expected %d", l, ex)
	}
	for _, deployment := range deployments {
		if deployment.Namespace == "" || deployment.Name == "" || deployment.Age == "" {
			t.Fatalf("Test failed! deployment must have all info")
		}
	}
}

func TestBuildDeploymentV1(t *testing.T) {
	data := `qdc-web-test                  qdc-web-test                                         0      0      0      0      169d`
	deployments := buildDeploymentList(data, "", []Pod{})

	deployment := deployments[0]
	if deployment.Namespace != "qdc-web-test" ||
		deployment.Name != "qdc-web-test" ||
		deployment.Replicas != 0 ||
		deployment.ReplicasExpected != 0 ||
		deployment.UpToDate != 0 ||
		deployment.Avaliable != 0 ||
		deployment.Age != "169d" {
		t.Fatalf("Test failed! deployment does not match data")
	}
}

func TestBuildDeploymentV2(t *testing.T) {
	data := `istio-system               grafana                                    1/1     1            1           133d`
	deployments := buildDeploymentList(data, "", []Pod{})

	deployment := deployments[0]
	if deployment.Namespace != "istio-system" ||
		deployment.Name != "grafana" ||
		deployment.Replicas != 1 ||
		deployment.ReplicasExpected != 1 ||
		deployment.UpToDate != 1 ||
		deployment.Avaliable != 1 ||
		deployment.Age != "133d" {
		t.Fatalf("Test failed! deployment does not match data")
	}
}
