# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go SDK for the S.EE URL shortener service (`github.com/sdotee/sdk.go`). Package name is `seesdk`. Zero external dependencies — uses only the Go standard library. Requires Go 1.21+.

## Commands

```bash
go build ./...                              # Build
go test ./...                               # Run all tests (requires SEE_API_KEY)
go test -run TestNewClient ./...            # Run a single test
go test -v ./...                            # Verbose test output
```

## Testing

All tests are **integration tests** that hit the live S.EE API. They require environment variables:

- `SEE_API_KEY` (required) — tests skip gracefully if not set
- `SEE_BASE_URL` (optional) — defaults to `https://s.ee/api/v1`

Tests use `setupTestClient(t)` as a shared helper that reads env vars and calls `t.Skip()` when credentials are missing.

## Architecture

Four source files, single package (`seesdk`), flat structure:

- **`client.go`** — `Client` struct, `Config`, `NewClient()` constructor, and private HTTP helpers (`doRequest` for JSON, `doMultipartRequest` for file uploads). Auth is done by setting the `Authorization` header directly (not Bearer prefix).
- **`api.go`** — All 13 public methods on `Client` (CRUD for short URLs, text, files; read-only for domains, tags, usage). Uses a generic `unmarshalResponse[T]` helper for JSON deserialization.
- **`models.go`** — Request/response structs with JSON tags. Responses follow a common `{code, data, message}` envelope pattern with nested `Data` structs.
- **`client_test.go`** — Integration tests covering the full API surface.

## API Endpoints Mapped to Methods

| Method | HTTP | Endpoint |
|---|---|---|
| `CreateShortURL` | POST | `/shorten` |
| `UpdateShortURL` | PUT | `/shorten` |
| `DeleteShortURL` | DELETE | `/shorten` |
| `CreateText` | POST | `/text` |
| `UpdateText` | PUT | `/text` |
| `DeleteText` | DELETE | `/text` |
| `UploadFile` | POST (multipart) | `/file/upload` |
| `DeleteFile` | GET | `/file/delete/:key` |
| `GetDomains` | GET | `/domains` |
| `GetFileDomains` | GET | `/file/domains` |
| `GetTextDomains` | GET | `/text/domains` |
| `GetTags` | GET | `/tags` |
| `GetUsage` | GET | `/usage` |

## Conventions

- File headers include copyright, author, and timestamp comments
- Request structs use `json:"field_name,omitempty"` tags; response structs use inline anonymous `Data` structs
- All public methods return `(*ResponseType, error)`
- Errors are wrapped with `fmt.Errorf("context: %w", err)`
- File uploads are capped at 100MB (`checkFileSize` validates via `Stat()` or `Len()` interface)
- `UsageNoLimit = -1` sentinel for unlimited usage quotas

## EditorConfig

4-space indentation, LF line endings, UTF-8 charset, final newlines inserted.
