package main

import (
	"fmt"
	"io/ioutil"
	//"strconv"
	"strings"
	"regexp"
)

const fabricWidth = 1000

func main() {

	input, _ := ioutil.ReadFile("input.txt")
	claims := strings.Split(string(input), "\r\n")

	fabric := make([][]int, fabricWidth)
	for i := 0; i < fabricWidth; i++ {
		fabric[i] = make([]int, fabricWidth)
	}

	for i := 1; i <= 2; i++ {

		//OUTER:
		for _, claim := range claims {

			matches := regexp.MustCompile("#\\d+ @ (\\d+),(\\d+): (\\d+)x(\\d+)").FindAllStringSubmatch(claim, -1)

			fmt.Println(matches[1])
			fmt.Printf("%T", matches[1])


			/*xI, _ := strconv.Atoi(matches[1])
			yI, _ := strconv.Atoi(matches[2])
			widthI, _ := strconv.Atoi(matches[3])
			heightI, _ := strconv.Atoi(matches[4])

			if i == 1 {

				for y := 0; y < heightI; y++ {
					for x := 0; x < widthI; x++ {
						fabric[xI+x][yI+y]++
					}
				}

			} else {

				for y := 0; y < heightI; y++ {
					for x := 0; x < widthI; x++ {
						if fabric[xI+x][yI+y] > 1 {
							continue OUTER
						}
					}
				}

				fmt.Printf("Number with no overlaps: %v", matches[0])

			}
*/
		}

	}
}