package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {

	input, _ := ioutil.ReadFile("input.txt")
	inputSplit := []int{}

	for _, v := range strings.Split(string(input), " ") {
		num, _ := strconv.Atoi(v)
		inputSplit = append(inputSplit, num)
	}

	fmt.Println(calculateMetaData(inputSplit))

}

func calculateMetaData(input []int) (int, int) {

	numG, numMD, amountMD, readMD := input[0], input[1], 0, 2

	for i := 0; i < numG; i++ {
		subAmount, subRead := calculateMetaData(append([]int{}, input[readMD:]...))
		amountMD += subAmount
		readMD += subRead
	}

	for i := 0; i < numMD; i++ {
		amountMD += input[readMD]
		readMD++
	}

	return amountMD, readMD

}