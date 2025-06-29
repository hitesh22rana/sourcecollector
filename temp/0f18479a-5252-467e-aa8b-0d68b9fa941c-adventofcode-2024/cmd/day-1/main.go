package main

import (
	"container/heap"
	"fmt"
	"math"
	"strconv"
	"strings"

	"advent-of-code-2024/pkg/ds/heaps"
	"advent-of-code-2024/pkg/utils"
)

func PartOne(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	leftMinHeap := heaps.NewMinHeap()
	rightMinHeap := heaps.NewMinHeap()

	for _, line := range data {
		lines := strings.Split(string(line), "   ")

		leftVal, _ := strconv.Atoi(lines[0])
		rightVal, _ := strconv.Atoi(lines[1])

		heap.Push(leftMinHeap, leftVal)
		heap.Push(rightMinHeap, rightVal)
	}

	ans := 0
	for leftMinHeap.Len() > 0 && rightMinHeap.Len() > 0 {
		leftVal := heap.Pop(leftMinHeap).(int)
		rightVal := heap.Pop(rightMinHeap).(int)

		distance := math.Abs(float64(leftVal - rightVal))
		ans += int(distance)
	}

	return ans
}

func PartTwo(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	keys := []int{}
	frequency := make(map[int]int)

	for _, line := range data {
		lines := strings.Split(string(line), "   ")

		leftVal, _ := strconv.Atoi(lines[0])
		rightVal, _ := strconv.Atoi(lines[1])

		keys = append(keys, leftVal)
		frequency[rightVal]++
	}

	ans := 0
	for _, key := range keys {
		ans += key * frequency[key]
	}

	return ans
}

func main() {
	var testCases []string

	// First part
	testCases = []string{
		"inputs/day-1/input1.txt",
		"inputs/day-1/input2.txt",
	}
	fmt.Println("Part One")
	for _, tc := range testCases {
		fmt.Println("For ", tc, ":", PartOne(tc))
	}

	// Second part
	testCases = []string{
		"inputs/day-1/input1.txt",
		"inputs/day-1/input2.txt",
	}
	fmt.Println("\nPart Two")
	for _, tc := range testCases {
		fmt.Println("For", tc, ":", PartTwo(tc))
	}
}
