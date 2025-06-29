package main

import (
	"fmt"
	"strconv"
	"strings"

	"advent-of-code-2024/pkg/utils"
)

func PartTwo(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	graphLowToHigh := make(map[int]map[int]struct{})
	graphHighToLow := make(map[int]map[int]struct{})
	index := 0
	for _, line := range data {
		if len(line) == 0 {
			break
		}
		index++

		pages := strings.Split(line, "|")
		highPriority, err := strconv.Atoi(pages[0])
		if err != nil {
			panic(err)
		}

		lowPriority, err := strconv.Atoi(pages[1])
		if err != nil {
			panic(err)
		}

		if _, ok := graphLowToHigh[lowPriority]; !ok {
			graphLowToHigh[lowPriority] = make(map[int]struct{})
		}
		graphLowToHigh[lowPriority][highPriority] = struct{}{}

		if _, ok := graphHighToLow[highPriority]; !ok {
			graphHighToLow[highPriority] = make(map[int]struct{})
		}
		graphHighToLow[highPriority][lowPriority] = struct{}{}
	}

	ans := 0
	for _, line := range data[index+1:] {
		pagesStr := strings.Split(line, ",")
		pages := make([]int, 0)
		for _, page := range pagesStr {
			pageNum, err := strconv.Atoi(page)
			if err != nil {
				panic(err)
			}

			pages = append(pages, pageNum)
		}

		flag := false
		for i := 0; i < len(pages); i++ {
			for j := i + 1; j < len(pages); j++ {
				if _, ok := graphLowToHigh[pages[i]][pages[j]]; ok {
					flag = true
					break
				}
			}
		}

		if !flag {
			continue
		}

		good := make([]int, 0)
		queue := make([]int, 0)

		intersection := make(map[int]int)
		vsSet := make(map[int]struct{})
		for _, v := range pages {
			vsSet[v] = struct{}{}
		}

		for _, v := range pages {
			count := 0
			if predecessors, exists := graphLowToHigh[v]; exists {
				for pred := range predecessors {
					if _, inVs := vsSet[pred]; inVs {
						count++
					}
				}
			}
			intersection[v] = count
			if count == 0 {
				queue = append(queue, v)
			}
		}

		for len(queue) > 0 {
			x := queue[0]
			queue = queue[1:]
			good = append(good, x)
			if successors, exists := graphHighToLow[x]; exists {
				for y := range successors {
					if _, inD := intersection[y]; inD {
						intersection[y]--
						if intersection[y] == 0 {
							queue = append(queue, y)
						}
					}
				}
			}
		}

		ans += good[len(good)/2]
	}

	return ans
}

func PartOne(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	graph := make(map[int]map[int]struct{})
	index := 0
	for _, line := range data {
		if len(line) == 0 {
			break
		}
		index++

		pages := strings.Split(line, "|")
		highPriority, err := strconv.Atoi(pages[0])
		if err != nil {
			panic(err)
		}

		lowPriority, err := strconv.Atoi(pages[1])
		if err != nil {
			panic(err)
		}

		if _, ok := graph[lowPriority]; !ok {
			graph[lowPriority] = make(map[int]struct{})
		}
		graph[lowPriority][highPriority] = struct{}{}
	}

	ans := 0
	for _, line := range data[index+1:] {
		pagesStr := strings.Split(line, ",")
		pages := make([]int, 0)
		for _, page := range pagesStr {
			pageNum, err := strconv.Atoi(page)
			if err != nil {
				panic(err)
			}

			pages = append(pages, pageNum)
		}

		flag := false
		for i := 0; i < len(pages); i++ {
			for j := i + 1; j < len(pages); j++ {
				if _, ok := graph[pages[i]][pages[j]]; ok {
					flag = true
					break
				}
			}
		}

		if !flag {
			ans += pages[len(pages)/2]
		}
	}

	return ans
}

func main() {
	var testCases []string

	// First part
	testCases = []string{
		"inputs/day-5/input1.txt",
		"inputs/day-5/input2.txt",
	}

	fmt.Println("Part One")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartOne(tc))
	}

	// Second part
	testCases = []string{
		"inputs/day-5/input1.txt",
		"inputs/day-5/input2.txt",
	}

	fmt.Println("\nPart Two")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartTwo(tc))
	}
}
