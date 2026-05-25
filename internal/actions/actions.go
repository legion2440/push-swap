// Пакет actions реализует операции push-swap над двумя стеками.
package actions

import (
	"push-swap/internal/stack"
)

// Выполняет инструкцию над стеками a и b.
// Записывает название инструкции в out, если out не nil.
func Execute(a, b *stack.Stack, instruction string, out []string) []string {
	switch instruction {
	case "pa":
		if b.Len() > 0 {
			a.Push(b.Pop())
		}
		return append(out, "pa")
	case "pb":
		if a.Len() > 0 {
			b.Push(a.Pop())
		}
		return append(out, "pb")
	case "sa":
		a.Swap()
		return append(out, "sa")
	case "sb":
		b.Swap()
		return append(out, "sb")
	case "ss":
		a.Swap()
		b.Swap()
		return append(out, "ss")
	case "ra":
		a.Rotate()
		return append(out, "ra")
	case "rb":
		b.Rotate()
		return append(out, "rb")
	case "rr":
		a.Rotate()
		b.Rotate()
		return append(out, "rr")
	case "rra":
		a.ReverseRotate()
		return append(out, "rra")
	case "rrb":
		b.ReverseRotate()
		return append(out, "rrb")
	case "rrr":
		a.ReverseRotate()
		b.ReverseRotate()
		return append(out, "rrr")
	default:
		return out
	}
}

// ExecuteSilent выполняет инструкцию без добавления в вывод (для checker).
func ExecuteSilent(a, b *stack.Stack, instruction string) {
	switch instruction {
	case "pa":
		if b.Len() > 0 {
			a.Push(b.Pop())
		}
	case "pb":
		if a.Len() > 0 {
			b.Push(a.Pop())
		}
	case "sa":
		a.Swap()
	case "sb":
		b.Swap()
	case "ss":
		a.Swap()
		b.Swap()
	case "ra":
		a.Rotate()
	case "rb":
		b.Rotate()
	case "rr":
		a.Rotate()
		b.Rotate()
	case "rra":
		a.ReverseRotate()
	case "rrb":
		b.ReverseRotate()
	case "rrr":
		a.ReverseRotate()
		b.ReverseRotate()
	}
}
