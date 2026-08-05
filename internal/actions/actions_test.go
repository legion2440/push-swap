// Тесты для пакета actions.
package actions

import (
	"testing"

	"push-swap/internal/stack"
)

func TestExecute_Push(t *testing.T) {
	a := stack.NewFromSlice([]int{1})
	b := stack.NewFromSlice([]int{2})
	ExecuteSilent(a, b, "pb")
	if a.Len() != 0 || b.Len() != 2 {
		t.Errorf("pb: a=%v b=%v", a.Data, b.Data)
	}
	if b.Data[0] != 1 {
		t.Errorf("pb: верх B должен быть 1, получили %d", b.Data[0])
	}
	ExecuteSilent(a, b, "pa")
	if a.Len() != 1 || b.Len() != 1 {
		t.Errorf("pa: a=%v b=%v", a.Data, b.Data)
	}
	if a.Data[0] != 1 {
		t.Errorf("pa: верх A должен быть 1, получили %d", a.Data[0])
	}
}

func TestExecute_Swap(t *testing.T) {
	a := stack.NewFromSlice([]int{2, 1})
	b := stack.New()
	ExecuteSilent(a, b, "sa")
	if a.Data[0] != 1 || a.Data[1] != 2 {
		t.Errorf("sa: получили %v, ожидалось [1 2]", a.Data)
	}
}

func TestExecute_Rotate(t *testing.T) {
	a := stack.NewFromSlice([]int{1, 2, 3})
	b := stack.New()
	ExecuteSilent(a, b, "ra")
	if a.Data[0] != 2 || a.Data[2] != 1 {
		t.Errorf("ra: получили %v, ожидалось [2 3 1]", a.Data)
	}
}

func TestExecute_ReverseRotate(t *testing.T) {
	a := stack.NewFromSlice([]int{1, 2, 3})
	b := stack.New()
	ExecuteSilent(a, b, "rra")
	if a.Data[0] != 3 || a.Data[1] != 1 {
		t.Errorf("rra: получили %v, ожидалось [3 1 2]", a.Data)
	}
}

func TestExecute_BStackInstructions(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want []int
	}{
		{name: "sb", cmd: "sb", want: []int{2, 1, 3}},
		{name: "rb", cmd: "rb", want: []int{2, 3, 1}},
		{name: "rrb", cmd: "rrb", want: []int{3, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := stack.New()
			b := stack.NewFromSlice([]int{1, 2, 3})
			ExecuteSilent(a, b, tt.cmd)
			for i, want := range tt.want {
				if b.Data[i] != want {
					t.Fatalf("%s: b=%v, want %v", tt.cmd, b.Data, tt.want)
				}
			}
		})
	}
}

func TestExecute_AllInstructions(t *testing.T) {
	tests := []struct {
		name  string
		a     []int
		b     []int
		cmd   string
		check func(t *testing.T, a, b *stack.Stack)
	}{
		{"ss", []int{2, 1}, []int{4, 3}, "ss", func(t *testing.T, a, b *stack.Stack) {
			if a.Data[0] != 1 || a.Data[1] != 2 || b.Data[0] != 3 || b.Data[1] != 4 {
				t.Error("ss: оба стека должны быть перевёрнуты")
			}
		}},
		{"rr", []int{1, 2, 3}, []int{4, 5, 6}, "rr", func(t *testing.T, a, b *stack.Stack) {
			if a.Data[0] != 2 || b.Data[0] != 5 {
				t.Error("rr: оба стека должны быть сдвинуты")
			}
		}},
		{"rrr", []int{1, 2, 3}, []int{4, 5, 6}, "rrr", func(t *testing.T, a, b *stack.Stack) {
			if a.Data[0] != 3 || b.Data[0] != 6 {
				t.Error("rrr: оба стека должны быть reverse-сдвинуты")
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := stack.NewFromSlice(tt.a)
			b := stack.NewFromSlice(tt.b)
			ExecuteSilent(a, b, tt.cmd)
			tt.check(t, a, b)
		})
	}
}

func TestExecute_NoOpOnInsufficientStack(t *testing.T) {
	a := stack.NewFromSlice([]int{1})
	b := stack.New()

	for _, cmd := range []string{"pa", "sa", "sb", "ra", "rb", "rra", "rrb"} {
		beforeA := append([]int(nil), a.Data...)
		beforeB := append([]int(nil), b.Data...)
		ExecuteSilent(a, b, cmd)
		if len(a.Data) != len(beforeA) || len(b.Data) != len(beforeB) {
			t.Fatalf("%s changed stack lengths: a=%v b=%v", cmd, a.Data, b.Data)
		}
		for i := range beforeA {
			if a.Data[i] != beforeA[i] {
				t.Fatalf("%s changed A: got %v want %v", cmd, a.Data, beforeA)
			}
		}
		for i := range beforeB {
			if b.Data[i] != beforeB[i] {
				t.Fatalf("%s changed B: got %v want %v", cmd, b.Data, beforeB)
			}
		}
	}
}
