package main

import (
	"fmt"

	"advent-of-code-2024/pkg/utils"
)

// Directions: N, NE, E, SE, S, SW, W, NW
var dx = []int{-1, -1, 0, 1, 1, 1, 0, -1}
var dy = []int{0, 1, 1, 1, 0, -1, -1, -1}

var diagonalDx = []int{-1, 1, 1, -1}
var diagonalDy = []int{1, 1, -1, -1}

// Word to search
var word string = "XMAS"

func isSafeDiagonally(x int, y int, n int, m int) bool {
	for i := 0; i < 4; i++ {
		_x := x + diagonalDx[i]
		_y := y + diagonalDy[i]

		if _x < 0 || _x >= n || _y < 0 || _y >= m {
			return false
		}
	}

	return true
}

func searchWord(grid [][]rune, word string, x int, y int, dx int, dy int) int {
	for k := 0; k < len(word); k++ {
		_x := x + k*dx
		_y := y + k*dy

		if _x < 0 || _x >= len(grid) || _y < 0 || _y >= len(grid[0]) {
			return 0
		}

		if grid[_x][_y] != rune(word[k]) {
			return 0
		}
	}

	return 1
}

func PartTwo(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	grid := make([][]rune, len(data))
	for i, line := range data {
		grid[i] = []rune(line)
	}

	ans := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] != 'A' || !isSafeDiagonally(i, j, len(grid), len(grid[0])) {
				continue
			}

			if ((grid[i-1][j-1] == 'M' && grid[i+1][j+1] == 'S') || (grid[i-1][j-1] == 'S' && grid[i+1][j+1] == 'M')) && ((grid[i-1][j+1] == 'M' && grid[i+1][j-1] == 'S') || (grid[i-1][j+1] == 'S' && grid[i+1][j-1] == 'M')) {
				ans++
			}
		}
	}

	return ans
}

func PartOne(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	grid := make([][]rune, len(data))
	for i, line := range data {
		grid[i] = []rune(line)
	}

	ans := 0
	for x := 0; x < len(grid); x++ {
		for y := 0; y < len(grid[0]); y++ {
			for dir := 0; dir < 8; dir++ {
				ans += searchWord(grid, word, x, y, dx[dir], dy[dir])
			}
		}
	}

	return ans
}

func main() {
	var testCases []string

	// First part
	testCases = []string{
		"inputs/day-4/input1.txt",
		"inputs/day-4/input2.txt",
	}

	fmt.Println("Part One")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartOne(tc))
	}

	// Second part
	testCases = []string{
		"inputs/day-4/input1.txt",
		"inputs/day-4/input2.txt",
	}

	fmt.Println("\nPart Two")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartTwo(tc))
	}
}
