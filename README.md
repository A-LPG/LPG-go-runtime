# LPG-go-runtime

Go runtime for [LPG2](https://github.com/A-LPG/LPG2).

## Install / coordinates

| Field | Value |
|-------|-------|
| Package | Go module `github.com/A-LPG/LPG-go-runtime` |
| Version | git tags |
| Compatible generator | LPG2 ≥ 2.3.0 — see [`ecosystem/compat.json`](https://github.com/A-LPG/LPG2/blob/main/ecosystem/compat.json) |

```bash
go get github.com/A-LPG/LPG-go-runtime@latest
```

## Minimum toolchain

Go 1.17+.

## Build and test

```bash
go test ./...
```

## Wiring generated files

1. Generate with `-programming_language=go -table` and `dtParserTemplateF.gi`
2. Place generated packages in your module and import this runtime

## Features

| Feature | Status |
|---------|--------|
| Deterministic parser | yes |
| Backtracking | yes |
| Nested automatic AST | yes |
| `%Recover` prosthetic AST | yes |

## Publish status

- Channel: Go module proxy via git tags
- Automation: tag releases on this repository

## Links

- Generator: https://github.com/A-LPG/LPG2
- Ecosystem: https://github.com/A-LPG/LPG2/blob/main/docs/ECOSYSTEM.md
