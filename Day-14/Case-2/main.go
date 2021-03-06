package main

import (
	"fmt"
)

type recipe struct {
	score int
	left  *recipe
	right *recipe
}

var recipeRange = []int{5, 4, 0, 5, 6, 1}

func main() {

	recipes := []recipe{}

	// Create all the initial RECIPES and add them to the list
	firstRecipe := recipe{score: 3}
	secondRecipe := recipe{score: 7, left: &firstRecipe, right: &firstRecipe}
	firstRecipe.left = &secondRecipe
	firstRecipe.right = &secondRecipe
	lastRecipe := &secondRecipe
	recipes = append(recipes, firstRecipe)
	recipes = append(recipes, secondRecipe)

	// Create the ELFS and assign a CURRENT RECIPE to them
	firstElf := &firstRecipe
	secondElf := &secondRecipe

OUTER:
	for {

		// Combine the SCORE of the 2 RECIPES and create NEW RECIPES based on the combined SCORE
		var newRecipeScore [2]int
		newRecipeScore[0], newRecipeScore[1] = (firstElf.score+secondElf.score)/10, (firstElf.score+secondElf.score)%10

	INNER:
		for s := 0; s < len(newRecipeScore); s++ {
			if s == 0 && newRecipeScore[s] == 0 {
				continue
			}
			newRecipe := recipe{score: newRecipeScore[s], left: lastRecipe, right: lastRecipe.right}
			lastRecipe.right = &newRecipe
			lastRecipe = &newRecipe
			recipes = append(recipes, newRecipe)

			// Check if last RECIPES match the RECIPE RANGE
			if len(recipes) >= len(recipeRange) {
				for r := len(recipes) - len(recipeRange); r < len(recipes); r++ {
					if recipes[r].score != recipeRange[len(recipeRange)-(len(recipes)-r)] {
						continue INNER
					}
				}
				break OUTER
			}
		}

		// Assign the new CURRENT RECIPES to the ELFS
		firstElf = firstElf.getRecipe(firstElf.score + 1)
		secondElf = secondElf.getRecipe(secondElf.score + 1)

	}

	fmt.Printf("The recipe range first appears after %v recipes.\n", len(recipes)-len(recipeRange))

}

func (current *recipe) getRecipe(offset int) *recipe {
	var actual *recipe
	if offset < 0 {
		actual = current.left.getRecipe(offset + 1)
	} else if offset > 0 {
		actual = current.right.getRecipe(offset - 1)
	} else if offset == 0 {
		actual = current
	}
	return actual
}
