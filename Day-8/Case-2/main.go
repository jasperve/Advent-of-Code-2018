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

	fmt.Println(calculateMd(inputSplit))

}

func calculateMd(input []int) (int, int) {

	numG, numMd, amountMd, amountMdGroup, readMd := input[0], input[1], 0, []int{} , 2

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