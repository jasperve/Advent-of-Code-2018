package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	input := splitToIntSlice("input.txt", " ")
	fmt.Println(calculateMd(input))
}

func splitToIntSlice(location string, sep string) (out []int) {
	input, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatalln("FATAL: Unable to open file at location: %v", location)
	}
	for _, token := range strings.Split(string(input), sep) {
		value, err := strconv.Atoi(token)
		if err != nil {
			log.Fatalln("FATAL: Unable to convert %v", token)
		}
		out = append(out, value)
	}
	return out
}

func calculateMd(input []int) (int, int) {

	numG, numMd, amountMd, amountMdGroup, readMd := input[0], input[1], 0, []int{}, 2

	for i := 0; i < numG; i++ {
		subAmount, subRead := calculateMd(append([]int{}, input[readMd:]...))
		amountMdGroup = append(amountMdGroup, subAmount)
		readMd += subRead
	}

	for i := 0; i < numMd; i++ {
		if numG == 0 {
			amountMd += input[readMd]
		} else {
			if input[readMd] <= len(amountMdGroup) {
				amountMd += amountMdGroup[input[readMd]-1]
			}
		}
		readMd++
	}

	return amountMd, readMd

}
