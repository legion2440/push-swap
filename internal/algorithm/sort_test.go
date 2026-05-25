// Тесты для пакета algorithm.
package algorithm

import (
	"push-swap/internal/actions"
	"push-swap/internal/stack"
	"testing"
)

func runInstructions(a, b *stack.Stack, instructions []string) {
	for _, instr := range instructions {
		actions.ExecuteSilent(a, b, instr)
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		input []int
	}{
		{[]int{2, 1}},
		{[]int{3, 2, 1}},
		{[]int{2, 1, 3, 6, 5, 8}},
		{[]int{4, 67, 3, 87, 23}},
		{[]int{1, 2, 3}},
		{[]int{5, 4, 3, 2, 1}},
		{[]int{1}},
		{[]int{0, -1, 1}},
		{[]int{7, 6, 5, 4, 3, 2, 1}},
	}
	for _, tt := range tests {
		a := stack.NewFromSlice(tt.input)
		instructions := Sort(a)
		b := stack.New()
		a2 := stack.NewFromSlice(tt.input)
		runInstructions(a2, b, instructions)
		if !a2.IsSorted() || b.Len() != 0 {
			t.Errorf("Sort(%v): результат не отсортирован. Инструкции: %v", tt.input, instructions)
		}
	}
}

// Лимиты по аудиту: 6 элементов < 9 инструкций, 5 элементов < 12.
func TestSort_InstructionLimits(t *testing.T) {
	tests := []struct {
		input    []int
		maxInstr  int
	}{
		{[]int{2, 1, 3, 6, 5, 8}, 8},
		{[]int{4, 67, 3, 87, 23}, 12},
		{[]int{5, 4, 3, 2, 1}, 11}, // аудит: < 12 инструкций
	}
	for _, tt := range tests {
		a := stack.NewFromSlice(tt.input)
		instructions := Sort(a)
		if len(instructions) > tt.maxInstr {
			t.Errorf("Sort(%v): %d инструкций, лимит %d", tt.input, len(instructions), tt.maxInstr)
		}
	}
}

func TestSortEmptyAndSingle(t *testing.T) {
	a := stack.NewFromSlice([]int{42})
	instructions := Sort(a)
	if len(instructions) != 0 {
		t.Errorf("Sort([42]) должен вернуть пустой список")
	}
}

func TestSortAlreadySorted(t *testing.T) {
	a := stack.NewFromSlice([]int{1, 2, 3, 4, 5})
	instructions := Sort(a)
	if len(instructions) != 0 {
		t.Errorf("Sort(отсортированного) должен вернуть пустой список, получили %v", instructions)
	}
}
