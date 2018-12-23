package main

import (
	"fmt"
)

const gridSize = 300
const gridSerialNumber = 7857

type coordinate struct {
	x    int
	y    int
	size int
}

func main() {

	grid := [][]int{}
	for x := 1; x < gridSize; x++ {
		row := []int{}
		for y := 1; y < gridSize; y++ {
			row = append(row, (y*(x+10)+gridSerialNumber)*(x+10)%1000/100-5)
		}
		grid = append(grid, row)
	}

	maxPower := coordinate{}
	maxPowerLevel := -2147483648

	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {

			var powerGridSizeLimit int
			powerGridSizeLimitX := gridSize - x
			powerGridSizeLimitY := gridSize - y

			if powerGridSizeLimitX <= powerGridSizeLimitY {
				powerGridSizeLimit = powerGridSizeLimitX
			} else {
				powerGridSizeLimit = powerGridSizeLimitY
			}

			powerLevel := 0

			for powerGridSize := 1; powerGridSize < powerGridSizeLimit; powerGridSize++ {

				if powerGridSize == 1 {
					powerLevel += grid[x][y]
					//fmt.Println("x: ", x, "y:", y, "amount:", grid[x][y])
				} else {
					for subX := x; subX < x+powerGridSize-1; subX++ {
						powerLevel += grid[subX][y+powerGridSize-1]
						//fmt.Println("x: ", subX, "y:", y+powerGridSize-1, "amount:", grid[subX][y+powerGridSize-1])
					}
					for subY := y; subY < y+powerGridSize; subY++ {
						powerLevel += grid[x+powerGridSize-1][subY]
						//fmt.Println("x: ", x+powerGridSize-1, "y:", subY, "amount:", grid[x+powerGridSize-1][subY])
					}
				}

				if powerLevel > maxPowerLevel {
					maxPower = coordinate{x: x + 1, y: y + 1, size: powerGridSize}
					maxPowerLevel = powerLevel
				}

			}

		}
	}

	fmt.Printf("Maximum power was found at: %v,%v,%v with a power of %v.", maxPower.x, maxPower.y, maxPower.size, maxPowerLevel)

}
