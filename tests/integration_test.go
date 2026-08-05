// Интеграционные тесты по сценариям проверки (аудит).
// Запуск: go test ./tests -v -run Integration
package tests

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var (
	pushSwapBin  string
	checkerBin   string
	pushSwapPath string
	checkerPath  string
	testBinDir   string
	buildOnce    sync.Once
	buildErr     error
)

func init() {
	if runtime.GOOS == "windows" {
		pushSwapBin = "push-swap.exe"
		checkerBin = "checker.exe"
	} else {
		pushSwapBin = "push-swap"
		checkerBin = "checker"
	}
}

func runPushSwap(t *testing.T, args ...string) (stdout, stderr string, err error) {
	ensureBinaries(t)
	root := projectRoot()
	cmd := exec.Command(pushSwapPath, args...)
	cmd.Dir = root
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func runChecker(t *testing.T, args []string, stdin string) (stdout, stderr string, err error) {
	ensureBinaries(t)
	root := projectRoot()
	cmd := exec.Command(checkerPath, args...)
	cmd.Dir = root
	cmd.Stdin = strings.NewReader(stdin)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func projectRoot() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(f))
}

func ensureBinaries(t *testing.T) {
	t.Helper()
	buildOnce.Do(func() {
		root := projectRoot()
		var err error
		testBinDir, err = os.MkdirTemp("", "push-swap-integration-*")
		if err != nil {
			buildErr = fmt.Errorf("create temp binary dir: %w", err)
			return
		}

		pushSwapPath = filepath.Join(testBinDir, pushSwapBin)
		checkerPath = filepath.Join(testBinDir, checkerBin)

		for _, target := range []struct {
			output string
			pkg    string
		}{
			{pushSwapPath, "./cmd/push-swap"},
			{checkerPath, "./cmd/checker"},
		} {
			cmd := exec.Command("go", "build", "-o", target.output, target.pkg)
			cmd.Dir = root
			if out, err := cmd.CombinedOutput(); err != nil {
				buildErr = fmt.Errorf("build %s: %w\n%s", target.pkg, err, out)
				return
			}
		}
	})
	if buildErr != nil {
		t.Fatal(buildErr)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	if testBinDir != "" {
		_ = os.RemoveAll(testBinDir)
	}
	os.Exit(code)
}

func TestIntegration_PushSwap_NoArgs_OutputsNothing(t *testing.T) {
	ensureBinaries(t)
	stdout, stderr, _ := runPushSwap(t)
	if stdout != "" || stderr != "" {
		t.Errorf("push-swap без аргументов должен выводить ничего; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestIntegration_PushSwap_AlreadySorted_OutputsNothing(t *testing.T) {
	ensureBinaries(t)
	stdout, stderr, _ := runPushSwap(t, "0 1 2 3 4 5")
	if stdout != "" || stderr != "" {
		t.Errorf("push-swap для отсортированного должен выводить ничего; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestIntegration_PushSwap_InvalidArg_Error(t *testing.T) {
	ensureBinaries(t)
	_, stderr, _ := runPushSwap(t, "0 one 2 3")
	if !strings.Contains(stderr, "Error") {
		t.Errorf("push-swap '0 one 2 3' должен вывести Error; stderr=%q", stderr)
	}
}

func TestIntegration_PushSwap_Duplicates_Error(t *testing.T) {
	ensureBinaries(t)
	_, stderr, _ := runPushSwap(t, "1 2 2 3")
	if !strings.Contains(stderr, "Error") {
		t.Errorf("push-swap '1 2 2 3' должен вывести Error; stderr=%q", stderr)
	}
}

func TestIntegration_PushSwap_Example_LessThan9Instructions(t *testing.T) {
	ensureBinaries(t)
	stdout, stderr, _ := runPushSwap(t, "2 1 3 6 5 8")
	if stderr != "" {
		t.Errorf("push-swap '2 1 3 6 5 8' не должен выводить в stderr; stderr=%q", stderr)
	}
	lines := countLines(stdout)
	if lines >= 9 {
		t.Errorf("push-swap '2 1 3 6 5 8' должен выдать < 9 инструкций; получили %d", lines)
	}
	runAndCheck(t, "2 1 3 6 5 8", stdout, "OK")
}

func TestIntegration_PushSwap_FiveNumbers_LessThan12Instructions(t *testing.T) {
	ensureBinaries(t)
	arg := "4 67 3 87 23"
	stdout, stderr, _ := runPushSwap(t, arg)
	if stderr != "" {
		t.Errorf("push-swap %q не должен выводить в stderr; stderr=%q", arg, stderr)
	}
	lines := countLines(stdout)
	if lines >= 12 {
		t.Errorf("push-swap для 5 чисел должен выдать < 12 инструкций; получили %d", lines)
	}
	runAndCheck(t, arg, stdout, "OK")
}

func TestIntegration_PushSwap_GeneratesValidSolution(t *testing.T) {
	ensureBinaries(t)
	tests := []string{
		"2 1",
		"3 2 1",
		"5 4 3 2 1",
		"2 1 3 6 5 8",
		"4 67 3 87 23",
	}
	for _, arg := range tests {
		stdout, stderr, _ := runPushSwap(t, arg)
		if stderr != "" {
			t.Errorf("push-swap %q: stderr=%q", arg, stderr)
			continue
		}
		runAndCheck(t, arg, stdout, "OK")
	}
}

func runAndCheck(t *testing.T, arg, instructions, want string) {
	t.Helper()
	stdout, stderr, _ := runChecker(t, []string{arg}, instructions)
	if stderr != "" {
		t.Errorf("checker %q: неожиданный stderr=%q", arg, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != want {
		t.Errorf("checker %q с инструкциями: получили %q, ожидалось %q", arg, got, want)
	}
}

func TestIntegration_Checker_NoArgs_OutputsNothing(t *testing.T) {
	ensureBinaries(t)
	stdout, stderr, _ := runChecker(t, nil, "")
	if stdout != "" || stderr != "" {
		t.Errorf("checker без аргументов должен выводить ничего; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestIntegration_Checker_InvalidArg_Error(t *testing.T) {
	ensureBinaries(t)
	_, stderr, _ := runChecker(t, []string{"0 one 2 3"}, "sa\n")
	if !strings.Contains(stderr, "Error") {
		t.Errorf("checker '0 one 2 3' должен вывести Error; stderr=%q", stderr)
	}
}

// Неизвестная инструкция должна давать Error
func TestIntegration_Checker_InvalidInstruction_Error(t *testing.T) {
	ensureBinaries(t)
	_, stderr, _ := runChecker(t, []string{"1 2 3"}, "sa\nxx\npa\n")
	if !strings.Contains(stderr, "Error") {
		t.Errorf("checker с неизвестной инструкцией должен вывести Error; stderr=%q", stderr)
	}
}

func TestIntegration_Checker_KO(t *testing.T) {
	ensureBinaries(t)
	stdout, stderr, _ := runChecker(t, []string{"0 9 1 8 2 7 3 6 4 5"}, "sa\npb\nrrr\n")
	if stderr != "" {
		t.Errorf("checker: stderr=%q", stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "KO" {
		t.Errorf("checker с sa,pb,rrr: получили %q, ожидалось KO", got)
	}
}

func TestIntegration_Checker_OK(t *testing.T) {
	ensureBinaries(t)
	stdout, stderr, _ := runChecker(t, []string{"0 9 1 8 2"}, "pb\nra\npb\nra\nsa\nra\npa\npa\n")
	if stderr != "" {
		t.Errorf("checker: stderr=%q", stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "OK" {
		t.Errorf("checker: получили %q, ожидалось OK", got)
	}
}

func TestIntegration_PushSwapPipeChecker_OK(t *testing.T) {
	ensureBinaries(t)
	arg := "4 67 3 87 23"
	stdout, _, _ := runPushSwap(t, arg)
	runAndCheck(t, arg, stdout, "OK")
}

// Пример из задания: rra pb sa rra pa сортирует "3 2 1 0"
func TestIntegration_Checker_ExampleFromTask_OK(t *testing.T) {
	ensureBinaries(t)
	stdout, stderr, _ := runChecker(t, []string{"3 2 1 0"}, "rra\npb\nsa\nrra\npa\n")
	if stderr != "" {
		t.Errorf("checker: stderr=%q", stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "OK" {
		t.Errorf("checker '3 2 1 0' с rra,pb,sa,rra,pa: получили %q, ожидалось OK", got)
	}
}

// Проверка: push-swap "2 1 3 6 5 8" даёт валидную последовательность из примера
func TestIntegration_PushSwap_ExampleSequence(t *testing.T) {
	ensureBinaries(t)
	stdout, _, _ := runPushSwap(t, "2 1 3 6 5 8")
	runAndCheck(t, "2 1 3 6 5 8", stdout, "OK")
}

func countLines(s string) int {
	n := 0
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		if strings.TrimSpace(sc.Text()) != "" {
			n++
		}
	}
	return n
}
