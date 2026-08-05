package tests

import (
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// This file mirrors the executable checks from the 01-edu push-swap audit.
// Run only the audit-facing suite with:
//
//	go test ./tests -run '^TestAudit_' -count=1 -v

func TestAudit_AllowedPackagesOnly(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", "-f", "{{if not .Standard}}{{.ImportPath}}{{end}}", "./cmd/...")
	cmd.Dir = projectRoot()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go list failed: %v\n%s", err, out)
	}

	for _, pkg := range strings.Fields(string(out)) {
		if strings.HasPrefix(pkg, "push-swap/") {
			continue
		}
		t.Errorf("non-standard dependency found: %s", pkg)
	}
}

func TestAudit_PushSwap_NoArgs(t *testing.T) {
	stdout, stderr, err := runPushSwap(t)
	if err != nil {
		t.Fatalf("push-swap failed: %v", err)
	}
	if stdout != "" || stderr != "" {
		t.Fatalf("expected no output; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestAudit_PushSwap_ExampleUnder9(t *testing.T) {
	arg := "2 1 3 6 5 8"
	stdout, stderr, err := runPushSwap(t, arg)
	if err != nil {
		t.Fatalf("push-swap failed: %v", err)
	}
	if stderr != "" {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
	if got := countLines(stdout); got >= 9 {
		t.Fatalf("got %d instructions, audit requires < 9", got)
	}
	assertCheckerResult(t, arg, stdout, "OK")
}

func TestAudit_PushSwap_AlreadySorted(t *testing.T) {
	stdout, stderr, err := runPushSwap(t, "0 1 2 3 4 5")
	if err != nil {
		t.Fatalf("push-swap failed: %v", err)
	}
	if stdout != "" || stderr != "" {
		t.Fatalf("expected no output; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestAudit_PushSwap_NonIntegerError(t *testing.T) {
	stdout, stderr, _ := runPushSwap(t, "0 one 2 3")
	if stdout != "" || stderr != "Error\n" {
		t.Fatalf("stdout=%q stderr=%q, want stderr Error\\n only", stdout, stderr)
	}
}

func TestAudit_PushSwap_DuplicateError(t *testing.T) {
	stdout, stderr, _ := runPushSwap(t, "1 2 2 3")
	if stdout != "" || stderr != "Error\n" {
		t.Fatalf("stdout=%q stderr=%q, want stderr Error\\n only", stdout, stderr)
	}
}

func TestAudit_PushSwap_FiveRandomNumbers_First(t *testing.T) {
	assertFiveNumberAuditCase(t, "4 67 3 87 23")
}

func TestAudit_PushSwap_FiveRandomNumbers_Second(t *testing.T) {
	assertFiveNumberAuditCase(t, "42 -7 19 0 5")
}

func TestAudit_Checker_NoArgs(t *testing.T) {
	stdout, stderr, err := runChecker(t, nil, "")
	if err != nil {
		t.Fatalf("checker failed: %v", err)
	}
	if stdout != "" || stderr != "" {
		t.Fatalf("expected no output; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestAudit_Checker_NonIntegerError(t *testing.T) {
	stdout, stderr, _ := runChecker(t, []string{"0 one 2 3"}, "")
	if stdout != "" || stderr != "Error\n" {
		t.Fatalf("stdout=%q stderr=%q, want stderr Error\\n only", stdout, stderr)
	}
}

func TestAudit_Checker_KO(t *testing.T) {
	stdin := "sa\npb\nrrr\n\n"
	stdout, stderr, err := runChecker(t, []string{"0 9 1 8 2 7 3 6 4 5"}, stdin)
	if err != nil {
		t.Fatalf("checker failed: %v", err)
	}
	if stdout != "KO\n" || stderr != "" {
		t.Fatalf("stdout=%q stderr=%q, want KO\\n", stdout, stderr)
	}
}

func TestAudit_Checker_OK(t *testing.T) {
	stdin := "pb\nra\npb\nra\nsa\nra\npa\npa\n\n"
	stdout, stderr, err := runChecker(t, []string{"0 9 1 8 2"}, stdin)
	if err != nil {
		t.Fatalf("checker failed: %v", err)
	}
	if stdout != "OK\n" || stderr != "" {
		t.Fatalf("stdout=%q stderr=%q, want OK\\n", stdout, stderr)
	}
}

func TestAudit_PushSwapPipeChecker_OK(t *testing.T) {
	arg := "4 67 3 87 23"
	instructions, stderr, err := runPushSwap(t, arg)
	if err != nil {
		t.Fatalf("push-swap failed: %v", err)
	}
	if stderr != "" {
		t.Fatalf("unexpected push-swap stderr: %q", stderr)
	}
	assertCheckerResult(t, arg, instructions, "OK")
}

func TestAudit_100RandomNumbers_Under700AndOK(t *testing.T) {
	rng := rand.New(rand.NewSource(20260805))
	values := rng.Perm(100)
	arg := joinInts(values)

	instructions, stderr, err := runPushSwap(t, arg)
	if err != nil {
		t.Fatalf("push-swap failed: %v", err)
	}
	if stderr != "" {
		t.Fatalf("unexpected push-swap stderr: %q", stderr)
	}
	if got := countLines(instructions); got >= 700 {
		t.Fatalf("got %d commands, audit requires < 700", got)
	} else {
		t.Logf("100-number audit case: %d commands", got)
	}
	assertCheckerResult(t, arg, instructions, "OK")
}

func assertFiveNumberAuditCase(t *testing.T, arg string) {
	t.Helper()
	instructions, stderr, err := runPushSwap(t, arg)
	if err != nil {
		t.Fatalf("push-swap failed: %v", err)
	}
	if stderr != "" {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
	if got := countLines(instructions); got >= 12 {
		t.Fatalf("got %d instructions, audit requires < 12", got)
	}
	assertCheckerResult(t, arg, instructions, "OK")
}

func assertCheckerResult(t *testing.T, arg, instructions, want string) {
	t.Helper()
	stdout, stderr, err := runChecker(t, []string{arg}, instructions)
	if err != nil {
		t.Fatalf("checker failed: %v", err)
	}
	if stderr != "" {
		t.Fatalf("unexpected checker stderr: %q", stderr)
	}
	if stdout != want+"\n" {
		t.Fatalf("checker output=%q, want %q", stdout, want+"\n")
	}
}

func joinInts(values []int) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = strconv.Itoa(value)
	}
	return strings.Join(parts, " ")
}
