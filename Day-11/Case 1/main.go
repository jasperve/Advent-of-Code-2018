package main

import "fmt"

const gridSize = 300
const gridSerialNumber = 7857

type coordinate struct {
	x int
	y int
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

	maxPower := coordinate{}
	maxPowerLevel := -2147483648

	for x := 1; x < gridSize-2; x++ {
		for y := 1; y < gridSize-2; y++ {

			powerLevel := 0
			for subX := x - 1; subX <= x+1; subX++ {
				for subY := y - 1; subY <= y+1; subY++ {
					powerLevel += grid[subX][subY]
				}
			}

			if powerLevel > maxPowerLevel {
				maxPower = coordinate{x: x, y: y}
				maxPowerLevel = powerLevel
			}

		}
	}

	fmt.Printf("Maximum power was found at: %v,%v with a power of %v.", maxPower.x, maxPower.y, maxPowerLevel)

}
