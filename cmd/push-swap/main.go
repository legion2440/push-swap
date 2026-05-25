// Программа push-swap вычисляет минимальную последовательность
// инструкций для сортировки стека целых чисел.
package main

import (
	"fmt"
	"os"

	"push-swap/internal/algorithm"
	"push-swap/internal/validation"
	"push-swap/internal/stack"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	nums, err := validation.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	if len(nums) == 0 {
		return
	}

	a := stack.NewFromSlice(nums)
	instructions := algorithm.Sort(a)

	for _, instr := range instructions {
		fmt.Println(instr)
	}
}
