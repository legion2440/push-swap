// Тесты производительности и корректности для больших стеков.
package algorithm

import (
	"math/rand"
	"testing"

	"push-swap/internal/actions"
	"push-swap/internal/stack"
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

func checkLarge100(t *testing.T, name string, input []int) int {
	t.Helper()
	a := stack.NewFromSlice(input)
	instructions := Sort(a)
	if len(instructions) >= 700 {
		t.Fatalf("%s: %d инструкций, требуется < 700", name, len(instructions))
	}
	verifySort(t, input, instructions)
	return len(instructions)
}

func TestSort_Large100Stress(t *testing.T) {
	maxInstructions := 0
	maxSeed := int64(-1)

	// Детерминированные случайные перестановки: значения не важны,
	// потому что large solver работает через относительные ранги.
	for seed := int64(0); seed < 500; seed++ {
		rng := rand.New(rand.NewSource(seed))
		input := rng.Perm(100)
		count := checkLarge100(t, "random seed "+formatSeed(seed), input)
		if count > maxInstructions {
			maxInstructions = count
			maxSeed = seed
		}
	}

	t.Logf("500 random permutations: max=%d instructions (seed=%d)", maxInstructions, maxSeed)
}

func TestSort_Large100StructuredCases(t *testing.T) {
	reverse := make([]int, 100)
	rotated := make([]int, 100)
	zigzag := make([]int, 0, 100)

	for i := 0; i < 100; i++ {
		reverse[i] = 99 - i
		rotated[i] = (i + 50) % 100
	}
	for low, high := 0, 99; low <= high; low, high = low+1, high-1 {
		zigzag = append(zigzag, low)
		if low != high {
			zigzag = append(zigzag, high)
		}
	}

	cases := []struct {
		name  string
		input []int
	}{
		{name: "reverse", input: reverse},
		{name: "rotated", input: rotated},
		{name: "zigzag", input: zigzag},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			count := checkLarge100(t, tc.name, tc.input)
			t.Logf("%s: %d instructions", tc.name, count)
		})
	}
}

func formatSeed(seed int64) string {
	if seed == 0 {
		return "0"
	}

	buf := make([]byte, 0, 20)
	for seed > 0 {
		buf = append(buf, byte('0'+seed%10))
		seed /= 10
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
