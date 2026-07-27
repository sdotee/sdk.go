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

Tests are split into two files:

- **`unit_test.go`** — Offline unit tests against a `httptest` mock server (helpers, header injection, upload routing, TUS chunk flow, validation). Always run; no credentials needed.
- **`client_test.go`** — Integration tests that hit the live S.EE API. They require environment variables:

- `SEE_API_KEY` (required) — tests skip gracefully if not set
- `SEE_BASE_URL` (optional) — defaults to `https://s.ee/api/v1`

Tests use `setupTestClient(t)` as a shared helper that reads env vars and calls `t.Skip()` when credentials are missing.

## Architecture

Single package (`seesdk`), flat structure:

- **`client.go`** — `Client` struct, `Config`, `NewClient()` constructor, `Version` constant, and private HTTP helpers (`doRequest` for JSON, `doMultipartRequest` for file uploads, shared `do`/`setCommonHeaders` for header injection and response handling). Auth is done by setting the `Authorization` header directly (not Bearer prefix). A `User-Agent: see-go-sdk/<version>` header (the `userAgent` constant) is sent on every request.
- **`api.go`** — Core public methods on `Client` (CRUD for short URLs, text, files; read-only for domains, tags, usage, visit stats, link/text/file history, private file download URLs; token check). Uses a generic `callAPI[T]` helper that combines request execution and JSON deserialization, plus `pagedEndpoint` for paginated GET endpoints and `readerSize` for size detection via `Stat()`/`Len()`. `SmartUploadFile` auto-switches between regular multipart upload (≤100MB) and TUS large file upload based on detected size.
- **`largefile.go`** — Large file upload support (up to 5GB) via the TUS 1.0 resumable protocol: `CreateLargeFileUpload`, `UploadLargeFileChunk` (PATCH), `GetLargeFileUploadOffset` (HEAD), `GetLargeFileUploadProgress`, `CompleteLargeFileUpload`, `CancelLargeFileUpload`, and the one-call `UploadLargeFile` convenience (handles fast/deduplicated uploads via `file_hash`). Private `doTUS` helper sets `Tus-Resumable: 1.0.0` and common headers.
- **`bio.go`** — Bio page methods: `CreateBioPage`, `UpdateBioPage`, `DeleteBioPage`, `GetBioPageHistory`.
- **`qrcode.go`** — QR code methods: `CreateQRCode`, `DeleteQRCode`, `GetQRCodeHistory`.
- **`models.go`** — Request/response structs with JSON tags. Responses follow a common `{code, data, message}` envelope pattern. `UploadFileData` is a shared named struct used by both `UploadFileResponse` and `GetFileHistoryResponse`.
- **`client_test.go`** — Integration tests covering the full API surface. Test fixtures live in `testdata/`.
- **`unit_test.go`** — Offline unit tests using `httptest` mock servers; run without any credentials.

## API Endpoints Mapped to Methods

| Method | HTTP | Endpoint |
|---|---|---|
| `CreateShortURL` | POST | `/shorten` |
| `UpdateShortURL` | PUT | `/shorten` |
| `DeleteShortURL` | DELETE | `/shorten` |
| `GetLinkVisitStat` | GET | `/link/visit-stat` |
| `CreateText` | POST | `/text` |
| `UpdateText` | PUT | `/text` |
| `DeleteText` | DELETE | `/text` |
| `UploadFile` | POST (multipart) | `/file/upload` |
| `SmartUploadFile` | — | delegates to `UploadFile` or `UploadLargeFile` by size |
| `CreateLargeFileUpload` | POST | `/file/large-file/create` |
| `UploadLargeFileChunk` | PATCH (TUS) | `/file/large-file-tus/:upload_id` |
| `GetLargeFileUploadOffset` | HEAD (TUS) | `/file/large-file-tus/:upload_id` |
| `GetLargeFileUploadProgress` | GET | `/file/large-file/progress` |
| `CompleteLargeFileUpload` | POST | `/file/large-file/complete` |
| `CancelLargeFileUpload` | DELETE | `/file/large-file/cancel` |
| `GetFileHistory` | GET | `/files` |
| `DeleteFile` | GET | `/file/delete/:key` |
| `GetPrivateFileDownloadURL` | GET | `/file/private/download-url` |
| `GetDomains` | GET | `/domains` |
| `GetFileDomains` | GET | `/file/domains` |
| `GetTextDomains` | GET | `/text/domains` |
| `GetTags` | GET | `/tags` |
| `GetUsage` | GET | `/usage` |
| `GetLinkHistory` | GET | `/links` |
| `GetTextHistory` | GET | `/texts` |
| `CheckToken` | POST | `/token/check` |
| `CreateBioPage` | POST | `/bio` |
| `UpdateBioPage` | PUT | `/bio` |
| `DeleteBioPage` | DELETE | `/bio` |
| `GetBioPageHistory` | GET | `/bios` |
| `CreateQRCode` | POST | `/qrcode` |
| `DeleteQRCode` | DELETE | `/qrcode` |
| `GetQRCodeHistory` | GET | `/qrcodes` |

Not implemented: `GET /shorten` (simple mode) — a signature-in-query bookmarklet variant redundant with `CreateShortURL`.

## Conventions

- File headers include copyright, author, and timestamp comments
- JSON request structs use `json:"field_name,omitempty"` tags; multipart requests (e.g. `UploadFileRequest`) use plain Go fields without JSON tags
- `doMultipartRequest` accepts a `fields map[string]string` for additional form fields (e.g. `is_private`, `domain`, `custom_slug`)
- All public methods return `(*ResponseType, error)`
- Errors are wrapped with `fmt.Errorf("context: %w", err)`; the shared `errNilFile` sentinel is returned for nil upload readers
- File uploads are capped at 100MB (`maxUploadFileSize`, validated via `readerSize` which checks `Stat()`/`Len()`); large file uploads are capped at 5GB (`maxLargeFileSize`)
- `UsageNoLimit = -1` sentinel for unlimited usage quotas
- `VisitStatPeriod*` constants for `GetLinkVisitStat` period values (`daily`, `monthly`, `totally`)
- `LargeFileUploadStatus*` constants for upload session status (1=uploading, 2=completed, 3=failed, 4=cancelled)
- Paginated GET endpoints (`/files`, `/links`, `/texts`, `/bios`, `/qrcodes`) share the `pagedEndpoint` helper; page ≤ 1 means first page
- SDK version lives in the `Version` constant in `client.go`; keep it in sync with git tags

## EditorConfig

4-space indentation, LF line endings, UTF-8 charset, final newlines inserted.
