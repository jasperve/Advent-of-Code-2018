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
	maxRadius := 0
	var maxRadiusNanobot nanobot

	input, _ := ioutil.ReadFile("input.txt")
	linesRegex := regexp.MustCompile("pos=<(-?\\d*),(-?\\d*),(-?\\d*)>,\\sr=(\\d*)(\\r\\n?|\\n)?")
	lines := linesRegex.FindAllStringSubmatch(string(input), -1)

	for _, line := range lines {

		x, _ := strconv.Atoi(line[1])
		y, _ := strconv.Atoi(line[2])
		z, _ := strconv.Atoi(line[3])
		radius, _ := strconv.Atoi(line[4])

		newNanobot := nanobot{x: x, y: y, z: z, radius: radius}
		nanobots = append(nanobots, newNanobot)

		if radius > maxRadius {
			maxRadius = radius
			maxRadiusNanobot = newNanobot
		}

	}

	nanobotsInRange := 0

	for n := 0; n < len(nanobots); n++ {

		distance := int(math.Abs(float64(nanobots[n].x-maxRadiusNanobot.x))) + int(math.Abs(float64(nanobots[n].y-maxRadiusNanobot.y))) + int(math.Abs(float64(nanobots[n].z-maxRadiusNanobot.z)))
		if distance <= maxRadiusNanobot.radius {
			nanobotsInRange++
		}

	}

	fmt.Println(nanobotsInRange)

}
