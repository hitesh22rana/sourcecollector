package main

import (
	"fmt"

	"advent-of-code-2024/pkg/utils"
)

var antiNode rune = '#'

type Position struct {
	x int
	y int
}

func isSafe(grid *[][]rune, x int, y int) bool {
	return x >= 0 && x < len(*grid) && y >= 0 && y < len((*grid)[0])
}

func isAntenna(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= 'a' && char <= 'z')
}

func makeAntiNodesPartTwo(grid *[][]rune, antennas *[]Position, antiNodes *map[Position]struct{}) {
	for i := 0; i < len(*antennas); i++ {
		firstAntenna := (*antennas)[i]

		for j := i + 1; j < len(*antennas); j++ {
			secondAntenna := (*antennas)[j]

			diff := Position{
				x: secondAntenna.x - firstAntenna.x,
				y: secondAntenna.y - firstAntenna.y,
			}

			var candidate Position

			candidate = firstAntenna
			for isSafe(grid, candidate.x, candidate.y) {
				if _, ok := (*antiNodes)[candidate]; !ok {
					(*antiNodes)[candidate] = struct{}{}
				}
				candidate.x -= diff.x
				candidate.y -= diff.y
			}

			candidate = secondAntenna
			for isSafe(grid, candidate.x, candidate.y) {
				if _, ok := (*antiNodes)[candidate]; !ok {
					(*antiNodes)[candidate] = struct{}{}
				}
				candidate.x += diff.x
				candidate.y += diff.y
			}
		}
	}
}

func makeAntiNodesPartOne(grid *[][]rune, antennas *[]Position) {
	for i := 0; i < len(*antennas); i++ {
		firstAntenna := (*antennas)[i]
		for j := i + 1; j < len(*antennas); j++ {
			secondAntenna := (*antennas)[j]

			antiNode1 := Position{
				x: 2*(firstAntenna.x-secondAntenna.x) + secondAntenna.x,
				y: 2*(firstAntenna.y-secondAntenna.y) + secondAntenna.y,
			}
			antiNode2 := Position{
				x: 2*(secondAntenna.x-firstAntenna.x) + firstAntenna.x,
				y: 2*(secondAntenna.y-firstAntenna.y) + firstAntenna.y,
			}

			if isSafe(grid, antiNode1.x, antiNode1.y) {
				(*grid)[antiNode1.x][antiNode1.y] = antiNode
			}

			if isSafe(grid, antiNode2.x, antiNode2.y) {
				(*grid)[antiNode2.x][antiNode2.y] = antiNode
			}
		}
	}
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

	antennasMap := make(map[rune][]Position, 0)
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if isAntenna(grid[i][j]) {
				if _, ok := antennasMap[grid[i][j]]; !ok {
					antennasMap[grid[i][j]] = make([]Position, 0)
				}
				antennasMap[grid[i][j]] = append(
					antennasMap[grid[i][j]],
					Position{x: i, y: j},
				)
			}
		}
	}

	antiNodes := make(map[Position]struct{}, 0)
	for _, antennas := range antennasMap {
		makeAntiNodesPartTwo(&grid, &antennas, &antiNodes)
	}

	return len(antiNodes)
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

	antennasMap := make(map[rune][]Position, 0)
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if isAntenna(grid[i][j]) {
				if _, ok := antennasMap[grid[i][j]]; !ok {
					antennasMap[grid[i][j]] = make([]Position, 0)
				}
				antennasMap[grid[i][j]] = append(
					antennasMap[grid[i][j]],
					Position{x: i, y: j},
				)
			}
		}
	}

	for _, antennas := range antennasMap {
		makeAntiNodesPartOne(&grid, &antennas)
	}

	ans := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] == antiNode {
				ans++
			}
		}
	}

	return ans
}

func main() {
	var testCases []string

	// First part
	testCases = []string{
		"inputs/day-8/input1.txt",
		"inputs/day-8/input2.txt",
	}

	fmt.Println("Part One")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartOne(tc))
	}

	// Second part
	testCases = []string{
		"inputs/day-8/input1.txt",
		"inputs/day-8/input2.txt",
	}

	fmt.Println("\nPart Two")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartTwo(tc))
	}
}
