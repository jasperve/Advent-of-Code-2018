package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
)

type nanobot struct {
	x      int
	y      int
	z      int
	radius int
}

func main() {

	nanobots := []nanobot{}
	minX, maxX, minY, maxY, minZ, maxZ := 2147483647, 0, 2147483647, 0, 2147483647, 0

	input, _ := ioutil.ReadFile("input.txt")
	linesRegex := regexp.MustCompile("pos=<(-?\\d*),(-?\\d*),(-?\\d*)>,\\sr=(\\d*)(\\r\\n?|\\n)?")
	lines := linesRegex.FindAllStringSubmatch(string(input), -1)

	for _, line := range lines {
		x, _ := strconv.Atoi(line[1])
		y, _ := strconv.Atoi(line[2])
		z, _ := strconv.Atoi(line[3])
		radius, _ := strconv.Atoi(line[4])

		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
		if z < minZ {
			minZ = z
		}
		if z > maxZ {
			maxZ = z
		}

		newNanobot := nanobot{x: x, y: y, z: z, radius: radius}
		nanobots = append(nanobots, newNanobot)

	}

	devideBy, multiplyBy := 10000000, 10000000
	possibleX, possibleY, possibleZ := 0, 0, 0

	for {

		maxWithinRange := 0
		minDistance := 0

		for x := minX / devideBy; x <= maxX/devideBy; x++ {
			for y := minY / devideBy; y <= maxY/devideBy; y++ {
				for z := minZ / devideBy; z <= maxZ/devideBy; z++ {

					withinRange := 0
					for n := 0; n < len(nanobots); n++ {
						if int(math.Abs(float64(x-nanobots[n].x/devideBy)))+int(math.Abs(float64(y-nanobots[n].y/devideBy)))+int(math.Abs(float64(z-nanobots[n].z/devideBy))) <= nanobots[n].radius/devideBy {
							withinRange++
						}
					}

					if withinRange >= maxWithinRange || (withinRange == maxWithinRange && (int(math.Abs(float64(x-0)))+int(math.Abs(float64(y-0)))+int(math.Abs(float64(z-0))) < minDistance || minDistance == 0)) {
						maxWithinRange = withinRange
						minDistance = int(math.Abs(float64(x-0))) + int(math.Abs(float64(y-0))) + int(math.Abs(float64(z-0)))
						possibleX, possibleY, possibleZ = x, y, z
					}

				}
			}
		}

		minX = (possibleX - 3) * multiplyBy
		maxX = (possibleX + 3) * multiplyBy
		minY = (possibleY - 3) * multiplyBy
		maxY = (possibleY + 3) * multiplyBy
		minZ = (possibleZ - 3) * multiplyBy
		maxZ = (possibleZ + 3) * multiplyBy

		if devideBy == 1 {
			break
		}

		devideBy /= 10
		multiplyBy /= 10

	}

	fmt.Println(int(math.Abs(float64(possibleX-0))) + int(math.Abs(float64(possibleY-0))) + int(math.Abs(float64(possibleZ-0))))

}
