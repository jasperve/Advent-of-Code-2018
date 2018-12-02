package main

import (
	"fmt"
	"strings"
	"io/ioutil"
)

func main() {

	var double, triple int

	input, _ := ioutil.ReadFile("input.txt")

	for _, id := range strings.Split(string(input), "\r\n"){

		idChars := make(map[string]int)

		for _, idChar := range id { idChars[string(idChar)]++	}

		doubleFound := false
		tripleFound := false

		for _, v := range idChars {
			if v == 2 {	doubleFound = true }
			if v == 3 {	tripleFound = true }
		}

		if doubleFound { double++ }
		if tripleFound { triple++ }

	}

	fmt.Printf("Checksum: %v", double*triple)

}