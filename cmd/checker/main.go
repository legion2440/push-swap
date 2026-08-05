// Программа checker проверяет, отсортирован ли стек после
// выполнения последовательности инструкций из stdin.
package main

import (
	"fmt"
	"os"

	"push-swap/internal/actions"
	"push-swap/internal/stack"
	"push-swap/internal/validation"
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

	instructions, err := validation.ParseInstructions()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	a := stack.NewFromSlice(nums)
	b := stack.New()

	for _, instr := range instructions {
		actions.ExecuteSilent(a, b, instr)
	}

	if a.IsSorted() && b.Len() == 0 {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}
