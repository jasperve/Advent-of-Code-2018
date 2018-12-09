package main

import "fmt"

const maxNumber = 71626 //71626
const numPlayers = 438  //438

func main() {

	players := make(map[int]int)

	for i := 1; i <= numPlayers; i++ {
		players[i] = 0
	}

	circle := []int{0}
	currentMarble := 0

	nextNumber := 1

	OUTER:
	for {

		for i := 1; i <= numPlayers; i++ {

			if len(circle) == 1 {
				circle = append(circle, nextNumber)
				currentMarble = 1
			} else {

				if nextNumber%23 == 0 {

					currentMarble -= 7

					if currentMarble < 0 {
						currentMarble = currentMarble + len(circle)
					}

					players[i] += nextNumber + circle[currentMarble]
					circle = append(circle[:currentMarble], circle[currentMarble+1:]...)

				} else {

					if currentMarble+2 == len(circle) {
						currentMarble += 2
					} else {
						currentMarble = (currentMarble + 2) % len(circle)
					}

					circle = append(circle, 0)
					copy(circle[currentMarble+1:], circle[currentMarble:])
					circle[currentMarble] = nextNumber

				}

			}

			nextNumber++
			if nextNumber > maxNumber { break OUTER }

		}

	}

	maxScore := 0

	for _, v := range players {
		if v > maxScore {
			maxScore = v
		}
	}

	fmt.Println(maxScore)

}