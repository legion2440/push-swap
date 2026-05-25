// Пакет validation отвечает за парсинг и проверку входных данных.
package validation

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseArgs парсит аргументы командной строки в слайс целых чисел.
// Аргументы могут быть переданы как "2 1 3" или как отдельные "2" "1" "3".
// Возвращает ошибку при невалидных данных или дубликатах.
func ParseArgs(args []string) ([]int, error) {
	if len(args) == 0 {
		return nil, nil
	}

	// Объединяем все аргументы и разбиваем по пробелам
	var allParts []string
	for _, arg := range args {
		parts := strings.Fields(arg)
		allParts = append(allParts, parts...)
	}

	if len(allParts) == 0 {
		return nil, nil
	}

	result := make([]int, 0, len(allParts))
	seen := make(map[int]bool)

	for _, part := range allParts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("неверный аргумент: %s", part)
		}
		if seen[num] {
			return nil, fmt.Errorf("дубликат: %d", num)
		}
		seen[num] = true
		result = append(result, num)
	}

	return result, nil
}

// ParseInstructions читает инструкции из stdin, по одной на строку.
// Возвращает слайс инструкций и ошибку при невалидной инструкции.
func ParseInstructions() ([]string, error) {
	var instructions []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Пустая строка — конец ввода (чтобы checker не зависал при Enter)
		if line == "" {
			break
		}
		if !IsValidInstruction(line) {
			return nil, fmt.Errorf("неизвестная инструкция: %s", line)
		}
		instructions = append(instructions, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения: %w", err)
	}

	return instructions, nil
}

// Валидные инструкции push-swap.
var validInstructions = map[string]bool{
	"pa": true, "pb": true,
	"sa": true, "sb": true, "ss": true,
	"ra": true, "rb": true, "rr": true,
	"rra": true, "rrb": true, "rrr": true,
}

// IsValidInstruction проверяет, является ли строка валидной инструкцией.
func IsValidInstruction(s string) bool {
	return validInstructions[s]
}
