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

var ipRegister = 1
var ipValue = 0

func main() {

	instructionRegex := regexp.MustCompile("([^ ]*)\\s(\\d*)\\s(\\d*)\\s(\\d*)\\r\\n?|\\n")
	
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
	

	register := [6]int{1,0,0,0,0,0}

	for ipValue < len(instructions) {
		
		register[ipRegister] = ipValue

		//fmt.Println("before:", register)
		//fmt.Println("instruction:", instructions[ipValue])

		switch instructions[ipValue].optcode {
		case "addr":
			register = addr(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "addi":
			register = addi(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "mulr":
			register = mulr(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "muli":
			register = muli(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "banr":
			register = banr(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "bani":
			register = bani(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "borr":
			register = borr(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "bori":
			register = bori(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "setr":
			register = setr(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "seti":
			register = seti(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "gtir":
			register = gtir(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "gtri":
			register = gtri(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "gtrr":
			register = gtrr(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "eqir":
			register = eqir(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "eqri":
			register = eqri(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		case "eqrr":
			register = eqrr(register, instructions[ipValue].inputA, instructions[ipValue].inputB, instructions[ipValue].outputC)
		}

		ipValue = register[ipRegister]
		ipValue++
		

		fmt.Println("after:", register)
		//fmt.Println("ipvalue", ipValue)
		//fmt.Println()
		//fmt.Println()


	}

	fmt.Println(register)

}


// addr (add register) stores into register C the result of adding register A and register B.
func addr(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = register[inputA] + register[inputB]
	return register	
} 


// addi (add immediate) stores into register C the result of adding register A and value B.
func addi(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = register[inputA] + inputB
	return register
} 


// mulr (multiply register) stores into register C the result of multiplying register A and register B.
func mulr(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = register[inputA] * register[inputB]
	return register
} 


// muli (multiply immediate) stores into register C the result of multiplying register A and value B.
func muli(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = register[inputA] * inputB
	return register
} 


// banr (bitwise AND register) stores into register C the result of the bitwise AND of register A and register B.
func banr(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = register[inputA] & register[inputB]
	return register
} 


// bani (bitwise AND immediate) stores into register C the result of the bitwise AND of register A and value B.
func bani(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = register[inputA] & inputB
	return register
} 


// borr (bitwise OR register) stores into register C the result of the bitwise OR of register A and register B.
func borr(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = register[inputA] | register[inputB]
	return register
} 


// bori (bitwise OR immediate) stores into register C the result of the bitwise OR of register A and value B.
func bori(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = register[inputA] | inputB
	return register
} 


// setr (set register) copies the contents of register A into register C. (Input B is ignored.)
func setr(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = register[inputA]
	return register
} 


// seti (set immediate) stores value A into register C. (Input B is ignored.)
func seti(register [6]int, inputA int, inputB int, outputC int) [6]int {
	register[outputC] = inputA
	return register
} 


// gtir (greater-than immediate/register) sets register C to 1 if value A is greater than register B. Otherwise, register C is set to 0.
func gtir(register [6]int, inputA int, inputB int, outputC int) [6]int {
	if inputA > register[inputB] { 
		register[outputC] = 1
	} else {
		register[outputC] = 0
	}
	return register
} 


// gtri (greater-than register/immediate) sets register C to 1 if register A is greater than value B. Otherwise, register C is set to 0.
func gtri(register [6]int, inputA int, inputB int, outputC int) [6]int {
	if register[inputA] > inputB { 
		register[outputC] = 1
	} else {
		register[outputC] = 0
	}
	return register
} 


// gtrr (greater-than register/register) sets register C to 1 if register A is greater than register B. Otherwise, register C is set to 0.
func gtrr(register [6]int, inputA int, inputB int, outputC int) [6]int {
	if register[inputA] > register[inputB] { 
		register[outputC] = 1
	} else {
		register[outputC] = 0
	}
	return register
} 


// eqir (equal immediate/register) sets register C to 1 if value A is equal to register B. Otherwise, register C is set to 0.
func eqir(register [6]int, inputA int, inputB int, outputC int) [6]int {
	if inputA == register[inputB] { 
		register[outputC] = 1
	} else {
		register[outputC] = 0
	}
	return register
} 


// eqri (equal register/immediate) sets register C to 1 if register A is equal to value B. Otherwise, register C is set to 0.
func eqri(register [6]int, inputA int, inputB int, outputC int) [6]int {
	if register[inputA] == inputB { 
		register[outputC] = 1
	} else {
		register[outputC] = 0
	}
	return register
} 


// eqrr (equal register/register) sets register C to 1 if register A is equal to register B. Otherwise, register C is set to 0.
func eqrr(register [6]int, inputA int, inputB int, outputC int) [6]int {
	if register[inputA] == register[inputB] { 
		register[outputC] = 1
	} else {
		register[outputC] = 0
	}
	return register
}