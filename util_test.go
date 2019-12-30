package main

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

type test struct {
	In  string
	Out int
}

func TestString2MilliCPU(t *testing.T) {

	tests := []test{
		test{In: "130m", Out: 130},
		test{In: "1", Out: 1000},
		test{In: "0.5", Out: 500},
		test{In: "1.64", Out: 1640},
	}

	for i, tes := range tests {
		log.Infof("test info %d -> %+v", i, tes)
		if result := String2MilliCPU(tes.In); result != tes.Out {
			t.Fatalf("Test failed! %d but expected %d", result, tes.Out)
		}
	}
}

func TestString2MiMemory(t *testing.T) {

	tests := []test{
		test{In: "123Mi", Out: 123},
		test{In: "129M", Out: 123},
		test{In: "128974848", Out: 123},
		test{In: "125952Ki", Out: 123},
	}

	for i, tes := range tests {
		log.Infof("test info %d -> %+v", i, tes)
		if result := String2MiMemory(tes.In); result != tes.Out {
			t.Fatalf("Test failed! %d but expected %d", result, tes.Out)
		}
	}
}
