package main

import (
	"bufio"
	"fmt"
	"os"
)

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

	for len(steps) > 0 {
		for i := 65; i <= 90; i++ {
			if step, ok := steps[uint8(i)]; ok {
				for _, v := range output {
					if _, ok := step[v]; ok { delete(step, v) }
				}
				if len(step) == 0 {
					output = append(output, uint8(i))
					delete(steps, uint8(i))
					break
				}
			}
		}
	}

	fmt.Println(string(output))

}