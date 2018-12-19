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


// addr (add register) stores into register C the result of adding register A and register B.
func addr(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = register[instruction[1]] + register[instruction[2]]
	return register	
} 


// addi (add immediate) stores into register C the result of adding register A and value B.
func addi(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = register[instruction[1]] + instruction[2]
	return register

} 


// mulr (multiply register) stores into register C the result of multiplying register A and register B.
func mulr(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = register[instruction[1]] * register[instruction[2]]
	return register
} 


// muli (multiply immediate) stores into register C the result of multiplying register A and value B.
func muli(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = register[instruction[1]] * instruction[2]
	return register
} 


// banr (bitwise AND register) stores into register C the result of the bitwise AND of register A and register B.
func banr(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = register[instruction[1]] & register[instruction[2]]
	return register
} 


// bani (bitwise AND immediate) stores into register C the result of the bitwise AND of register A and value B.
func bani(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = register[instruction[1]] & instruction[2]
	return register
} 


// borr (bitwise OR register) stores into register C the result of the bitwise OR of register A and register B.
func borr(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = register[instruction[1]] | register[instruction[2]]
	return register
} 


// bori (bitwise OR immediate) stores into register C the result of the bitwise OR of register A and value B.
func bori(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = register[instruction[1]] | instruction[2]
	return register
} 


// setr (set register) copies the contents of register A into register C. (Input B is ignored.)
func setr(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = register[instruction[1]]
	return register
} 


// seti (set immediate) stores value A into register C. (Input B is ignored.)
func seti(register [4]int, instruction [4]int) (result [4]int) {
	register[instruction[3]] = instruction[1]
	return register
} 


// gtir (greater-than immediate/register) sets register C to 1 if value A is greater than register B. Otherwise, register C is set to 0.
func gtir(register [4]int, instruction [4]int) (result [4]int) {
	if instruction[1] > register[instruction[2]] { 
		register[instruction[3]] = 1
	} else {
		register[instruction[3]] = 0
	}
	return register
} 


// gtri (greater-than register/immediate) sets register C to 1 if register A is greater than value B. Otherwise, register C is set to 0.
func gtri(register [4]int, instruction [4]int) (result [4]int) {
	if register[instruction[1]] > instruction[2] { 
		register[instruction[3]] = 1
	} else {
		register[instruction[3]] = 0
	}
	return register
} 


// gtrr (greater-than register/register) sets register C to 1 if register A is greater than register B. Otherwise, register C is set to 0.
func gtrr(register [4]int, instruction [4]int) (result [4]int) {
	if register[instruction[1]] > register[instruction[2]] { 
		register[instruction[3]] = 1
	} else {
		register[instruction[3]] = 0
	}
	return register
} 


// eqir (equal immediate/register) sets register C to 1 if value A is equal to register B. Otherwise, register C is set to 0.
func eqir(register [4]int, instruction [4]int) (result [4]int) {
	if instruction[1] == register[instruction[2]] { 
		register[instruction[3]] = 1
	} else {
		register[instruction[3]] = 0
	}
	return register
} 


// eqri (equal register/immediate) sets register C to 1 if register A is equal to value B. Otherwise, register C is set to 0.
func eqri(register [4]int, instruction [4]int) (result [4]int) {
	if register[instruction[1]] == instruction[2] { 
		register[instruction[3]] = 1
	} else {
		register[instruction[3]] = 0
	}
	return register
} 


// eqrr (equal register/register) sets register C to 1 if register A is equal to register B. Otherwise, register C is set to 0.
func eqrr(register [4]int, instruction [4]int) (result [4]int) {
	if register[instruction[1]] == register[instruction[2]] { 
		register[instruction[3]] = 1
	} else {
		register[instruction[3]] = 0
	}
	return register
}