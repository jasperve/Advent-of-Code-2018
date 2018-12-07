package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type elf struct {
	active bool
	step uint8
}

const amountElves = 2

var comm chan int
var elves []elf
var steps map[uint8]map[uint8]bool
var output []uint8

func main() {

	comm = make(chan int)
	go listen()

	steps = make(map[uint8]map[uint8]bool)

	file, _ := os.Open("input-test.txt")
	input := bufio.NewScanner(file)
	for input.Scan() {
		line := input.Text()
		if _, ok := steps[line[36]]; !ok {
			steps[line[36]] = make(map[uint8]bool)
			steps[line[36]][line[5]] = true
		} else {
			steps[line[36]][line[5]] = true
		}
		if _, ok := steps[line[5]]; !ok {
			steps[line[5]] = make(map[uint8]bool)
		}
	}

	output = []uint8{}

	for i := 0; i < amountElves; i++ {
		elves = append(elves, elf{active: false, step: 0,})
	}

	start := time.Now()

	for len(steps) > 0 {

		for j := 0; j < len(elves); j++ {

			if !elves[j].active {

				for i := 'A'; i <= 'Z'; i++ {
					if step, ok := steps[uint8(i)]; ok {
						for _, v := range output {
							if _, ok := step[v]; ok { delete(step, v) }
						}
						if len(step) == 0 {

							go startWork(j, int(i)-64 ,comm)
							elves[j].step = uint8(i)
							elves[j].active = true
							delete(steps, uint8(i))
							break
						}
					}
				}

			}

		}

	}

	activeElf := true

	for activeElf {
		activeElf = false
		for _, v := range elves {
			if v.active {
				activeElf = true
			}
		}
	}

	end := time.Now()


	fmt.Println(string(output))
	fmt.Println(end.Sub(start))

}

func listen() {
	for number := range comm {
		elves[number].active = false
		output = append(output, elves[number].step)
		elves[number].step = 0

	}
}

func startWork(elf int, taskLength int, finished chan int) {
	time.Sleep(time.Duration(taskLength) * time.Second)
	finished <- elf
}