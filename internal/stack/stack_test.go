// Тесты для пакета stack.
package stack

import (
	"testing"
)

func TestStackPushPop(t *testing.T) {
	s := New()
	s.Push(1)
	s.Push(2)
	if s.Len() != 2 {
		t.Errorf("Len() = %d, ожидалось 2", s.Len())
	}
	if s.Top() != 2 {
		t.Errorf("Top() = %d, ожидалось 2", s.Top())
	}
	val := s.Pop()
	if val != 2 {
		t.Errorf("Pop() = %d, ожидалось 2", val)
	}
}

func TestStackSwap(t *testing.T) {
	s := NewFromSlice([]int{2, 1, 3})
	s.Swap()
	if s.Data[0] != 1 || s.Data[1] != 2 {
		t.Errorf("Swap: получили %v, ожидалось [1 2 3]", s.Data)
	}
}

func TestStackRotate(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3})
	s.Rotate()
	if s.Data[0] != 2 || s.Data[2] != 1 {
		t.Errorf("Rotate: получили %v, ожидалось [2 3 1]", s.Data)
	}
}

func TestStackReverseRotate(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3})
	s.ReverseRotate()
	if s.Data[0] != 3 || s.Data[1] != 1 {
		t.Errorf("ReverseRotate: получили %v, ожидалось [3 1 2]", s.Data)
	}
}

func TestStackCopy(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3})
	c := s.Copy()
	c.Swap()
	if s.Data[0] != 1 {
		t.Error("Copy: оригинал не должен измениться")
	}
	if c.Data[0] != 2 {
		t.Error("Copy: копия должна измениться")
	}
}

func TestStackEmptyOperations(t *testing.T) {
	s := New()
	if s.Len() != 0 {
		t.Error("Пустой стек: Len должен быть 0")
	}
	s.Swap()
	s.Rotate()
	s.ReverseRotate()
	if s.Len() != 0 {
		t.Error("Операции на пустом стеке не должны менять Len")
	}
}

func TestStackIsSorted(t *testing.T) {
	sorted := NewFromSlice([]int{1, 2, 3})
	if !sorted.IsSorted() {
		t.Error("IsSorted() должен быть true для [1 2 3]")
	}
	unsorted := NewFromSlice([]int{3, 1, 2})
	if unsorted.IsSorted() {
		t.Error("IsSorted() должен быть false для [3 1 2]")
	}
}
