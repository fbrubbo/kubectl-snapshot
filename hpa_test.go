package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestBuildHpaList(t *testing.T) {
	b, err := ioutil.ReadFile("test-data/hpa.txt")
	if err != nil {
		log.Fatal(err)
	}
	data := string(b)

	hpas := buildHpaList(data, "", []Pod{})
	ex := 18
	if l := len(hpas); l != ex {
		t.Fatalf("Test failed! found %d expected %d", l, ex)
	}
	for _, hpa := range hpas {
		if hpa.Namespace == "" || hpa.Name == "" || hpa.ReferenceKind == "" || hpa.ReferenceName == "" || hpa.Age == "" {
			t.Fatalf("Test failed! hpa must have all info")
		}
	}
}

func TestBuildHpaMap(t *testing.T) {
	data := `default        nginx-1-hpa                                             Deployment/nginx-1                 <unknown>/80%   1         5         3          33d
default        paymentservice                                          Deployment/paymentservice          4%/80%          2         20        2          87d`
	hpas := buildHpaList(data, "", []Pod{})

	hpa := hpas[0]
	if hpa.Namespace != "default" ||
		hpa.Name != "nginx-1-hpa" ||
		hpa.ReferenceKind != "Deployment" ||
		hpa.ReferenceName != "nginx-1" ||
		hpa.UsageCPU != -1 ||
		hpa.Target != 80 ||
		hpa.MinPods != 1 ||
		hpa.MaxPods != 5 ||
		hpa.Replicas != 3 ||
		hpa.Age != "33d" {
		t.Fatalf("Test failed! hpa does not match data")
	}

	hpa = hpas[1]
	if hpa.Namespace != "default" ||
		hpa.Name != "paymentservice" ||
		hpa.ReferenceKind != "Deployment" ||
		hpa.ReferenceName != "paymentservice" ||
		hpa.UsageCPU != 4 ||
		hpa.Target != 80 ||
		hpa.MinPods != 2 ||
		hpa.MaxPods != 20 ||
		hpa.Replicas != 2 ||
		hpa.Age != "87d" {
		t.Fatalf("Test failed! hpa does not match data")
	}
}
