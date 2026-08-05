# Push-Swap

Two Go CLI programs for sorting integers through a constrained set of operations on two stacks: `push-swap` generates an instruction sequence, while `checker` executes that sequence and validates the result.

Small stacks use exact solutions, larger stacks use an LIS-based strategy with combined-rotation optimization, and the product code depends only on the Go standard library. The repository also contains an audit-facing test suite that mirrors the executable checks from the school rubric.

· [Русская версия](README_RU.md)

## 📋 TOC

- [🚀 Quick start](#-quick-start)
- [📝 About](#-about)
- [✨ Features](#-features)
- [🔄 Architecture](#-architecture)
- [🎛️ Instructions](#️-instructions)
- [🧠 Sorting algorithm](#-sorting-algorithm)
- [✅ Checker](#-checker)
- [🧪 Tests and audit](#-tests-and-audit)
- [🧰 Technology stack](#-technology-stack)
- [📁 Project structure](#-project-structure)
- [⚠️ Notes](#️-notes)
- [🧑‍💻 Authors](#-authors)

## 🚀 Quick start

### Prerequisites

- Go `1.21+`
- Git
- Bash / Git Bash or PowerShell
- GNU Make is optional

### Clone

```bash
git clone https://01.tomorrow-school.ai/git/nyestaye/push-swap
cd push-swap
```

### Build

Linux / macOS / Git Bash with Make:

```bash
make all
```

Manual build:

```bash
go build -o checker ./cmd/checker
go build -o push-swap ./cmd/push-swap
```

Windows PowerShell:

```powershell
.\build.ps1
```

The script creates `checker.exe` and `push-swap.exe`.

### First run

```bash
./push-swap "2 1 3 6 5 8"
```

The program prints only sorting instructions, one per line.

Validate a generated sequence:

```bash
ARG="4 67 3 87 23"
./push-swap "$ARG" | ./checker "$ARG"
```

Expected result:

```text
OK
```

## 📝 About

The project operates on two stacks:

- `A` contains the input integers; the first integer is at the top of the stack;
- `B` starts empty and is used as auxiliary storage.

The goal is to leave `A` sorted in ascending order from top to bottom and `B` empty while using only the allowed push-swap instructions.

The repository builds two programs:

| Program | Responsibility |
| --- | --- |
| `push-swap` | generate an instruction sequence that sorts the input stack |
| `checker` | read instructions from stdin, execute them, and print `OK` or `KO` |

Invalid integer arguments and duplicates produce `Error` on stderr. `checker` also reports `Error` for unknown or incorrectly formatted instructions.

When no arguments are supplied, both programs exit without output.

## ✨ Features

### Sorting

- all 11 allowed push-swap operations;
- negative and arbitrary `int` values;
- rank normalization for order-only reasoning;
- exact small-stack sorting;
- LIS-based large-stack sorting;
- combined rotations through `rr` and `rrr`;
- final rotation that places the minimum element on top.

### Validation

- integer input parsing;
- duplicate detection;
- `Error\n` on stderr for invalid input;
- strict instruction-token validation;
- every instruction must be newline-terminated;
- unknown and malformed instructions are rejected.

### Verification

- unit tests for stack, actions, validation, and algorithm modules;
- exhaustive permutations for sizes `2..6`;
- exhaustive `7! = 5040` permutations through the large/LIS path;
- deterministic 100-element stress tests;
- dedicated `tests/audit_test.go` mirroring the executable school audit checks;
- dependency validation confirming product code uses only the Go standard library.

## 🔄 Architecture

```text
                    command-line integers
                             |
                             v
                    +------------------+
                    |    validation    |
                    +---------+--------+
                              |
                              v
                         +---------+
                         | stack A |
                         +---------+

        push-swap                           checker
            |                                  |
            v                                  v
 +----------------------+          +-----------------------+
 | sorting algorithm    |          | stdin instructions    |
 | 2-3 direct           |          | strict validation     |
 | 4-6 exact BFS        |          +-----------+-----------+
 | 7+ LIS + insertion   |                      |
 +----------+-----------+                      v
            |                         +--------------------+
            v                         | execute operations |
 +----------------------+             | on stacks A and B  |
 | instruction sequence |             +---------+----------+
 +----------+-----------+                       |
            |                                   v
            v                               OK / KO
          stdout
```

Shared stack semantics live in `internal/stack` and `internal/actions`, so the solver and checker execute the same operation behavior.

## 🎛️ Instructions

The project supports exactly the 11 instructions defined by the assignment.

| Instruction | Action |
| --- | --- |
| `pa` | push the top element of `B` to `A` |
| `pb` | push the top element of `A` to `B` |
| `sa` | swap the first two elements of `A` |
| `sb` | swap the first two elements of `B` |
| `ss` | execute `sa` and `sb` |
| `ra` | rotate `A`: the first element becomes the last |
| `rb` | rotate `B` |
| `rr` | execute `ra` and `rb` |
| `rra` | reverse rotate `A`: the last element becomes the first |
| `rrb` | reverse rotate `B` |
| `rrr` | execute `rra` and `rrb` |

An operation on a stack with too few elements leaves that stack unchanged.

## 🧠 Sorting algorithm

### 2–3 elements

Two- and three-element inputs use short direct operation sequences.

### 4–6 elements

The solver uses breadth-first search over the complete two-stack state space.

Values are first converted to relative ranks, so the search depends only on ordering. For six values the complete state space contains only:

```text
(6 + 1) × 6! = 5040 states
```

BFS therefore guarantees a shortest valid instruction sequence for these sizes.

### 7+ elements

Larger inputs use an LIS-based strategy:

1. convert values to ranks;
2. compute a Longest Increasing Subsequence;
3. keep LIS elements in `A` and push the others to `B`;
4. calculate the cheapest insertion for each element moved back from `B`;
5. combine compatible rotations with `rr` or `rrr`;
6. rotate the minimum rank to the top after all elements return to `A`.

The large solver does not claim a mathematically global minimum for arbitrary input sizes. It optimizes the generated sequence and is validated against the assignment's performance limits.

## ✅ Checker

`checker` receives the same stack `A` through command-line arguments and reads instructions from stdin.

Audit `KO` example:

```bash
echo -e "sa\npb\nrrr\n" | ./checker "0 9 1 8 2 7 3 6 4 5"
```

Result:

```text
KO
```

Audit `OK` example:

```bash
echo -e "pb\nra\npb\nra\nsa\nra\npa\npa\n" | ./checker "0 9 1 8 2"
```

Result:

```text
OK
```

Instruction parsing is strict: each token must exactly match one of the 11 commands and must be terminated by a newline. Whitespace inside an instruction line is not normalized.

## 🧪 Tests and audit

### Full test suite

```bash
go test ./... -count=1
```

Static analysis:

```bash
go vet ./...
```

Repo-local contracts:

```bash
python scripts/validate_agent_contracts.py
```

### Audit-facing suite

`tests/audit_test.go` mirrors the executable checks from the 01-edu audit list.

```bash
go test ./tests -run '^TestAudit_' -count=1 -v
```

The suite checks, among other cases:

- `push-swap` with no arguments → no output;
- `2 1 3 6 5 8` → valid solution and `< 9` instructions;
- already sorted input → no output;
- non-integer input → `Error`;
- duplicates → `Error`;
- two independent 5-number cases → valid solution and `< 12` instructions;
- checker `KO` example;
- checker `OK` example;
- `push-swap | checker` → `OK`;
- 100 unique numbers → `< 700` instructions and `OK`;
- standard-library-only product dependencies.

### Additional solver coverage

The small solver is checked exhaustively across every permutation of sizes `2..6`.

The large solver additionally covers:

- all `5040` permutations of size `7`;
- random cases across several larger sizes;
- `500` deterministic permutations of 100 elements;
- structured 100-element cases: reverse, rotated, and zigzag.

Stress tests do not only count operations: every generated sequence is executed again through the operation layer and must leave `A` sorted and `B` empty.

## 🧰 Technology stack

| Area | Technology |
| --- | --- |
| Language | Go `1.21` |
| Runtime dependencies | Go standard library only |
| Tests | Go `testing` |
| Build | Go toolchain, Make, PowerShell |
| CLI input | command-line arguments + stdin |
| Repository checks | Go tests, `go vet`, repo-local validator |

`go.mod` contains no third-party dependencies.

## 📁 Project structure

```text
push-swap/
├── agent/
│   ├── modules/
│   ├── schemas/
│   ├── dependency-graph.json
│   └── module-index.json
├── cmd/
│   ├── checker/
│   │   └── main.go
│   └── push-swap/
│       └── main.go
├── internal/
│   ├── actions/
│   │   ├── actions.go
│   │   └── actions_test.go
│   ├── algorithm/
│   │   ├── bench_test.go
│   │   ├── sort.go
│   │   └── sort_test.go
│   ├── stack/
│   │   ├── stack.go
│   │   └── stack_test.go
│   └── validation/
│       ├── validation.go
│       └── validation_test.go
├── scripts/
│   └── validate_agent_contracts.py
├── tests/
│   ├── audit_test.go
│   └── integration_test.go
├── .gitignore
├── AGENTS.md
├── build.ps1
├── go.mod
├── Makefile
├── README.md
└── README_RU.md
```

`agent/` and `AGENTS.md` contain repo-local navigation metadata and validation contracts. They are not part of the runtime path of either CLI program.

## ⚠️ Notes

- Product code uses only standard Go packages.
- Built `checker`, `push-swap`, and Windows `.exe` binaries are ignored by Git.
- `push-swap` prints nothing for already sorted input and when no arguments are supplied.
- `checker` prints nothing when no arguments are supplied.
- Errors are written as `Error\n` to stderr.
- Exact BFS guarantees shortest sequences for sizes `4..6`; the `7+` solver uses an optimized LIS-based strategy without claiming global optimality.
- Audit performance targets are `< 9` for `2 1 3 6 5 8`, `< 12` for tested 5-number cases, and `< 700` for 100 unique numbers.

## 🧑‍💻 Authors

1. Әсет Байырша ([@abaiyrsh](https://github.com/abaiyrsh))
2. Nazar Yestayev ([@nyestaye](https://github.com/nyestaye))
3. David Kortel ([@dkortel](https://github.com/dkortel))
