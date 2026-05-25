// Тесты для пакета validation.
package validation

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		args     []string
		wantErr  bool
		expected []int
	}{
		{[]string{"1", "2", "3"}, false, []int{1, 2, 3}},
		{[]string{"2 1 3"}, false, []int{2, 1, 3}},
		{[]string{}, false, nil},
		{[]string{"1", "one"}, true, nil},
		{[]string{"1", "2", "2"}, true, nil},
		{[]string{"-5", "0", "100"}, false, []int{-5, 0, 100}},
	}
	for _, tt := range tests {
		got, err := ParseArgs(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseArgs(%v): err = %v, wantErr = %v", tt.args, err, tt.wantErr)
			continue
		}
		if !tt.wantErr {
			if len(got) != len(tt.expected) {
				t.Errorf("ParseArgs(%v): len = %d, ожидалось %d", tt.args, len(got), len(tt.expected))
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("ParseArgs(%v): [%d] = %d, ожидалось %d", tt.args, i, got[i], tt.expected[i])
				}
			}
		}
	}
}

func TestParseArgs_EdgeCases(t *testing.T) {
	tests := []struct {
		args     []string
		wantErr  bool
		expected []int
	}{
		{[]string{"0"}, false, []int{0}},
		{[]string{"-1", "-2", "-3"}, false, []int{-1, -2, -3}},
		{[]string{"  5  4  3  "}, false, []int{5, 4, 3}},
		{[]string{"1", "2", "3", "4", "5"}, false, []int{1, 2, 3, 4, 5}},
		{[]string{"abc"}, true, nil},
		{[]string{"1", "2", "1"}, true, nil},
		{[]string{"2147483647"}, false, []int{2147483647}},
	}
	for _, tt := range tests {
		got, err := ParseArgs(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseArgs(%v): err = %v, wantErr = %v", tt.args, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && len(got) == len(tt.expected) {
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("ParseArgs(%v): [%d] = %d, ожидалось %d", tt.args, i, got[i], tt.expected[i])
				}
			}
		}
	}
}

func TestIsValidInstruction(t *testing.T) {
	valid := []string{"pa", "pb", "sa", "sb", "ss", "ra", "rb", "rr", "rra", "rrb", "rrr"}
	for _, instr := range valid {
		if !IsValidInstruction(instr) {
			t.Errorf("IsValidInstruction(%q) = false, ожидалось true", instr)
		}
	}
	invalid := []string{"", "xx", "paa", "rraa", "p"}
	for _, instr := range invalid {
		if IsValidInstruction(instr) {
			t.Errorf("IsValidInstruction(%q) = true, ожидалось false", instr)
		}
	}
}
