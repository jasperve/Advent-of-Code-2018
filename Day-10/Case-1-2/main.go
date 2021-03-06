package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type star struct {
	x    int
	y    int
	velX int
	velY int
}

func main() {

	stars := []star{}

	file, _ := os.Open("input.txt")
	input := bufio.NewScanner(file)

	for input.Scan() {

		regex, _ := regexp.Compile("position=<.*?(\\-?\\d*),.*?(\\-?\\d*)> velocity=<.*?(\\-?\\d*),.*?(\\-?\\d*)>")
		matches := regex.FindStringSubmatch(input.Text())

		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])
		velX, _ := strconv.Atoi(matches[3])
		velY, _ := strconv.Atoi(matches[4])

		newStar := star{x: x, y: y, velX: velX, velY: velY}
		stars = append(stars, newStar)

	}

	var locations map[string]bool
	iteration := 0

MAIN:
	for {

		minX, maxX, minY, maxY := 2147483647, -2147483647, 2147483647, -2147483647

		locations = make(map[string]bool)

		for _, star := range stars {

			newX := star.x + star.velX*iteration
			newY := star.y + star.velY*iteration
			neighbour := false

			for x := newX - 1; x <= newX+1; x++ {
				for y := newY - 1; y <= newY+1; y++ {

					if _, ok := locations[fmt.Sprintf("%v,%v", x, y)]; ok {
						locations[fmt.Sprintf("%v,%v", x, y)] = true
						neighbour = true
					}

				}
			}

			locations[fmt.Sprintf("%v,%v", newX, newY)] = neighbour

			if newX < minX {
				minX = newX
			}
			if newX > maxX {
				maxX = newX
			}
			if newY < minY {
				minY = newY
			}
			if newY > maxY {
				maxY = newY
			}

		}

		for _, neighbour := range locations {
			if neighbour == false {
				iteration++
				continue MAIN
			}
		}

		modX := minX * -1
		modY := minY * -1

		img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{maxX + modX + 10, maxY + modY + 10}})

		cyan := color.RGBA{100, 200, 200, 0xff}

		for coordinate, _ := range locations {
			coordinateSplit := strings.Split(coordinate, ",")
			x, _ := strconv.Atoi(coordinateSplit[0])
			y, _ := strconv.Atoi(coordinateSplit[1])
			img.Set(x+modX+5, y+modY+5, cyan)
		}

		f, _ := os.Create("output.png")
		png.Encode(f, img)

		fmt.Printf("After %v seconds the stars line up!", iteration)

		break

	}

}
