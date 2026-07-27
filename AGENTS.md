# AGENTS.md

Guidance for AI coding agents working in this repository.

## Project Overview

Go SDK for the [S.EE](https://s.ee) URL shortener service (`github.com/sdotee/sdk.go`). The package name is `seesdk`. It has **zero external dependencies** — only the Go standard library is used. Requires Go 1.21+.

## Repository Layout

Single flat package (`seesdk`) with four source files:

- **`client.go`** — `Client` struct, `Config`, `NewClient()` constructor, and private HTTP helpers (`doRequest` for JSON requests, `doMultipartRequest` for file uploads). Authentication sets the `Authorization` header directly (no `Bearer` prefix).
- **`api.go`** — All 14 public methods on `Client` (CRUD for short URLs, text, and files; read-only for domains, tags, usage, and file history). Uses a generic `unmarshalResponse[T]` helper for JSON deserialization.
- **`models.go`** — Request/response structs with JSON tags. Responses follow a common `{code, data, message}` envelope pattern. `UploadFileData` is a shared named struct used by both `UploadFileResponse` and `GetFileHistoryResponse`.
- **`client_test.go`** — Integration tests covering the full API surface. Test fixtures live in `testdata/`.
- **`examples/main.go`** — Usage example.

## Build and Test Commands

```bash
go build ./...                              # Build
go vet ./...                                # Static analysis
gofmt -l .                                  # Check formatting
go test ./...                               # Run all tests (requires SEE_API_KEY)
go test -run TestNewClient ./...            # Run a single test
go test -v ./...                            # Verbose test output
```

## Testing

All tests are **integration tests** that hit the live S.EE API. They require environment variables:

- `SEE_API_KEY` (required) — tests skip gracefully via `t.Skip()` if not set
- `SEE_BASE_URL` (optional) — defaults to `https://s.ee/api/v1`

Tests use the shared `setupTestClient(t)` helper that reads the env vars and skips when credentials are missing. Without `SEE_API_KEY`, `go test ./...` passes with all tests skipped — this is expected in sandboxed environments.

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
| `GetFileHistory` | GET | `/files` |
| `DeleteFile` | GET | `/file/delete/:key` |
| `GetDomains` | GET | `/domains` |
| `GetFileDomains` | GET | `/file/domains` |
| `GetTextDomains` | GET | `/text/domains` |
| `GetTags` | GET | `/tags` |
| `GetUsage` | GET | `/usage` |

## Coding Conventions

- File headers include copyright, author, and timestamp comments — preserve this style in new files.
- JSON request structs use `json:"field_name,omitempty"` tags; multipart requests (e.g. `UploadFileRequest`) use plain Go fields without JSON tags.
- `doMultipartRequest` accepts a `fields map[string]string` for additional form fields (e.g. `is_private`, `domain`, `custom_slug`).
- All public methods return `(*ResponseType, error)`.
- Errors are wrapped with `fmt.Errorf("context: %w", err)`.
- File uploads are capped at 100MB (`checkFileSize` validates via `Stat()` or `Len()` interface).
- `UsageNoLimit = -1` is a sentinel value for unlimited usage quotas.
- Do not add external dependencies — keep the SDK standard-library only.

## Style / EditorConfig

Per `.editorconfig`: 4-space indentation, LF line endings, UTF-8 charset, final newlines inserted. Always run `gofmt` on Go files before committing.
