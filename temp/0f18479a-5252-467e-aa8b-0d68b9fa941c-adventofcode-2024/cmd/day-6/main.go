package main

import (
	"fmt"

	"advent-of-code-2024/pkg/utils"
)

type point struct {
	x   int
	y   int
	dir int
}

var directions [][]int = [][]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func dfs(grid *[][]rune, x int, y int, dir *int, vis *map[point]struct{}) {
	if x < 0 || x >= len(*grid) || y < 0 || y >= len((*grid)[0]) {
		return
	}

	if (*grid)[x][y] == '#' {
		x += -1 * directions[*dir][0]
		y += -1 * directions[*dir][1]

		*dir = ((*dir) + 1) % 4
	}

	(*vis)[point{x, y, *dir}] = struct{}{}

	x += directions[*dir][0]
	y += directions[*dir][1]
	dfs(grid, x, y, dir, vis)
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

	x, y := 0, 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j] == '^' {
				x, y = i, j
				break
			}
		}
	}

	ans := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			r, c := x, y

			dir := 0
			seen := make(map[point]struct{})
			for {
				p := point{r, c, dir}
				if _, ok := seen[p]; ok {
					ans++
					break
				}

				seen[p] = struct{}{}

				dr := r + directions[dir][0]
				dc := c + directions[dir][1]
				if dr >= 0 && dr < len(grid) && dc >= 0 && dc < len(grid[0]) {
					if grid[dr][dc] == '#' || (dr == i && dc == j) {
						dir = (dir + 1) % 4
					} else {
						r, c = dr, dc
					}
				} else {
					break
				}
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

	x, y := 0, 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j] == '^' {
				grid[i][j] = '.'
				x, y = i, j
				break
			}
		}
	}

	vis := make(map[point]struct{})
	dir := 0
	dfs(&grid, x, y, &dir, &vis)
	return len(vis)
}

func main() {
	var testCases []string

	// First part
	testCases = []string{
		"inputs/day-6/input1.txt",
		"inputs/day-6/input2.txt",
	}

	fmt.Println("Part One")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartOne(tc))
	}

	// Second part
	testCases = []string{
		"inputs/day-6/input1.txt",
		"inputs/day-6/input2.txt",
	}

	fmt.Println("\nPart Two")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartTwo(tc))
	}
}
