package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

const (
	sand  = 0
	clay  = 1
	water = 2

	up    = 0
	left  = 1
	right = 2
	down  = 3
)

type coordinate struct {
	x int
	y int
}

var minY, maxY, minX, maxX = 2147483647, 0, 2147483647, 0
var grid map[coordinate]int

func main() {

	input, _ := ioutil.ReadFile("input-jasper.txt")
	linesRegex := regexp.MustCompile("(x|y)=(\\d*),\\s(x|y)=(\\d*)\\.*(\\d*)")
	lines := linesRegex.FindAllStringSubmatch(string(input), -1)

	// Create a grid and prefill it with sand
	grid = make(map[coordinate]int)

	// Fill the grid with the clay positions
	for _, line := range lines {
		coordinateValue1, _ := strconv.Atoi(line[2])
		coordinateValue2, _ := strconv.Atoi(line[4])
		coordinateValue3, _ := strconv.Atoi(line[5])
		if line[1] == "x" {
			for y := coordinateValue2; y <= coordinateValue3; y++ {
				grid[coordinate{coordinateValue1, y}] = clay
			}
		} else if line[1] == "y" {
			for x := coordinateValue2; x <= coordinateValue3; x++ {
				grid[coordinate{x, coordinateValue1}] = clay
			}
		}
	}

	// Determine the grid size
	for k := range grid {
		if k.x < minX {
			minX = k.x
		}
		if k.x > maxX {
			maxX = k.x
		}
		if k.y < minY {
			minY = k.y
		}
		if k.y > maxY {
			maxY = k.y
		}
	}

	// Increase the grid a bit to the left and right to allow overflow of the buckets on the edges
	minX -= 5
	maxX += 5

	// Fill the empty spots in the grid with sand
	for x := minX; x <= maxX; x++ {
		for y := minY; y <=maxY; y++ {
			if _, ok := grid[coordinate{x, y}]; !ok {
				grid[coordinate{x, y}] = sand
			}
		}
	}

	fmt.Println(minX, maxX, minY, maxY)

	// Start filling the grid with water
	fillWater(500, minY)

	// Count the amount of water while flowing
	whileFlowingCounter := 0
	for _, g := range grid {
		if g == water {
			whileFlowingCounter++
		}
	}

	fmt.Println("While the well is flowing:", whileFlowingCounter)
	displayGrid("WhileFlowing.png")

	// Remove all standing water
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if grid[coordinate{x, y}] == water && (x == minX || x == maxX) {
				grid[coordinate{x, y}] = sand
			}
			if grid[coordinate{x, y}] == water && (grid[coordinate{x-1, y}] == sand || grid[coordinate{x+1, y}] == sand) {
				grid[coordinate{x, y}] = sand
			}
		}
		for x := maxX; x >= minX; x-- {
			if grid[coordinate{x, y}] == water && (x == minX || x == maxX) {
				grid[coordinate{x, y}] = sand
			}
			if grid[coordinate{x, y}] == water && (grid[coordinate{x-1, y}] == sand || grid[coordinate{x+1, y}] == sand) {
				grid[coordinate{x, y}] = sand
			}
		}
	}

	afterFlowingCounter := 0
	for _, g := range grid {
		if g == water {
			afterFlowingCounter++
		}
	}

	fmt.Println("After the well has stopped flowing:", afterFlowingCounter)
	displayGrid("AfterFlowing.png")

}


func fillWater(x int, y int) {

	for y <= maxY {

		// If the CURRENT SQUARE consists of SAND
		if grid[coordinate{x, y}] == sand {
			grid[coordinate{x, y}] = water
		}

		// If the SQUARE BELOW consists of SAND or WATER
		if grid[coordinate{x, y + 1}] == sand || grid[coordinate{x, y + 1}] == water {
			y++
		} else if grid[coordinate{x, y + 1}] == clay {

			for subY := y; subY >= minY; subY-- {

				edges := []coordinate{}

				// Look left
				for subX := x; subX >= minX; subX-- {

					grid[coordinate{subX, subY}] = water

					if subX != x &&
						grid[coordinate{subX, subY}] != clay &&
						grid[coordinate{subX, subY + 1}] != clay &&
						grid[coordinate{subX + 1, subY}] != clay &&
						grid[coordinate{subX - 1, subY}] == sand &&
						grid[coordinate{subX + 1, subY + 1}] == clay {

						edges = append(edges, coordinate{subX, subY + 1})
						break

					}

					if grid[coordinate{subX - 1, subY}] == clay {
						break
					}

				}

				// Look right
				for subX := x; subX < maxX; subX++ {

					grid[coordinate{subX, subY}] = water

					if subX != x &&
						grid[coordinate{subX, subY}] != clay &&
						grid[coordinate{subX, subY + 1}] != clay &&
						grid[coordinate{subX - 1, subY}] != clay &&
						grid[coordinate{subX + 1, subY}] == sand &&
						grid[coordinate{subX - 1, subY + 1}] == clay {

						edges = append(edges, coordinate{subX, subY + 1})
						break
					}

					if grid[coordinate{subX + 1, subY}] == clay {
						break
					}

				}

				if len(edges) > 0 {
					for _, edge := range edges {
						if grid[coordinate{edge.x, edge.y + 1}] != water {
							fillWater(edge.x, edge.y)
						}
					}
					return
				}

			}

		}

	}

}


func displayGrid(location string) {

	img := image.NewRGBA(image.Rectangle{image.Point{minX - 2, minY - 2}, image.Point{maxX + 2, maxY + 2}})

	cyan := color.RGBA{100, 200, 200, 0xff}
	red := color.RGBA{255, 0, 0, 0xff}
	blue := color.RGBA{0, 0, 255, 0xff}

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if grid[coordinate{x, y}] == clay {
				img.Set(x, y, cyan)
			} 
			if grid[coordinate{x, y}] == water {
				img.Set(x, y, blue)
			}
		}
	}

	img.Set(500, 0, red)

	file, _ := os.Create(location)
	png.Encode(file, img)

}
