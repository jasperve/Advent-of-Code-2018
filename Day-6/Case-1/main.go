package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type coordinate struct {
	x        int
	y        int
	infinite bool
	count    int
}

func main() {

	file, _ := os.Open("input.txt")

	coordinates := []coordinate{}

	input := bufio.NewScanner(file)
	for input.Scan() {

		line := strings.Split(input.Text(), ", ")
		x, _ := strconv.Atoi(line[0])
		y, _ := strconv.Atoi(line[1])

		c := coordinate{x: x, y: y}
		coordinates = append(coordinates, c)

	}

	maxX, maxY := 0, 0

	for _, v := range coordinates {
		if v.x > maxX {
			maxX = v.x
		}
		if v.y > maxY {
			maxY = v.y
		}
	}

	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {

			nearestIndex := -1

			var minSteps = maxX + maxY
			for i, v := range coordinates {

				steps := int(math.Abs(float64(v.x-x))) + int(math.Abs(float64(v.y-y)))
				if steps < minSteps {
					minSteps = steps
					nearestIndex = i
				} else if steps == minSteps {
					nearestIndex = -1
				}

			}

			if (y == 0 || y == maxY || x == 0 || x == maxX) && nearestIndex != -1 {
				coordinates[nearestIndex].infinite = true
			} else if nearestIndex != -1 {
				coordinates[nearestIndex].count++
			}

		}
	}

	biggestArea := -1
	countbiggestArea := -1

	for i, v := range coordinates {
		fmt.Println(v.x, v.y, v.count, v.infinite)
		if !v.infinite && v.count > countbiggestArea {
			biggestArea = i
		}
	}

	fmt.Printf("Biggest area is around coordinates: %v, %v with a size of %v\n", coordinates[biggestArea].x, coordinates[biggestArea].y, coordinates[biggestArea].count)

}
