# `services/studio`

Small, testable helpers for the **Studio** content pipeline. Business rules and HTTP handlers live in **`business/v1`** and **`api/v1`**; this package holds pure utilities.

## Contents

| File | Purpose |
|------|---------|
| `fractional_index.go` | **`GeneratePositionBetween(prev, next string)`** — produces lexicographically ordered keys between neighbours for Piece ordering within a phase (fractional indexing). |
| `fractional_index_test.go` | Unit tests for index generation edge cases |
| `phase_template.go` | Default phase scaffolding when a new studio is created |

## References

- Pipeline overview: [`docs/studio-pipeline.md`](../../docs/studio-pipeline.md)
- Pieces entity: [`entities/piece.go`](../../entities/piece.go) (`position` column)
