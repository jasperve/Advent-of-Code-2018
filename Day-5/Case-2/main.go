package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var remBytes map[int]int
var input []byte


func main() {

	alphabet := "abcdefghijklmnopqrstuvwxyz"

	for j:=0; j < len(alphabet); j++ {

		input, _ = ioutil.ReadFile("input.txt")
		input = []byte(strings.Replace(string(input), string(alphabet[j]), "", -1))
		input = []byte(strings.Replace(string(input), string(alphabet[j]-32), "", -1))

		remBytes = make(map[int]int)
		for i := 0; i < len(input)-1; i++ {
			i += checkBytes(i, i+1) / 2
		}
		for i := len(remBytes) - 1; i >= 0; i-- {
			input = append(input[:remBytes[i]], input[remBytes[i]+1:]...)
		}

		fmt.Printf("Letter %v gives this many chars left: %v\n", string(alphabet[j]), len(input))

	}

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
