// Тесты для пакета validation.
package validation

import (
	"reflect"
	"strings"
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
		{[]string{""}, true, nil},
		{[]string{"   "}, true, nil},
		{[]string{"1", ""}, true, nil},
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
	invalid := []string{"", "xx", "paa", "rraa", "p", " sa", "sa ", "sa\t"}
	for _, instr := range invalid {
		if IsValidInstruction(instr) {
			t.Errorf("IsValidInstruction(%q) = true, ожидалось false", instr)
		}
	}
}

func TestParseInstructions_StrictFormatting(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{name: "valid", input: "sa\npb\nrrr\n", want: []string{"sa", "pb", "rrr"}},
		{name: "echo trailing blank", input: "sa\npb\nrrr\n\n", want: []string{"sa", "pb", "rrr"}},
		{name: "multiple trailing blanks", input: "sa\n\n\n", want: []string{"sa"}},
		{name: "blank-only input", input: "\n", want: nil},
		{name: "leading space", input: " sa\n", wantErr: true},
		{name: "trailing space", input: "sa \n", wantErr: true},
		{name: "tab suffix", input: "sa\t\n", wantErr: true},
		{name: "blank line between instructions", input: "sa\n\npb\n", wantErr: true},
		{name: "unknown instruction", input: "sa\nxx\n", wantErr: true},
		{name: "missing final newline", input: "sa", wantErr: true},
		{name: "missing final newline after valid prefix", input: "sa\npb", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInstructions(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseInstructions(%q): err=%v, wantErr=%v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseInstructions(%q)=%v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
