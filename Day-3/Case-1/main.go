package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const fabricWidth = 1000

func main() {

	fabric := make([][]int, fabricWidth)

	for i := 0; i < fabricWidth; i++ {
		fabric[i] = make([]int, fabricWidth)
	}

	input, _ := ioutil.ReadFile("input.txt")

	claims := strings.Split(string(input), "\r\n")

	for _, claim := range claims {

		//number := strings.TrimSpace(claim[:strings.Index(claim, "@")])
		location := strings.TrimSpace(claim[strings.Index(claim, "@")+2 : strings.Index(claim, ":")])
		size := strings.TrimSpace(claim[strings.Index(claim, ":")+2:])
		xS := location[:strings.Index(location, ",")]
		yS := location[strings.Index(location, ",")+1:]
		widthS := size[:strings.Index(size, "x")]
		heightS := size[strings.Index(size, "x")+1:]

		xI, _ := strconv.Atoi(xS)
		yI, _ := strconv.Atoi(yS)
		widthI, _ := strconv.Atoi(widthS)
		heightI, _ := strconv.Atoi(heightS)

		for y := 0; y < heightI; y++ {
			for x := 0; x < widthI; x++ {
				fabric[xI+x][yI+y]++
			}
		}

	}

	totalOverlapping := 0

	for i := 0; i < fabricWidth; i++ {
		for j := 0; j < fabricWidth; j++ {
			if fabric[i][j] == 1 {
				totalOverlapping++
			}
		}
	}

	fmt.Printf("Total overlapping fields: %v\n", totalOverlapping)

}
