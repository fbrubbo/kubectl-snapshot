package main

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

func TestDateFormat(t *testing.T) {
	now := time.Now()
	fmt.Println(now.Format("kubectl-snapshot-2006-01-02-1504-pods.csv"))
}

func TestRegexp(t *testing.T) {
	match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	fmt.Println(match)
}

var data = [][]string{{"Line1", "Hello Readers of"}, {"Line2", "golangcode.com"}}

func TestCSV(t *testing.T) {
	// now := time.Now()
	// fileName := now.Format("kubectl-snapshot-2006-01-02-1504-pods.csv")

	// file, err := os.Create(fileName)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// writer := csv.NewWriter(file)
	// defer writer.Flush()

	// for _, value := range data {
	// 	err := writer.Write(value)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}
