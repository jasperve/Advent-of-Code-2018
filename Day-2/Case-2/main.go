package main

import (
	"fmt"
	"strings"
	"io/ioutil"
)

func main() {

	input, _ := ioutil.ReadFile("input.txt")

	ids := strings.Split(string(input), "\r\n")

	for _, id := range ids {

		for _, compId := range ids {

			numDiffChar := 0
			diffCharPosition := 0

			for i:=0; i < len(id); i++ {
				if id[i] != compId[i] {
					numDiffChar++
					diffCharPosition = i
				}
			}

			if numDiffChar == 1 {
				fmt.Printf("The correct id is: %v%v\n", id[0:diffCharPosition], id[diffCharPosition+1:len(id)])
				return
			}

		}

	}

}