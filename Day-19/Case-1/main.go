package main

import (
	"fmt"
	"strconv"
	"regexp"
	"io/ioutil"
)

type instruction struct {
	optcode string
	inputA int
	inputB int
	outputC int
}

var instructionPointer = 1

func main() {

	instructionRegex := regexp.MustCompile("([^ ]*)\\s(\\d*)\\s(\\d*)\\s(\\d*)")
	
	instructions := []instruction{}

	input, _ := ioutil.ReadFile("input.txt")
	lines := instructionRegex.FindAllStringSubmatch(string(input), -1)
	
	for _, line := range lines { 
		optcode := string(line[1])
		inputA, _ := strconv.Atoi(string(line[2]))
		inputB, _ := strconv.Atoi(string(line[3]))
		outputC, _ := strconv.Atoi(string(line[4]))
		
		newInstruction := instruction {
			optcode: optcode,
			inputA: inputA,
			inputB: inputB,
			outputC: outputC,
		}

		instructions = append(instructions, newInstruction)
	}
	
	fmt.Println(instructions[0].optcode)


}