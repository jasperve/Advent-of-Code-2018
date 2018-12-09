package main

import (
	"container/ring"
	"fmt"
)

func main() {

	nPlayers := 438
	lastMarble := 716260000

	players := make([]int, nPlayers)
	r := &ring.Ring{Value: 0}
	for i := 1; i <= lastMarble; i++ {
		if i%23 == 0 {
			r = r.Move(-8)
			players[(i-1)%nPlayers] += i + r.Unlink(1).Value.(int)
			r = r.Next()
		} else {
			r = r.Next().Link(&ring.Ring{Value: i}).Prev()
		}
	}
	var highScore int
	for _, s := range players {
		if s > highScore {
			highScore = s
		}
	}
	fmt.Println(highScore)

}
