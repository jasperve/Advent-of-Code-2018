package main

import (
	"fmt"
	"math"
)

const gridSize = 300
const subGridSize = 6 // Appears to be optimal
const gridSerialNumber = 7857

type coordinate struct {
	x    int
	y    int
	size int
}

func main() {

	grid := [][]int{}
	for x := 1; x <= gridSize; x++ {
		row := []int{}
		for y := 1; y <= gridSize; y++ {
			row = append(row, (y*(x+10)+gridSerialNumber)*(x+10)%1000/100-5)
		}
		grid = append(grid, row)
	}

	subGrid := [][]int{}
	for x := 0; x < gridSize/subGridSize; x++ {
		subRow := []int{}
		for y := 0; y < gridSize/subGridSize; y++ {
			total := 0
			for subX := x * subGridSize; subX < (x+1)*subGridSize; subX++ {
				for subY := y * subGridSize; subY < (y+1)*subGridSize; subY++ {
					total += grid[subX][subY]
				}
			}
			subRow = append(subRow, total)
		}
		subGrid = append(subGrid, subRow)
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

			for powerGridSize := 1; powerGridSize <= powerGridSizeLimit; powerGridSize++ {

				powerLevel := 0
				fromX := float64(x)
				fromY := float64(y)
				tillX := x + powerGridSize - 1
				tillY := y + powerGridSize - 1

				firstSubGridX := int(math.Ceil(fromX / subGridSize))
				firstSubGridY := int(math.Ceil(fromY / subGridSize))
				lastSubGridX := tillX / subGridSize
				lastSubGridY := tillY / subGridSize

				//Calculate all the subGrids that can be used
				for iX := firstSubGridX; iX < lastSubGridX; iX++ {
					for iY := firstSubGridY; iY < lastSubGridY; iY++ {
						powerLevel += subGrid[iX][iY]
					}
				}

				//Calculates everything left of the subgrid (including above and below)
				for iX := x; iX < firstSubGridX*subGridSize && iX <= tillX; iX++ {
					for iY := y; iY <= tillY; iY++ {
						powerLevel += grid[iX][iY]
					}
				}

				//Calculates everything right of the subgrid (including above and below)
				for iX := lastSubGridX * subGridSize; iX <= tillX && lastSubGridX > 0; iX++ {
					for iY := y; iY <= tillY; iY++ {
						powerLevel += grid[iX][iY]
					}
				}

				//Calculates everything above the subgrid
				for iX := firstSubGridX * subGridSize; iX < lastSubGridX*subGridSize; iX++ {
					for iY := y; iY < firstSubGridY*subGridSize; iY++ {
						powerLevel += grid[iX][iY]
					}
				}

				//Calculates everything below the subgrid
				for iX := firstSubGridX * subGridSize; iX < lastSubGridX*subGridSize; iX++ {
					for iY := lastSubGridY * subGridSize; iY <= tillY; iY++ {
						powerLevel += grid[iX][iY]
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
