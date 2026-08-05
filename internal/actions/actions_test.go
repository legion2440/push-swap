// Тесты для пакета actions.
package actions

import (
	"push-swap/internal/stack"
	"testing"
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
