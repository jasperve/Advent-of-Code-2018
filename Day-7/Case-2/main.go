package main

import (
	"bufio"
	"fmt"
	"os"
)

type elf struct {
	timeLeft int
	step uint8
}

const amountElves = 4

func main() {

	steps := make(map[uint8]map[uint8]bool)

	file, _ := os.Open("input.txt")
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

	output := []uint8{}
	time := 0

	elves := []elf{}
	for i := 0; i < amountElves; i++ {
		elves = append(elves, elf{timeLeft: 0, step: 0,})
	}
	activeElf := true

	for len(steps) > 0 || activeElf {

		activeElf = false

		for j := 0; j < len(elves); j++ {

			if elves[j].timeLeft == 0 {

				if elves[j].step != 0 {
					output = append(output, elves[j].step)
					elves[j].step = 0
				}

				for i := 'A'; i <= 'Z'; i++ {
					if step, ok := steps[uint8(i)]; ok {
						for _, v := range output {
							if _, ok := step[v]; ok { delete(step, v) }
						}
						if len(step) == 0 {
							elves[j].step = uint8(i)
							elves[j].timeLeft = int(i) - 4
							delete(steps, uint8(i))
							break
						}
					}
				}

			}

			if elves[j].timeLeft > 0 {
				activeElf = true
				elves[j].timeLeft--
			}

		}

		if activeElf { time++ }

	}

	fmt.Println(string(output))
	fmt.Println(time)

}