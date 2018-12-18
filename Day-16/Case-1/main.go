package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"regexp"
	"strconv"
)

type sample struct {
	registerBefore [4]int
	instruction [4]int
	registerAfter [4]int
}


func main() {

	registerRegex := regexp.MustCompile("\\[(\\d*), (\\d*), (\\d*), (\\d*)\\]")
	instructionRegex := regexp.MustCompile("(\\d*)\\s*(\\d*)\\s*(\\d*)\\s*(\\d*)")
	
	samples := []sample{}

	inputSamples, _ := ioutil.ReadFile("samples.txt")
	for _, match := range strings.Split(string(inputSamples), "\r\n\r\n"){

		lines := strings.Split(match, "\r\n")

		sRegisterBeforeVal := registerRegex.FindStringSubmatch(lines[0])
		sRegisterBefore := [4]int{}
		sRegisterBefore[0], _ = strconv.Atoi(sRegisterBeforeVal[1])
		sRegisterBefore[1], _ = strconv.Atoi(sRegisterBeforeVal[2])
		sRegisterBefore[2], _ = strconv.Atoi(sRegisterBeforeVal[3])
		sRegisterBefore[3], _ = strconv.Atoi(sRegisterBeforeVal[4])

		sInstVal := instructionRegex.FindStringSubmatch(lines[1])
		sInstruction := [4]int{}
		sInstruction[0], _ = strconv.Atoi(sInstVal[1])
		sInstruction[1], _ = strconv.Atoi(sInstVal[2])
		sInstruction[2], _ = strconv.Atoi(sInstVal[3])
		sInstruction[3], _ = strconv.Atoi(sInstVal[4])

		sRegisterAfterVal := registerRegex.FindStringSubmatch(lines[2])
		sRegisterAfter := [4]int{}
		sRegisterAfter[0], _ = strconv.Atoi(sRegisterAfterVal[1])
		sRegisterAfter[1], _ = strconv.Atoi(sRegisterAfterVal[2])
		sRegisterAfter[2], _ = strconv.Atoi(sRegisterAfterVal[3])
		sRegisterAfter[3], _ = strconv.Atoi(sRegisterAfterVal[4])
		
		newSample := sample {
			registerBefore: sRegisterBefore,
			instruction: sInstruction,
			registerAfter: sRegisterAfter,
		}

		samples = append(samples, newSample)

	}

	totalCounter := 0
	for s := 0; s < len(samples); s++ {

		sampleCounter := 0

		if checkIfEqual(samples[s].registerAfter, addr(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, addi(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, mulr(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, muli(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, banr(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, bani(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, borr(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, bori(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, setr(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, seti(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, gtir(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, gtri(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, gtrr(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, eqir(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, eqri(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }
		if checkIfEqual(samples[s].registerAfter, eqrr(samples[s].registerBefore, samples[s].instruction)) { sampleCounter++ }

		if sampleCounter >= 3 { totalCounter++ }
	}

	fmt.Println(totalCounter)

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


// Check if 2 arrays contain the same values
func checkIfEqual(a, b [4]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}