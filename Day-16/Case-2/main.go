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

type optcode struct {
	found bool
	number int
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

	register := [4]int{ 0, 0, 0, 0 }
	instructions := [][4]int{}

	inputInstructions, _ := ioutil.ReadFile("instructions.txt")
	for _, match := range strings.Split(string(inputInstructions), "\r\n"){

		instructionVal := instructionRegex.FindStringSubmatch(match)
		instruction := [4]int{}
		instruction[0], _ = strconv.Atoi(instructionVal[1])
		instruction[1], _ = strconv.Atoi(instructionVal[2])
		instruction[2], _ = strconv.Atoi(instructionVal[3])
		instruction[3], _ = strconv.Atoi(instructionVal[4])
		instructions = append(instructions, instruction)

	}

	numberCodes := []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	knownOpcodes := make(map[string]*optcode) 
	knownOpcodes["addr"] = &optcode { found: false, number: -1 }
	knownOpcodes["addi"] = &optcode { found: false, number: -1 }
	knownOpcodes["mulr"] = &optcode { found: false, number: -1 }
	knownOpcodes["muli"] = &optcode { found: false, number: -1 }
	knownOpcodes["banr"] = &optcode { found: false, number: -1 }
	knownOpcodes["bani"] = &optcode { found: false, number: -1 }
	knownOpcodes["borr"] = &optcode { found: false, number: -1 }
	knownOpcodes["bori"] = &optcode { found: false, number: -1 }
	knownOpcodes["setr"] = &optcode { found: false, number: -1 }
	knownOpcodes["seti"] = &optcode { found: false, number: -1 }	
	knownOpcodes["gtir"] = &optcode { found: false, number: -1 }
	knownOpcodes["gtri"] = &optcode { found: false, number: -1 }
	knownOpcodes["gtrr"] = &optcode { found: false, number: -1 }		
	knownOpcodes["eqir"] = &optcode { found: false, number: -1 }
	knownOpcodes["eqri"] = &optcode { found: false, number: -1 }		
	knownOpcodes["eqrr"] = &optcode { found: false, number: -1 }
	
	for len(numberCodes) > 0 {

		for n := 0; n < len(numberCodes); n++ {

			codesCounter := make(map[string]int)

			for s := 0; s < len(samples); s++ {

				if samples[s].instruction[0] != numberCodes[n] { continue }

				for k, v := range knownOpcodes {
					if v.found == false {
						if _, ok := codesCounter[k]; !ok {
							codesCounter[k] = 0
						}

						if k == "addr" && checkIfEqual(samples[s].registerAfter, addr(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "addi" && checkIfEqual(samples[s].registerAfter, addi(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "mulr" && checkIfEqual(samples[s].registerAfter, mulr(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "muli" && checkIfEqual(samples[s].registerAfter, muli(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "banr" && checkIfEqual(samples[s].registerAfter, banr(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "bani" && checkIfEqual(samples[s].registerAfter, bani(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "borr" && checkIfEqual(samples[s].registerAfter, borr(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "bori" && checkIfEqual(samples[s].registerAfter, bori(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "setr" && checkIfEqual(samples[s].registerAfter, setr(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "seti" && checkIfEqual(samples[s].registerAfter, seti(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "gtir" && checkIfEqual(samples[s].registerAfter, gtir(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "gtri" && checkIfEqual(samples[s].registerAfter, gtri(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "gtrr" && checkIfEqual(samples[s].registerAfter, gtrr(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "eqir" && checkIfEqual(samples[s].registerAfter, eqir(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "eqri" && checkIfEqual(samples[s].registerAfter, eqri(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
						if k == "eqrr" && checkIfEqual(samples[s].registerAfter, eqrr(samples[s].registerBefore, samples[s].instruction)) { codesCounter[k]++ }
					}
				}
			
			}

			matchingOpcodes := 0
			for _, v := range codesCounter {
				if v > 0 { matchingOpcodes++ }
			}

			if matchingOpcodes == 1 {
				for k, v := range codesCounter {
					if v > 0 { 
						knownOpcodes[k].found = true
						knownOpcodes[k].number = numberCodes[n]
						numberCodes = append(numberCodes[:n], numberCodes[n+1:]...)
						break
					}
				}
				break
			}

		}

	}

	for i := 0; i < len(instructions); i++ {

		opcodeName := ""

		for k, v := range knownOpcodes {
			if v.number == instructions[i][0] {
				opcodeName = k
				break
			}
		}

		switch opcodeName {
		case "eqri":
			register = eqri(register, instructions[i])
		case "mulr":
			register = mulr(register, instructions[i])
		case "gtri":
			register = gtri(register, instructions[i])
		case "gtrr":
			register = gtrr(register, instructions[i])
		case "banr":
			register = banr(register, instructions[i])
		case "addi":
			register = addi(register, instructions[i])
		case "seti":
			register = seti(register, instructions[i])
		case "gtir":
			register = gtir(register, instructions[i])
		case "muli":
			register = muli(register, instructions[i])
		case "bori":
			register = bori(register, instructions[i])
		case "setr":
			register = setr(register, instructions[i])
		case "addr":
			register = addr(register, instructions[i])
		case "bani":
			register = bani(register, instructions[i])
		case "borr":
			register = borr(register, instructions[i])
		case "eqir":
			register = eqir(register, instructions[i])
		case "eqrr":
			register = eqrr(register, instructions[i])
		}

	}

	fmt.Println(register)

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