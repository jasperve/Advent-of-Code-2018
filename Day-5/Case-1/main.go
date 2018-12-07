package main

import (
	"fmt"
	"io/ioutil"
)

var remBytes map[int]int
var input []byte

func main() {

	remBytes = make(map[int]int)
	input, _ = ioutil.ReadFile("input.txt")

	for i := 0; i < len(input)-1; i++ {
		i += checkBytes(i, i+1) / 2
	}

	for i := len(remBytes) - 1; i >= 0; i-- {
		input = append(input[:remBytes[i]], input[remBytes[i]+1:]...)
	}

	fmt.Println(len(input))

}

func checkBytes(i int, j int) int {

	numRemBytes := 0
	if input[i] == input[j]-32 || input[i] == input[j]+32 {
		remBytes[i] = i
		remBytes[j] = j

		i -= 1
		j += 1

		for {
			if _, ok := remBytes[i]; ok {
				i--
			} else {
				break
			}
		}

		for {
			if _, ok := remBytes[j]; ok {
				j++
			} else {
				break
			}
		}

		if i >= 0 && j < len(input) {
			numRemBytes = checkBytes(i, j)
		}

		return 2 + numRemBytes
	}
	return numRemBytes

}