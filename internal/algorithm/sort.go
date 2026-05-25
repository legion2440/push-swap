// Пакет algorithm реализует алгоритм сортировки для push-swap.
package algorithm

import (
	"math"
	"push-swap/internal/stack"
	"sort"
)

// Sort возвращает минимальную последовательность инструкций для сортировки стека.
func Sort(a *stack.Stack) []string {
	if a.Len() <= 1 {
		return nil
	}
	if a.IsSorted() {
		return nil
	}

	b := stack.New()
	result := make([]string, 0)

	switch a.Len() {
	case 2:
		result = sort2(a, result)
	case 3:
		result = sort3(a, result)
	case 4:
		result = sort4(a, b, result)
	case 5:
		result = sort5(a, b, result)
	case 6:
		result = sort6(a, b, result)
	default:
		result = sortLarge(a, b, result)
	}

	return result
}

// sort2 сортирует 2 элемента.
func sort2(a *stack.Stack, result []string) []string {
	if a.Top() > a.Data[1] {
		a.Swap()
		result = append(result, "sa")
	}
	return result
}

// sort3 сортирует 3 элемента оптимальной последовательностью.
func sort3(a *stack.Stack, result []string) []string {
	x, y, z := a.Data[0], a.Data[1], a.Data[2]

	// Все 6 перестановок
	switch {
	case x < y && y < z: // 1 2 3 — уже отсортировано
	case x < y && z < x: // 2 3 1
		a.ReverseRotate()
		result = append(result, "rra")
	case x < y && z < y: // 1 3 2
		a.ReverseRotate()
		result = append(result, "rra")
		a.Swap()
		result = append(result, "sa")
	case y < x && x < z: // 2 1 3
		a.Swap()
		result = append(result, "sa")
	case y < x && z < y: // 3 2 1 — проверяем до z<x, т.к. оба true для 3 2 1
		a.Rotate()
		result = append(result, "ra")
		a.Swap()
		result = append(result, "sa")
	case y < x && z < x: // 3 1 2
		a.Rotate()
		result = append(result, "ra")
	}
	return result
}

// sort4: выталкиваем 1 минимальный в B, сортируем 3, возвращаем.
func sort4(a *stack.Stack, b *stack.Stack, result []string) []string {
	minIdx := 0
	for i := 1; i < a.Len(); i++ {
		if a.Data[i] < a.Data[minIdx] {
			minIdx = i
		}
	}
	for minIdx > 0 {
		a.Rotate()
		result = append(result, "ra")
		minIdx--
	}
	if a.Len() > 0 {
		b.Push(a.Pop())
		result = append(result, "pb")
	}
	result = sort3(a, result)
	if b.Len() > 0 {
		a.Push(b.Pop())
		result = append(result, "pa")
	}
	return result
}

// sort5: выталкиваем 2 минимальных в B, сортируем 3, возвращаем обратно.
func sort5(a *stack.Stack, b *stack.Stack, result []string) []string {
	// Оптимальный случай: стек полностью убывает (например 5 4 3 2 1).
	if isDescending(a.Data) {
		a.ReverseRotate()
		result = append(result, "rra")
		b.Push(a.Pop())
		result = append(result, "pb")
		a.ReverseRotate()
		result = append(result, "rra")
		b.Push(a.Pop())
		result = append(result, "pb")
		result = sort3(a, result)
		for b.Len() > 0 {
			a.Push(b.Pop())
			result = append(result, "pa")
		}
		return result
	}

	m1, m2 := findTwoSmallest(a.Data)

	// Два минимальных уже сверху (например ... 2 1).
	if a.Len() >= 2 && a.Data[0] == m1 && a.Data[1] == m2 {
		b.Push(a.Pop())
		result = append(result, "pb")
		b.Push(a.Pop())
		result = append(result, "pb")
		result = sort3(a, result)
		for b.Len() > 0 {
			a.Push(b.Pop())
			result = append(result, "pa")
		}
		return result
	}

	for _, target := range []int{m1, m2} {
		for a.Data[0] != target {
			a.Rotate()
			result = append(result, "ra")
		}
		b.Push(a.Pop())
		result = append(result, "pb")
	}
	result = sort3(a, result)
	for b.Len() > 0 {
		a.Push(b.Pop())
		result = append(result, "pa")
	}
	if a.Len() >= 2 && a.Data[0] > a.Data[1] {
		a.Swap()
		result = append(result, "sa")
	}
	return result
}

// isDescending проверяет, убывает ли стек сверху вниз.
func isDescending(data []int) bool {
	if len(data) < 2 {
		return false
	}
	for i := 0; i < len(data)-1; i++ {
		if data[i] <= data[i+1] {
			return false
		}
	}
	return true
}

// sort6: оптимизированная сортировка 6 элементов.
func sort6(a *stack.Stack, b *stack.Stack, result []string) []string {
	// Специальный случай: два верхних — два минимальных (в обратном порядке).
	// Оптимальная последовательность: pb pb ra sa rrr pa pa = 7 инструкций.
	min1, min2 := findTwoSmallest(a.Data)
	if a.Data[0] == min2 && a.Data[1] == min1 {
		b.Push(a.Pop())
		result = append(result, "pb")
		b.Push(a.Pop())
		result = append(result, "pb")
		a.Rotate()
		result = append(result, "ra")
		a.Swap()
		result = append(result, "sa")
		a.ReverseRotate()
		b.ReverseRotate()
		result = append(result, "rrr")
		for b.Len() > 0 {
			a.Push(b.Pop())
			result = append(result, "pa")
		}
		return result
	}
	// Общий случай: выталкиваем 2 минимальных
	for pushCount := 0; pushCount < 2; pushCount++ {
		minIdx := 0
		for i := 1; i < a.Len(); i++ {
			if a.Data[i] < a.Data[minIdx] {
				minIdx = i
			}
		}
		for minIdx > 0 {
			a.Rotate()
			result = append(result, "ra")
			minIdx--
		}
		if a.Len() > 0 {
			b.Push(a.Pop())
			result = append(result, "pb")
		}
	}
	result = sort4(a, b, result)
	for b.Len() > 0 {
		a.Push(b.Pop())
		result = append(result, "pa")
	}
	if a.Len() >= 2 && a.Data[0] > a.Data[1] {
		a.Swap()
		result = append(result, "sa")
	}
	return result
}

// findTwoSmallest возвращает два минимальных значения из слайса.
func findTwoSmallest(data []int) (int, int) {
	if len(data) < 2 {
		return 0, 0
	}
	m1, m2 := data[0], data[1]
	if m2 < m1 {
		m1, m2 = m2, m1
	}
	for i := 2; i < len(data); i++ {
		if data[i] < m1 {
			m2 = m1
			m1 = data[i]
		} else if data[i] < m2 {
			m2 = data[i]
		}
	}
	return m1, m2
}

// sortLarge использует алгоритм LIS для больших стеков (7+ элементов).
func sortLarge(a *stack.Stack, b *stack.Stack, result []string) []string {
	rankMap := buildRankMap(a.Data)

	// Фаза 1: находим самую длинную возрастающую подпоследовательность (LIS)
	keep := computeLIS(a.Data, rankMap)

	// Фаза 2: отправляем в B всё, что не входит в LIS
	result = pushNonLISToB(a, b, rankMap, keep, result)

	// Фаза 3: возвращаем из B в A с минимальной стоимостью поворотов
	result = sendToA(a, b, rankMap, result)

	// Фаза 4: минимальный элемент — наверх стека
	result = rotateToMin(a, rankMap, result)

	return result
}

// buildRankMap присваивает каждому значению ранг от 0 (минимум) до n-1.
func buildRankMap(data []int) map[int]int {
	sorted := make([]int, len(data))
	copy(sorted, data)
	sort.Ints(sorted)
	rankMap := make(map[int]int, len(data))
	for i, v := range sorted {
		rankMap[v] = i
	}
	return rankMap
}

// computeLIS находит ранги, входящие в самую длинную возрастающую подпоследовательность.
func computeLIS(data []int, rankMap map[int]int) map[int]bool {
	n := len(data)
	ranks := make([]int, n)
	for i, v := range data {
		ranks[i] = rankMap[v]
	}

	// tailIdx[i] — индекс последнего элемента LIS длины i+1
	tailIdx := make([]int, n)
	prev := make([]int, n)
	for i := range prev {
		prev[i] = -1
	}

	size := 0
	for i, rank := range ranks {
		lo, hi := 0, size
		for lo < hi {
			mid := (lo + hi) / 2
			if ranks[tailIdx[mid]] < rank {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		if lo > 0 {
			prev[i] = tailIdx[lo-1]
		}
		tailIdx[lo] = i
		if lo == size {
			size++
		}
	}

	keep := make(map[int]bool)
	for k := tailIdx[size-1]; k >= 0; k = prev[k] {
		keep[ranks[k]] = true
	}
	return keep
}

// pushNonLISToB отправляет в B элементы, не входящие в LIS.
func pushNonLISToB(a, b *stack.Stack, rankMap map[int]int, keep map[int]bool, result []string) []string {
	size := a.Len()
	for i := 0; i < size; i++ {
		rank := rankMap[a.Top()]
		if keep[rank] {
			a.Rotate()
			result = append(result, "ra")
		} else {
			b.Push(a.Pop())
			result = append(result, "pb")
		}
	}
	return result
}

// sendToA возвращает элементы из B в A, выбирая каждый раз самый дешёвый ход.
func sendToA(a, b *stack.Stack, rankMap map[int]int, result []string) []string {
	for b.Len() > 0 {
		bIdx, aIdx, bRev, aRev := findBestMove(a, b, rankMap)
		result = applyRotations(a, b, aIdx, aRev, bIdx, bRev, result)
		a.Push(b.Pop())
		result = append(result, "pa")
	}
	return result
}

// findBestMove ищет элемент в B с минимальной суммарной стоимостью вставки в A.
func findBestMove(a, b *stack.Stack, rankMap map[int]int) (bIdx, aIdx int, bRev, aRev bool) {
	bestCost := math.MaxInt
	bSize := b.Len()

	for i := 0; i < bSize; i++ {
		rank := rankMap[b.Data[i]]
		bRot, bReverse := rotationCost(i, bSize)
		aRot, aReverse, targetIdx := insertCost(a, rankMap, rank)
		cost := combinedCost(bRot, aRot, bReverse, aReverse)

		if cost < bestCost {
			bestCost = cost
			bIdx = i
			aIdx = targetIdx
			bRev = bReverse
			aRev = aReverse
		}
	}
	return
}

// rotationCost возвращает минимальное число поворотов и направление.
func rotationCost(idx, size int) (cost int, reverse bool) {
	if idx <= size/2 {
		return idx, false
	}
	return size - idx, true
}

// insertCost находит позицию вставки ранга в отсортированный стек A.
func insertCost(a *stack.Stack, rankMap map[int]int, rank int) (cost int, reverse bool, targetIdx int) {
	size := a.Len()
	if size == 0 {
		return 0, false, 0
	}
	targetIdx = findInsertIndex(a, rankMap, rank)
	cost, reverse = rotationCost(targetIdx, size)
	return
}

// findInsertIndex находит индекс, куда нужно вставить элемент с заданным рангом.
func findInsertIndex(a *stack.Stack, rankMap map[int]int, rank int) int {
	size := a.Len()

	for i := 0; i < size; i++ {
		cur := rankMap[a.Data[i]]
		next := rankMap[a.Data[(i+1)%size]]

		if cur < next {
			// Обычный участок без перехода max→min
			if rank > cur && rank < next {
				return (i + 1) % size
			}
		} else {
			// Граница цикла: cur — максимум, next — минимум
			if rank > cur || rank < next {
				return (i + 1) % size
			}
		}
	}
	return 0
}

// combinedCost считает суммарную стоимость поворотов обоих стеков.
func combinedCost(bRot, aRot int, bRev, aRev bool) int {
	if bRev == aRev {
		if aRot > bRot {
			return aRot
		}
		return bRot
	}
	return bRot + aRot
}

// applyRotations выполняет повороты обоих стеков с оптимизацией rr/rrr.
func applyRotations(a, b *stack.Stack, aIdx int, aRev bool, bIdx int, bRev bool, result []string) []string {
	aCnt, _ := rotationCost(aIdx, a.Len())
	bCnt, _ := rotationCost(bIdx, b.Len())

	if !aRev && !bRev {
		common := aCnt
		if bCnt < common {
			common = bCnt
		}
		for i := 0; i < common; i++ {
			a.Rotate()
			b.Rotate()
			result = append(result, "rr")
		}
		for i := common; i < aCnt; i++ {
			a.Rotate()
			result = append(result, "ra")
		}
		for i := common; i < bCnt; i++ {
			b.Rotate()
			result = append(result, "rb")
		}
	} else if aRev && bRev {
		common := aCnt
		if bCnt < common {
			common = bCnt
		}
		for i := 0; i < common; i++ {
			a.ReverseRotate()
			b.ReverseRotate()
			result = append(result, "rrr")
		}
		for i := common; i < aCnt; i++ {
			a.ReverseRotate()
			result = append(result, "rra")
		}
		for i := common; i < bCnt; i++ {
			b.ReverseRotate()
			result = append(result, "rrb")
		}
	} else if aRev {
		for i := 0; i < aCnt; i++ {
			a.ReverseRotate()
			result = append(result, "rra")
		}
		for i := 0; i < bCnt; i++ {
			b.Rotate()
			result = append(result, "rb")
		}
	} else {
		for i := 0; i < aCnt; i++ {
			a.Rotate()
			result = append(result, "ra")
		}
		for i := 0; i < bCnt; i++ {
			b.ReverseRotate()
			result = append(result, "rrb")
		}
	}
	return result
}

// rotateToMin поворачивает стек A так, чтобы минимальный элемент был сверху.
func rotateToMin(a *stack.Stack, rankMap map[int]int, result []string) []string {
	size := a.Len()
	minIdx := 0
	for i := 1; i < size; i++ {
		if rankMap[a.Data[i]] < rankMap[a.Data[minIdx]] {
			minIdx = i
		}
	}
	cost, reverse := rotationCost(minIdx, size)
	if reverse {
		for i := 0; i < cost; i++ {
			a.ReverseRotate()
			result = append(result, "rra")
		}
	} else {
		for i := 0; i < cost; i++ {
			a.Rotate()
			result = append(result, "ra")
		}
	}
	return result
}
