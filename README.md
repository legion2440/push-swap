# Push-Swap

Two Go CLI programs for sorting integers with a restricted set of operations on two stacks. `push-swap` generates the operations; `checker` replays them and verifies the result.

The solver uses exact search for small inputs and an LIS-based strategy for larger ones. Product code uses only the Go standard library.

· [Русская версия](README_RU.md)  
· [School repository](https://01.tomorrow-school.ai/git/nyestaye/push-swap)

## 📋 TOC

- [🚀 Quick start](#-quick-start)
- [📝 About](#-about)
- [🎛️ Operations](#️-operations)
- [🧠 Sorting strategy](#-sorting-strategy)
- [✅ Checker](#-checker)
- [🧪 Tests and audit](#-tests-and-audit)
- [📁 Project structure](#-project-structure)
- [🧑‍💻 Authors](#-authors)

## 🚀 Quick start

### Requirements

- Go `1.21+`
- Git
- Bash; Git Bash is supported on Windows

### Clone

```bash
git clone https://01.tomorrow-school.ai/git/nyestaye/push-swap
cd push-swap
```

### Build

The same commands work in Bash on Linux/macOS and in Git Bash on Windows:

```bash
go build ./cmd/checker
go build ./cmd/push-swap
```

Go creates `checker` and `push-swap` on Linux/macOS and `checker.exe` and `push-swap.exe` on Windows.

### Run

Examples below use Linux/macOS binary names. In Git Bash on Windows, use the same commands with `.exe`.

```bash
./push-swap "2 1 3 6 5 8"
```

Validate the generated sequence:

```bash
ARG="4 67 3 87 23"
./push-swap "$ARG" | ./checker "$ARG"
```

Expected result:

```text
OK
```

## 📝 About

The assignment operates on two stacks:

- `A` starts with the input integers; the first integer is on top;
- `B` starts empty;
- the target state is `A` sorted in ascending order and `B` empty.

`push-swap` may print only the allowed operations, one per line. `checker` reads those operations from stdin and prints `OK` when the final state is correct, otherwise `KO`.

Invalid integers, duplicates, unknown operations, and malformed checker input produce `Error` on stderr. With no command-line arguments, both programs print nothing.

## 🎛️ Operations

The assignment defines exactly 11 operations:

- `pa`, `pb` — push the top element from one stack to the other;
- `sa`, `sb`, `ss` — swap the first two elements of one or both stacks;
- `ra`, `rb`, `rr` — rotate one or both stacks forward;
- `rra`, `rrb`, `rrr` — rotate one or both stacks backwards.

An operation that cannot affect a stack because it contains too few elements leaves that stack unchanged.

## 🧠 Sorting strategy

### 2-3 elements

Small direct cases use short hand-written sequences.

### 4-6 elements

The solver runs breadth-first search over the complete two-stack state space after converting values to relative ranks. For six values the full state space contains only `7 × 6! = 5040` states, so BFS returns a shortest valid sequence for these sizes.

### 7+ elements

Larger inputs use an LIS-based strategy:

1. normalize values to ranks;
2. compute a Longest Increasing Subsequence;
3. keep LIS elements in `A` and push the others to `B`;
4. return elements from `B` using the cheapest insertion position;
5. combine compatible rotations with `rr` and `rrr`;
6. rotate the minimum element to the top.

This path is optimized for the assignment limits; it does not claim a globally minimal sequence for arbitrary input sizes.

## ✅ Checker

Known `KO` case from the assignment audit:

```bash
echo -e "sa\npb\nrrr\n" | ./checker "0 9 1 8 2 7 3 6 4 5"
```

```text
KO
```

Known `OK` case:

```bash
echo -e "pb\nra\npb\nra\nsa\nra\npa\npa\n" | ./checker "0 9 1 8 2"
```

```text
OK
```

Checker instruction parsing is strict: each operation must match one of the 11 allowed tokens and be newline-terminated.

## 🧪 Tests and audit

Run the complete test suite:

```bash
go test ./... -count=1
```

Static checks:

```bash
go vet ./...
gofmt -l $(find . -name '*.go' -type f)
```

The executable cases from the 01-edu audit list are collected in `tests/audit_test.go` and can be run separately:

```bash
go test ./tests -run '^TestAudit_' -count=1 -v
```

That suite covers the official no-output and error cases, the checker `OK` / `KO` examples, `push-swap | checker`, the `< 9` six-number case, `< 12` five-number cases, the `< 700` 100-number limit, and the standard-library-only requirement.

Additional solver coverage includes:

- every permutation of sizes `2..6`;
- all `5040` permutations of size `7` through the large-solver path;
- varied larger sizes;
- `500` deterministic 100-element permutations, with every generated sequence replayed for correctness;
- structured reverse, rotated, and zigzag 100-element inputs.

## 📁 Project structure

```text
push-swap/
├── cmd/
│   ├── checker/
│   └── push-swap/
├── internal/
│   ├── actions/
│   ├── algorithm/
│   ├── stack/
│   └── validation/
├── tests/
│   ├── audit_test.go
│   └── integration_test.go
├── agent/
├── scripts/
├── go.mod
├── Makefile
├── README.md
└── README_RU.md
```

`cmd/` contains the two executables. Shared stack and operation semantics live under `internal/`; tests exercise both the internal logic and the compiled CLI programs.

## 🧑‍💻 Authors

1. Әсет Байырша ([@abaiyrsh](https://github.com/abaiyrsh))
2. Nazar Yestayev ([@nyestaye](https://github.com/nyestaye))
3. David Kortel ([@dkortel](https://github.com/dkortel))
