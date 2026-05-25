// Пакет stack реализует структуру данных стек для проекта push-swap.
package stack

import (
	"strconv"
)

// Стек целых чисел. Верхний элемент — первый в слайсе.
type Stack struct {
	Data []int
}

// New создаёт новый пустой стек.
func New() *Stack {
	return &Stack{
		Data: make([]int, 0),
	}
}

// NewFromSlice создаёт стек из слайса. Первый элемент — верх стека.
func NewFromSlice(data []int) *Stack {
	s := &Stack{
		Data: make([]int, len(data)),
	}
	copy(s.Data, data)
	return s
}

// Len возвращает количество элементов в стеке.
func (s *Stack) Len() int {
	return len(s.Data)
}

// Top возвращает верхний элемент стека. Паника при пустом стеке.
func (s *Stack) Top() int {
	return s.Data[0]
}

// Push добавляет элемент на верх стека.
func (s *Stack) Push(val int) {
	s.Data = append([]int{val}, s.Data...)
}

// Pop удаляет и возвращает верхний элемент. Паника при пустом стеке.
func (s *Stack) Pop() int {
	val := s.Data[0]
	s.Data = s.Data[1:]
	return val
}

// Swap меняет местами два верхних элемента.
func (s *Stack) Swap() {
	if len(s.Data) < 2 {
		return
	}
	s.Data[0], s.Data[1] = s.Data[1], s.Data[0]
}

// Rotate сдвигает все элементы вверх: первый становится последним.
func (s *Stack) Rotate() {
	if len(s.Data) < 2 {
		return
	}
	first := s.Data[0]
	s.Data = append(s.Data[1:], first)
}

// ReverseRotate сдвигает все элементы вниз: последний становится первым.
func (s *Stack) ReverseRotate() {
	if len(s.Data) < 2 {
		return
	}
	last := s.Data[len(s.Data)-1]
	s.Data = append([]int{last}, s.Data[:len(s.Data)-1]...)
}

// Copy создаёт копию стека.
func (s *Stack) Copy() *Stack {
	return NewFromSlice(s.Data)
}

// IsSorted проверяет, отсортирован ли стек по возрастанию.
func (s *Stack) IsSorted() bool {
	for i := 0; i < len(s.Data)-1; i++ {
		if s.Data[i] > s.Data[i+1] {
			return false
		}
	}
	return true
}

// String возвращает строковое представление стека (для отладки).
func (s *Stack) String() string {
	if len(s.Data) == 0 {
		return "[]"
	}
	result := "[" + strconv.Itoa(s.Data[0])
	for i := 1; i < len(s.Data); i++ {
		result += " " + strconv.Itoa(s.Data[i])
	}
	return result + "]"
}
