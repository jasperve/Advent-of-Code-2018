package main

import (
	"fmt"
)

func main() {

	results := []int{}
	register := []int{0,0,0,0,0,0}

	LOOP1:
	for {

		register[3] = register[5] | 65536
		register[5] = 9010242

		for {

			register[1] = register[3] & 255
			register[5] += register[1]
			register[5] = ((register[5] & 16777215) * 65899) & 16777215

			if 256 > register[3] {
				
				for r := 0; r < len(results); r++ {
					if results[r] == register[5] {
						break LOOP1
					}
				}
				results = append(results, register[5])
				continue LOOP1

			} else {

				register[1] = 0

				LOOP3:
				for {
				
					register[4] = register[1] + 1 
					register[4] *= 256

					if register[4] > register[3] {
						register[3] = register[1]
						break LOOP3
					} else {
						register[1]++
					}
							
				}

			}

		}

	}

	fmt.Println(results[len(results)-1])

}