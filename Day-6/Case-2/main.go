package main

import (

	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"math"

)

type coordinate struct {
	x int
	y int
}

func main() {

	file, _ := os.Open("input.txt")

	coordinates := []coordinate{}

	input := bufio.NewScanner(file)
	for input.Scan() {

		line := strings.Split(input.Text(), ", ")
		x, _ := strconv.Atoi(line[0])
		y, _ := strconv.Atoi(line[1])

		c := coordinate { x: x, y: y, }
		coordinates = append(coordinates, c)

	}

	maxX, maxY := 0, 0

	for _, v := range coordinates {
		if v.x > maxX { maxX = v.x }
		if v.y > maxY { maxY = v.y }
	}

	numInRange := 0

	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {

			distance := 0
			for _, v := range coordinates {
				distance += int(math.Abs(float64(v.x-x))) + int(math.Abs(float64(v.y-y)))
			}

			if distance < 10000 {
				numInRange++
			}
		}
	}

	fmt.Printf("The number of locations which have a range of < 10000 to each coordinate is %v\n", numInRange)

}