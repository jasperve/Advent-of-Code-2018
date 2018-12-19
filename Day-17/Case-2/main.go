package main

import (
	"io/ioutil"
	"fmt"
	"regexp"
	"strconv"
	"image"
	"image/color"
	"image/png"
	"os"
	//"time"
)

const (

	sand = 0
	clay = 1
	water = 2
	temp = 3

	up = 0
	left = 1
	right = 2
	down = 3

)

var minY, maxY, minX, maxX = 10000, 0, 0, 0
var grid map[int]map[int]int

func main() {

	input, _ := ioutil.ReadFile("input.txt")
	linesRegex := regexp.MustCompile("(x|y)=(\\d*),\\s(x|y)=(\\d*)\\.*(\\d*)")
	lines := linesRegex.FindAllStringSubmatch(string(input), -1)

	// Calculate the grid dimensions based on the input
	for _, line := range lines {
		coordinateValue1, _ := strconv.Atoi(line[2])
		coordinateValue2, _ := strconv.Atoi(line[4])
		coordinateValue3, _ := strconv.Atoi(line[5])
		if line[1] == "x" {
			if coordinateValue1 < minX || minX == 0 { minX = coordinateValue1 }
			if coordinateValue1 > maxX || maxX == 0 { maxX = coordinateValue1 }
			if coordinateValue2 < minY { minY = coordinateValue2 }
			if coordinateValue3 < minY { minY = coordinateValue3 }
		} else if line[1] == "y" {
			if coordinateValue1 < minY { minY = coordinateValue1 }
			if coordinateValue1 > maxY || maxY == 0 { maxY = coordinateValue1 }
			if coordinateValue2 < minX || minX == 0 { minX = coordinateValue2 }
			if coordinateValue3 < minX || minX == 0 { maxX = coordinateValue3 }
		}
	}

	minX = minX -5
	maxX = maxX +5

	// Create a grid and prefill it with sand
	grid = make(map[int]map[int]int)
	for y := minY; y <= maxY; y++ {
		row := make(map[int]int)
		for x := minX; x <= maxX; x++ {
			row[x] = sand
		}
		grid[y] = row
	}

	// Fill the grid with the clay positions
	for _, line := range lines {
		coordinateValue1, _ := strconv.Atoi(line[2])
		coordinateValue2, _ := strconv.Atoi(line[4])
		coordinateValue3, _ := strconv.Atoi(line[5])
		if line[1] == "x" {
			for y := coordinateValue2; y <= coordinateValue3; y++ {
				grid[y][coordinateValue1] = clay
			}
		} else if line[1] == "y" {
			for x := coordinateValue2; x <= coordinateValue3; x++ {
				grid[coordinateValue1][x] = clay
			}
		}
	}

	// Start filling the grid with water
	fillWater(minY, 500)
	displayFlow()

	counter := 0

	for y := minY; y <= maxY; y++ {
		for x:= minX; x < maxX; x++ {
			if grid[y][x] == water { counter++ }
		}
	}

	fmt.Println(counter)

}


// Returns true for obstruction found, false for clear path
func fillWater(y int, x int) {

	//grid[y][x] = water

	//maxY = 50
	OUTER:
	for y <= maxY {

		// If the CURRENT SQUARE consists of SAND
		//if grid[y][x] == sand  {
			//grid[y][x] = waterf
		//} 

		// If the SQUARE BELOW consists of SAND
		if grid[y+1][x] == sand || grid[y+1][x] == water { 
			y++
		} else if grid[y+1][x] == clay {

			borderYLeft, borderYRight, border, borderXLeft, borderXRight := -1, -1, -1, -1, -1
			
			// FIND MAX WATER LEVEL LEFT
			FINDLOOPLEFT:
			for subY := y; subY >= minY; subY-- {
				for subX := x-1; subX >= minX; subX-- {
					if grid[subY][subX] != clay && grid[subY+1][subX] != clay && grid[subY][subX+1] != clay && grid[subY][subX-1] != clay && grid[subY+1][subX+1] == clay  {
						borderYLeft = subY
						borderXLeft = subX
						break FINDLOOPLEFT
					} else if grid[subY][subX] == clay {
						break
					}
				}
			}

			// FIND MAX WATER LEVEL RIGHT
			FINDLOOPRIGHT:
			for subY := y; subY >= minY; subY-- {
				for subX := x+1; subX <= maxX; subX++ {
					if grid[subY][subX] != clay && grid[subY+1][subX] != clay && grid[subY][subX-1] != clay && grid[subY][subX+1] != clay && grid[subY+1][subX-1] == clay  {
						borderYRight = subY
						borderXRight = subX
						break FINDLOOPRIGHT
					} else if grid[subY][subX] == clay {
						break
					}
				}
			}
			
			if borderYLeft > borderYRight { 
				border = borderYLeft
			} else if borderYRight > borderYLeft { 
				border = borderYRight 
			} else if borderYRight == borderYLeft { 
				border = borderYRight 
			} 
						
			// FILL LEFT SIDE TO MAX LEVEL
			for subY := y; subY >= border+1; subY-- {
				for subX := x; subX >= minX; subX-- {
					grid[subY][subX] = water
					if grid[subY][subX-1] == clay {	break }
										
				}
			}

			// FILL RIGHT SIDE TO MAX LEVEL
			for subY := y; subY >= border+1; subY-- {
				for subX := x; subX < maxX; subX++ {
					grid[subY][subX] = water
					if grid[subY][subX+1] == clay {	break }
				}
			}

			if borderYLeft == borderYRight {
				if grid[border][borderXLeft] == sand { grid[border][borderXLeft] = temp }
				if grid[border][borderXRight] == sand { grid[border][borderXLeft] = temp }
			}

			if borderYLeft == borderYRight {
				if grid[border][borderXLeft] == temp { 
					grid[border][borderXLeft] = water
					fillWater(border, borderXLeft) 
				}
				if grid[border][borderXRight] == temp { 
					grid[border][borderXRight] = water
					fillWater(border, borderXRight) 
				}
			} else if borderYLeft >= borderYRight {
				//grid[border][borderXLeft] = water
				fillWater(border, borderXLeft) 
			} else if borderYRight > borderYLeft {
				//grid[border][borderXRight] = water
				fillWater(border, borderXRight)
			}
			
			break OUTER

		}

	}

}

func displayFlow() {

	// Create a image based on the grid
	img := image.NewRGBA(image.Rectangle{image.Point{minX-2, minY-2}, image.Point{maxX+2, maxY+2}})

	cyan := color.RGBA{100, 200, 200, 0xff}
	red := color.RGBA{255, 0, 0, 0xff}
	blue := color.RGBA{0, 0, 255, 0xff}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if grid[y][x] == clay { img.Set(x, y, cyan) }
			if grid[y][x] == water { img.Set(x, y, blue) }
		}
	}
	
	img.Set(500, 0, red)

	file, _ := os.Create("output.png")
	png.Encode(file, img)

	fmt.Println("file created")

}
