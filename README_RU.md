# Push-Swap

Две CLI-программы на Go для сортировки целых чисел через ограниченный набор операций над двумя стеками: `push-swap` генерирует последовательность команд, а `checker` выполняет её и проверяет результат.

Для маленьких стеков используются точные решения, для больших — LIS-based алгоритм с оптимизацией совместных вращений. Проект использует только стандартную библиотеку Go и содержит отдельный набор тестов, повторяющий executable-пункты школьного audit-листа.

· [English version](README.md)

## 📋 Содержание

- [🚀 Быстрый старт](#-быстрый-старт)
- [📝 О проекте](#-о-проекте)
- [✨ Возможности](#-возможности)
- [🔄 Архитектура](#-архитектура)
- [🎛️ Инструкции](#️-инструкции)
- [🧠 Алгоритм сортировки](#-алгоритм-сортировки)
- [✅ Checker](#-checker)
- [🧪 Тесты и audit](#-тесты-и-audit)
- [🧰 Технологии](#-технологии)
- [📁 Структура проекта](#-структура-проекта)
- [⚠️ Примечания](#️-примечания)
- [🧑‍💻 Авторы](#-авторы)

## 🚀 Быстрый старт

### Требования

- Go `1.21+`
- Git
- Bash / Git Bash или PowerShell
- GNU Make — опционально

### Клонирование

```bash
git clone https://github.com/legion2440/push-swap.git
cd push-swap
```

### Сборка

Linux / macOS / Git Bash с Make:

```bash
make all
```

Вручную:

```bash
go build -o checker ./cmd/checker
go build -o push-swap ./cmd/push-swap
```

Windows PowerShell:

```powershell
.\build.ps1
```

Скрипт создаёт `checker.exe` и `push-swap.exe`.

### Первый запуск

```bash
./push-swap "2 1 3 6 5 8"
```

Программа выводит только инструкции сортировки, по одной на строку.

Проверка результата:

```bash
ARG="4 67 3 87 23"
./push-swap "$ARG" | ./checker "$ARG"
```

Ожидаемый результат:

```text
OK
```

## 📝 О проекте

В проекте используются два стека:

- `A` содержит входные числа; первый integer находится наверху стека;
- `B` изначально пуст и используется как вспомогательный.

Цель — получить стек `A`, отсортированный по возрастанию сверху вниз, и пустой стек `B`, используя только разрешённые push-swap инструкции.

Репозиторий содержит две программы:

| Программа | Назначение |
| --- | --- |
| `push-swap` | генерирует последовательность инструкций для сортировки |
| `checker` | читает инструкции из stdin, выполняет их и печатает `OK` или `KO` |

Некорректные integer-аргументы и дубликаты приводят к `Error` в stderr. `checker` также возвращает `Error` для неизвестных или неправильно отформатированных инструкций.

Если аргументы не переданы, обе программы завершаются без вывода.

## ✨ Возможности

### Sorting

- все 11 разрешённых операций push-swap;
- корректная работа с отрицательными и произвольными `int` значениями;
- rank normalization для алгоритма сортировки;
- точная сортировка маленьких стеков;
- LIS-based сортировка больших стеков;
- объединённые вращения через `rr` и `rrr`;
- финальный поворот минимального элемента наверх.

### Validation

- проверка integer input;
- обнаружение дубликатов;
- `Error\n` в stderr для invalid input;
- строгая проверка instruction token;
- каждая инструкция должна быть завершена `\n`;
- неизвестные и malformed инструкции отклоняются.

### Verification

- unit-тесты для stack, actions, validation и algorithm;
- exhaustive permutations для размеров `2..6`;
- exhaustive `7! = 5040` permutations через large/LIS path;
- deterministic stress для 100 элементов;
- отдельный `tests/audit_test.go`, повторяющий executable-пункты школьного audit-листа;
- проверка, что product code использует только Go standard library.

## 🔄 Архитектура

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

Общая семантика stack operations находится в `internal/stack` и `internal/actions`, поэтому solver и checker используют одинаковое поведение инструкций.

## 🎛️ Инструкции

Проект поддерживает ровно 11 инструкций из задания.

| Инструкция | Действие |
| --- | --- |
| `pa` | перенести верхний элемент `B` в `A` |
| `pb` | перенести верхний элемент `A` в `B` |
| `sa` | поменять местами два верхних элемента `A` |
| `sb` | поменять местами два верхних элемента `B` |
| `ss` | выполнить `sa` и `sb` |
| `ra` | rotate `A`: верхний элемент становится последним |
| `rb` | rotate `B` |
| `rr` | выполнить `ra` и `rb` |
| `rra` | reverse rotate `A`: последний элемент становится верхним |
| `rrb` | reverse rotate `B` |
| `rrr` | выполнить `rra` и `rrb` |

Операция над стеком с недостаточным количеством элементов не меняет его состояние.

## 🧠 Алгоритм сортировки

### 2–3 элемента

Для двух и трёх элементов используются короткие прямые последовательности операций.

### 4–6 элементов

Используется breadth-first search по полному state space двух стеков.

Значения сначала преобразуются в относительные ранги, поэтому поиск зависит только от порядка чисел. Для шести элементов полный state space содержит всего:

```text
(6 + 1) × 6! = 5040 states
```

BFS гарантирует кратчайшую последовательность разрешённых инструкций для этих размеров.

### 7+ элементов

Для больших стеков используется LIS-based стратегия:

1. значения преобразуются в ранги;
2. вычисляется Longest Increasing Subsequence;
3. элементы LIS остаются в `A`, остальные отправляются в `B`;
4. для каждого элемента `B` рассчитывается cheapest insertion обратно в `A`;
5. совместимые вращения объединяются в `rr` или `rrr`;
6. после возврата всех элементов минимальный rank поворачивается наверх.

Large solver не заявляет математически глобальный minimum для произвольного размера; он оптимизирует количество операций и проходит performance limits из audit.

## ✅ Checker

`checker` получает тот же stack `A` через command-line arguments и читает инструкции из stdin.

Пример `KO` из audit:

```bash
echo -e "sa\npb\nrrr\n" | ./checker "0 9 1 8 2 7 3 6 4 5"
```

Результат:

```text
KO
```

Пример `OK`:

```bash
echo -e "pb\nra\npb\nra\nsa\nra\npa\npa\n" | ./checker "0 9 1 8 2"
```

Результат:

```text
OK
```

Instruction parsing строгий: token должен точно совпадать с одной из 11 команд и быть завершён newline. Пробелы внутри instruction line не нормализуются.

## 🧪 Тесты и audit

### Полный test suite

```bash
go test ./... -count=1
```

Статический анализ:

```bash
go vet ./...
```

Repo-local contracts:

```bash
python scripts/validate_agent_contracts.py
```

### Audit-facing suite

`tests/audit_test.go` зеркалит executable-проверки 01-edu audit-листа.

```bash
go test ./tests -run '^TestAudit_' -count=1 -v
```

Проверяются в том числе:

- `push-swap` без аргументов → без вывода;
- `2 1 3 6 5 8` → valid solution и `< 9` инструкций;
- уже отсортированный input → без вывода;
- non-integer → `Error`;
- duplicates → `Error`;
- два независимых 5-number cases → valid solution и `< 12`;
- checker `KO` example;
- checker `OK` example;
- `push-swap | checker` → `OK`;
- 100 unique numbers → `< 700` и `OK`;
- только Go standard-library dependencies.

### Дополнительное покрытие solver

Small solver проверяется exhaustively на всех permutations размеров `2..6`.

Large solver получает отдельные проверки:

- все `5040` permutations размера `7`;
- random cases разных размеров;
- `500` deterministic permutations по 100 элементов;
- structured 100-element cases: reverse, rotated и zigzag.

Stress test не просто считает операции: каждая полученная последовательность повторно исполняется через actions и должна оставить `A` отсортированным, а `B` пустым.

## 🧰 Технологии

| Область | Технология |
| --- | --- |
| Язык | Go `1.21` |
| Runtime dependencies | Go standard library only |
| Tests | `testing` |
| Build | Go toolchain, Make, PowerShell |
| CLI input | command-line arguments + stdin |
| Repository checks | Go tests, `go vet`, repo-local validator |

В `go.mod` нет сторонних dependencies.

## 📁 Структура проекта

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

`agent/` и `AGENTS.md` содержат repo-local navigation metadata и validation contracts. Они не участвуют в runtime двух CLI-программ.

## ⚠️ Примечания

- Product code использует только стандартные Go packages.
- Скомпилированные `checker`, `push-swap` и Windows `.exe` игнорируются Git.
- `push-swap` ничего не выводит для already-sorted input и при отсутствии аргументов.
- `checker` ничего не выводит при отсутствии аргументов.
- Ошибки печатаются как `Error\n` в stderr.
- Для маленьких стеков `4..6` exact BFS гарантирует shortest sequence; для `7+` используется оптимизированная LIS-based стратегия без заявления global optimality.
- Audit performance targets: `< 9` для `2 1 3 6 5 8`, `< 12` для tested 5-number cases и `< 700` для 100 unique numbers.

## 🧑‍💻 Авторы

1. Әсет Байырша ([@abaiyrsh](https://github.com/abaiyrsh))
2. Nazar Yestayev ([@nyestaye](https://github.com/nyestaye))
3. David Kortel ([@dkortel](https://github.com/dkortel))
