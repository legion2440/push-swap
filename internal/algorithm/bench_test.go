// Тесты производительности и корректности для больших стеков.
package algorithm

import (
	"math/rand"
	"push-swap/internal/actions"
	"push-swap/internal/stack"
	"testing"
)

func verifySort(t *testing.T, input []int, instructions []string) {
	t.Helper()
	a := stack.NewFromSlice(input)
	b := stack.New()
	for _, instr := range instructions {
		actions.ExecuteSilent(a, b, instr)
	}
	if !a.IsSorted() || b.Len() != 0 {
		t.Fatalf("Sort(%v): не отсортировано, %d инструкций", input, len(instructions))
	}
}

func TestSort_Large100(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for trial := 0; trial < 5; trial++ {
		input := make([]int, 100)
		used := make(map[int]bool)
		for i := 0; i < 100; i++ {
			for {
				v := rng.Intn(10000)
				if !used[v] {
					used[v] = true
					input[i] = v
					break
				}
			}
		}
		a := stack.NewFromSlice(input)
		instructions := Sort(a)
		if len(instructions) >= 700 {
			t.Errorf("Sort(100): %d инструкций, лимит 700", len(instructions))
		}
		verifySort(t, input, instructions)
	}
}
