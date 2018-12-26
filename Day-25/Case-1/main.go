package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type constellation struct {
	stars []*star
}

type star struct {
	x int
	y int
	z int
	t int
}

func main() {

	constellations := []*constellation{}
	stars := []*star{}

	file, _ := os.Open("input.txt")

	input := bufio.NewScanner(file)
	for input.Scan() {

		coordinates := strings.Split(input.Text(), ",")
		x, _ := strconv.Atoi(coordinates[0])
		y, _ := strconv.Atoi(coordinates[1])
		z, _ := strconv.Atoi(coordinates[2])
		t, _ := strconv.Atoi(coordinates[3])

		newStar := star{x: x, y: y, z: z, t: t}
		stars = append(stars, &newStar)

	}

	for s := 0; s < len(stars); s++ {

		constellationsInRange := []*constellation{}

		for c := 0; c < len(constellations); c++ {

			for cs := 0; cs < len(constellations[c].stars); cs++ {

				distance := int(math.Abs(float64((constellations[c].stars[cs].t)-stars[s].t))) +
					int(math.Abs(float64((constellations[c].stars[cs].x)-stars[s].x))) +
					int(math.Abs(float64((constellations[c].stars[cs].y)-stars[s].y))) +
					int(math.Abs(float64((constellations[c].stars[cs].z)-stars[s].z)))

				if distance <= 3 {
					constellationsInRange = append(constellationsInRange, constellations[c])
					break
				}
			}

		}

		if len(constellationsInRange) == 0 {
			newConstellation := constellation{stars: []*star{stars[s]}}
			constellations = append(constellations, &newConstellation)
		} else if len(constellationsInRange) == 1 {
			constellationsInRange[0].stars = append(constellationsInRange[0].stars, stars[s])
		} else if len(constellationsInRange) > 1 {
			constellationsInRange[0].stars = append(constellationsInRange[0].stars, stars[s])
			for cR := 1; cR < len(constellationsInRange); cR++ {
				constellationsInRange[0].stars = append(constellationsInRange[0].stars, constellationsInRange[cR].stars...)
				for c := 0; c < len(constellations); c++ {
					if constellationsInRange[cR] == constellations[c] {
						constellations = append(constellations[:c], constellations[c+1:]...)
					}
				}

			}
		}

	}

	fmt.Println(len(constellations))

}
