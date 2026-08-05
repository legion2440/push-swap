# Agent navigation

1. Read `agent/module-index.json`.
2. Identify the module that owns the requested feature using roots and entrypoints.
3. Read only `agent/modules/<module-id>.json` for local work.
4. Inspect production code/config before making implementation claims.
5. Read `agent/dependency-graph.json` before changing cross-module contracts.
6. Treat README and agent metadata as navigation, not implementation evidence.
7. Run affected Go tests after production changes.
8. Run `python scripts/validate_agent_contracts.py` after agent-metadata changes.
9. Do not read every manifest unless the task is repository-wide.
10. If ownership is ambiguous, inspect candidate roots and the dependency graph instead of guessing.

## Project guardrails

- Canonical subject: `https://github.com/01-edu/public/tree/master/subjects/push-swap`.
- Stay within the subject and audit scope.
- Implementation language: Go.
- Product code may use only the Go standard library.
- Required executables: `checker` and `push-swap`.
- Audit performance targets: `< 9` instructions for `2 1 3 6 5 8`, `< 12` for any tested 5-number case, and `< 700` for 100 random unique numbers.
