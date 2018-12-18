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
)

const (

	sand = 0
	clay = 1
	water = 2

	up = 0
	left = 1
	right = 2
	down = 3

)

var minY, maxY, minX, maxX int
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

	// Create a grid and prefill it with sand
	grid = make(map[int]map[int]int)
	for y := minY; y <= maxY; y++ {
		row := make(map[int]int)
		for x := minX; x <= maxX+2; x++ {
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
				grid[y][coordinateValue1+1] = clay
			}
		} else if line[1] == "y" {
			for x := coordinateValue2; x <= coordinateValue3; x++ {
				grid[coordinateValue1][x+1] = clay
			}
		}
	}

	// Start filling the grid with water
	fillWater(0, 500)
	displayFlow()

	/*for y := 201; y < 300; y++ {
		for x:= 450; x < 550; x++ {
			if y == 201 && x == 500 { 
				fmt.Printf("&")
				continue
			}
			if grid[y][x] == sand { fmt.Printf(".")}
			if grid[y][x] == clay { fmt.Printf("#")}
		}
		fmt.Printf("\n")
	}*/

}


// Returns true for obstruction found, false for clear path
func fillWater(y int, x int) {

	OUTER:
	for y <= maxY {

		// If the CURRENT SQUARE consists of SAND
		if grid[y][x] == sand  {
			fmt.Println("filling", y, ",", x, "with water")
			grid[y][x] = water
		} 

		// If the SQUARE BELOW consists of SAND
		if grid[y+1][x] == sand { 
			y++
		} else if grid[y+1][x] == clay || grid[y+1][x] == water {

			borderYLeft, borderYRight, border, borderXLeft, borderXRight := -1, -1, -1, -1, -1
			
			// FIND MAX WATER LEVEL LEFT
			FINDLOOPLEFT:
			for subY := y; subY >= minY; subY-- {
				for subX := x; subX >= minX; subX-- {
					if (grid[subY+1][subX] == sand && grid[subY+1][subX+1] == clay) {
						borderYLeft = subY
						borderXLeft = subX
						break FINDLOOPLEFT
					}
					if grid[subY][subX] == clay {
						break
					}
				}
			}

			// FIND MAX WATER LEVEL RIGHT
			FINDLOOPRIGHT:
			for subY := y; subY >= minY; subY-- {
				for subX := x; subX <= maxX; subX++ {
					if (grid[subY+1][subX] == sand && grid[subY+1][subX-1] == clay) {

						borderYRight = subY
						borderXRight = subX
						break FINDLOOPRIGHT
					}
					if grid[subY][subX] == clay {
						break
					}
				}
			}
			
			fmt.Println("border left", borderYLeft, "borderRight", borderYRight)
			if borderYLeft > borderYRight { 
				border = borderYLeft
			} else if borderYRight > borderYLeft { 
				border = borderYRight 
			} else if borderYRight == borderYLeft { 
				border = borderYRight 
			} 
			
			fmt.Println("fill to level", border)

			// FILL LEFT SIDE TO MAX LEVEL
			FILLLOOPLEFT:
			for subY := y; subY >= border; subY-- {
				for subX := x; subX >= minX; subX-- {
					if (grid[subY+1][subX] == sand && grid[subY+1][subX+1] == clay) {
						break FILLLOOPLEFT
					}

					fmt.Println("filling", subY, ",", subX, "with water")
					grid[subY][subX] = water
					
					if grid[subY][subX-1] == clay || grid[subY][subX-1] == water {
						break
					}
				}
			}

			// FILL RIGHT SIDE TO MAX LEVEL
			FILLLOOPRIGHT:
			for subY := y; subY >= border; subY-- {
				for subX := x; subX <= maxX; subX++ {
					if (grid[subY+1][subX] == sand && grid[subY+1][subX-1] == clay) {
						break FILLLOOPRIGHT
					}

					fmt.Println("filling", subY, ",", subX, "with water")
					grid[subY][subX] = water

					if grid[subY][subX+1] == clay || grid[subY][subX+1] ==  water {
						break
					}
				}
			}

			/*if borderYLeft == borderYRight {
				fillWater(border, borderXLeft)
				fillWater(border, borderXRight)
			} else*/ 
			if borderYLeft >= borderYRight {
				fillWater(border, borderXLeft)
			} else if borderYRight > borderYLeft {
				fillWater(border, borderXRight)
			}

			fmt.Println(borderXLeft, borderXRight)
			break OUTER

		}

	}

}

func displayFlow() {

	// Create a image based on the grid
	margin := 10
	img := image.NewRGBA(image.Rectangle{image.Point{minX, minY}, image.Point{maxX + 2 * margin, maxY + 2 * margin}})

	cyan := color.RGBA{100, 200, 200, 0xff}
	red := color.RGBA{255, 0, 0, 0xff}
	blue := color.RGBA{0, 0, 255, 0xff}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if grid[y][x] == clay { img.Set(x + margin, y + margin, cyan) }
			if grid[y][x] == water { img.Set(x + margin, y + margin, blue) }
		}
	}
	
	img.Set(510, 10, red)

	file, _ := os.Create("output.png")
	png.Encode(file, img)

	fmt.Println("file created")

}
