# See Go SDK

Official Golang SDK for [S.EE](https://s.ee) URL shortener service. Create, manage, and track short URLs with ease.

## Features

- 🔗 Create short URLs with custom slugs
- 📝 Create text/paste with syntax highlighting
- 📂 File upload and sharing (public or private)
- 📦 Large file uploads up to 5GB (TUS resumable protocol, with instant deduplicated uploads)
- 🚀 Smart upload that automatically picks the best upload strategy by file size
- 📜 Link, text, and file history with pagination
- 🔑 Temporary download URLs for private files
- 👤 Bio pages with custom links
- 📱 Dynamic QR codes (PNG/SVG/PDF)
- 🔒 Password-protected links
- ⏰ Expiration time support
- 🏷️ Tag management for organization
- 🌐 Multiple domain support
- 📈 View account usage and link visit statistics
- ✅ API token validation

## Installation

```bash
go get github.com/sdotee/sdk.go
```

## Quick Start

Initialize the client with your API credentials:

```go
import seesdk "github.com/sdotee/sdk.go"

client := seesdk.NewClient(seesdk.Config{
    BaseURL: "https://api.s.ee",
    APIKey:  "your-api-key-here",
})
```

Create your first short URL:

```go
resp, err := client.CreateShortURL(seesdk.CreateShortURLRequest{
    TargetURL: "https://www.example.com/very/long/url",
    Domain:    "s.ee",
    Title:     "My Link",
})

fmt.Printf("Short URL: %s\n", resp.Data.ShortURL)
```

## Usage Examples

### Domain and Tag Management

```go
// Get available domains
domains, _ := client.GetDomains()
fmt.Println(domains.Data.Domains)

// Get available tags
tags, _ := client.GetTags()
for _, tag := range tags.Data.Tags {
    fmt.Printf("%s (ID: %d)\n", tag.Name, tag.ID)
}
```

### Advanced Short URL Creation

Create a custom branded link with expiration and password protection:

```go
expireAt := time.Now().Add(30 * 24 * time.Hour).Unix()

resp, err := client.CreateShortURL(seesdk.CreateShortURLRequest{
    TargetURL:  "https://www.example.com/campaign",
    Domain:     "s.ee",
    CustomSlug: "summer-sale",
    ExpireAt:   expireAt,
    Password:   "secret123",
    Title:      "Summer Sale Campaign",
    TagIDs:     []int64{1, 2},
})
```

### Statistics

```go
// Get account usage statistics
usage, _ := client.GetUsage()
fmt.Printf("Links created today: %d/%d\n",
    usage.Data.LinkCountDay,
    usage.Data.LinkCountDayLimit)
fmt.Printf("Storage used: %s MB / %s MB\n",
    usage.Data.StorageUsageMB,
    usage.Data.StorageUsageLimitMB)

// Get visit statistics for a short URL
// Period: seesdk.VisitStatPeriodDaily, seesdk.VisitStatPeriodMonthly,
// or seesdk.VisitStatPeriodTotally (empty string defaults to all-time)
stat, _ := client.GetLinkVisitStat("s.ee", "summer-sale", seesdk.VisitStatPeriodTotally)
fmt.Printf("Total visits: %d\n", stat.Data.VisitCount)
```

### Update and Delete

```go
// Update existing short URL
client.UpdateShortURL(seesdk.UpdateShortURLRequest{
    Domain:    "s.ee",
    Slug:      "summer-sale",
    TargetURL: "https://www.example.com/new-campaign",
    Title:     "Updated Campaign",
})

// Delete short URL
client.DeleteShortURL(seesdk.DeleteURLRequest{
    Domain: "s.ee",
    Slug:   "summer-sale",
})
```

### Text Management

```go
// Create a new text/paste
textResp, err := client.CreateText(seesdk.CreateTextRequest{
    Content:    "fmt.Println(\"Hello World\")",
    Domain:     "s.ee",
    Title:      "Go Hello World",
    TextType:   "source_code",
    CustomSlug: "hello-go",
})
fmt.Printf("Text URL: %s\n", textResp.Data.ShortURL)

// Update text
client.UpdateText(seesdk.UpdateTextRequest{
    Domain:  "s.ee",
    Slug:    "hello-go",
    Content: "fmt.Println(\"Hello Updated World\")",
    Title:   "Updated Go Hello World",
})

// Delete text
client.DeleteText(seesdk.DeleteTextRequest{
    Domain: "s.ee",
    Slug:   "hello-go",
})
```

### File Management

```go
// Upload a file
file, _ := os.Open("image.png")
defer file.Close()

uploadResp, err := client.UploadFile(seesdk.UploadFileRequest{
    Filename: "image.png",
    File:     file,
})
fmt.Printf("File URL: %s\n", uploadResp.Data.URL)
fmt.Printf("Delete Key: %s\n", uploadResp.Data.Hash)

// Upload a private file with custom domain and slug
privateResp, err := client.UploadFile(seesdk.UploadFileRequest{
    Filename:   "secret.pdf",
    File:       file,
    IsPrivate:  true,
    Domain:     "s.ee",
    CustomSlug: "my-file",
})

// Get file upload history (paginated, 30 per page)
history, _ := client.GetFileHistory(1)
for _, f := range history.Data {
    fmt.Printf("%s - %s\n", f.Filename, f.URL)
}

// Get available domains for file sharing
fileDomains, _ := client.GetFileDomains()
fmt.Println(fileDomains.Data.Domains)

// Get a temporary download URL for a private file (valid for ~1 hour)
dlResp, _ := client.GetPrivateFileDownloadURL(int64(privateResp.Data.FileID))
fmt.Printf("Download URL: %s (expires at %d)\n", dlResp.Data.URL, dlResp.Data.ExpiresAt)

// Delete file using hash
client.DeleteFile(uploadResp.Data.Hash)
```

### Large File Upload (up to 5GB)

For files larger than the 100MB `UploadFile` limit, use the TUS-based large file API:

```go
// One-call convenience: creates a session, uploads in 16MB chunks, completes
f, _ := os.Open("video.mp4")
defer f.Close()
info, _ := f.Stat()

resp, err := client.UploadLargeFile(seesdk.CreateLargeFileUploadRequest{
    FileName: "video.mp4",
    FileSize: info.Size(),
    // FileHash: "<sha256>", // optional: enables instant deduplicated upload
    // IsPrivate: 1,         // optional: 0 = public (default), 1 = private
}, f)
fmt.Printf("File URL: %s\n", resp.Data.File.URL)
```

Or drive the session manually for resumable uploads:

```go
// 1. Create an upload session
createResp, _ := client.CreateLargeFileUpload(seesdk.CreateLargeFileUploadRequest{
    FileName: "video.mp4",
    FileSize: info.Size(),
})
uploadID := createResp.Data.UploadID

// 2. Upload chunks (TUS PATCH); resume with GetLargeFileUploadOffset
offset, _ := client.GetLargeFileUploadOffset(uploadID)
offset, _ = client.UploadLargeFileChunk(uploadID, offset, chunk)

// 3. Check progress at any time
progress, _ := client.GetLargeFileUploadProgress(uploadID)
fmt.Printf("Progress: %.1f%%\n", progress.Data.Progress)

// 4. Complete (or cancel) the session
completeResp, _ := client.CompleteLargeFileUpload(uploadID)
// client.CancelLargeFileUpload(uploadID)
```

### Smart Upload

`SmartUploadFile` automatically chooses between the regular upload (≤100MB) and
the TUS large file upload (up to 5GB) based on the file size:

```go
f, _ := os.Open("any-file.bin")
defer f.Close()

resp, err := client.SmartUploadFile(seesdk.UploadFileRequest{
    Filename: "any-file.bin",
    File:     f,
})
fmt.Printf("File URL: %s\n", resp.Data.URL)
```

The size is detected via `Stat()` (e.g. `*os.File`) or `Len()` (e.g.
`*bytes.Reader`). When the size cannot be determined, the regular upload is used.

### History

```go
// Short link creation history (paginated)
links, _ := client.GetLinkHistory(1)
for _, l := range links.Data {
    fmt.Printf("%s -> %s (%d visits)\n", l.ShortURL, l.TargetURL, l.VisitCount)
}

// Text creation history (paginated)
texts, _ := client.GetTextHistory(1)
for _, tx := range texts.Data {
    fmt.Printf("%s: %s\n", tx.ShortURL, tx.ContentPreview)
}
```

### Bio Pages

```go
// Create a bio page with custom links
bioResp, _ := client.CreateBioPage(seesdk.CreateBioPageRequest{
    Title:       "My Bio",
    Description: "About me",
    MastodonURL: "https://mastodon.social/@me",
    CustomLinks: []seesdk.BioCustomLink{
        {Title: "Blog", URL: "https://blog.example.com"},
    },
})
fmt.Printf("Bio page: %s\n", bioResp.Data.ShortURL)

// Update a bio page
client.UpdateBioPage(seesdk.UpdateBioPageRequest{
    ID:    bioResp.Data.BioPageID,
    Title: "My Updated Bio",
})

// List bio pages (paginated)
bios, _ := client.GetBioPageHistory(1)
fmt.Printf("Total bio pages: %d\n", bios.Data.Total)

// Delete a bio page
client.DeleteBioPage(bioResp.Data.BioPageID)
```

### QR Codes

```go
// Create a dynamic QR code (PNG/SVG/PDF are generated server-side)
qrResp, _ := client.CreateQRCode(seesdk.CreateQRCodeRequest{
    TargetURL: "https://example.com",
    Title:     "My QR Code",
})
fmt.Printf("QR PNG: %s\n", qrResp.Data.PNGURL)

// List QR codes (paginated)
qrcodes, _ := client.GetQRCodeHistory(1)
fmt.Printf("Total QR codes: %d\n", qrcodes.Data.Total)

// Delete a QR code
client.DeleteQRCode(seesdk.DeleteQRCodeRequest{Domain: "s.ee", Slug: qrResp.Data.Slug})
```

### Token Validation

```go
check, _ := client.CheckToken("your-api-key-here")
if check.Data.Valid {
    fmt.Printf("Token valid until %d\n", check.Data.ExpiresAt)
}
```

## API Reference

### Client Configuration

| Field   | Type          | Required | Description                    |
| ------- | ------------- | -------- | ------------------------------ |
| BaseURL | string        | Yes      | API endpoint URL               |
| APIKey  | string        | Yes      | Your authentication token      |
| Timeout | time.Duration | No       | Request timeout (default: 30s) |

### Methods

**CreateShortURL(req CreateShortURLRequest)** - Create a new short URL

**UpdateShortURL(req UpdateShortURLRequest)** - Modify an existing short URL

**DeleteShortURL(req DeleteURLRequest)** - Remove a short URL

**CreateText(req CreateTextRequest)** - Create a new text entry

**UpdateText(req UpdateTextRequest)** - Modify an existing text entry

**DeleteText(req DeleteTextRequest)** - Remove a text entry

**UploadFile(req UploadFileRequest)** - Upload a file (max 100MB)

**SmartUploadFile(req UploadFileRequest)** - Upload a file, automatically switching to the TUS large file upload above 100MB

**UploadLargeFile(req CreateLargeFileUploadRequest, r io.Reader)** - Upload a file up to 5GB (full TUS flow)

**CreateLargeFileUpload(req CreateLargeFileUploadRequest)** - Create a TUS upload session

**UploadLargeFileChunk(uploadID string, offset int64, chunk []byte)** - Upload a chunk (TUS PATCH)

**GetLargeFileUploadOffset(uploadID string)** - Get the current upload offset (TUS HEAD)

**GetLargeFileUploadProgress(uploadID string)** - Get upload session progress

**CompleteLargeFileUpload(uploadID string)** - Finalize an upload session

**CancelLargeFileUpload(uploadID string)** - Cancel an upload session

**GetFileHistory(page int)** - Get paginated file upload history (30 per page)

**GetLinkHistory(page int)** - Get paginated short link creation history

**GetTextHistory(page int)** - Get paginated text creation history

**DeleteFile(deleteKey string)** - Delete a file using the delete key

**GetPrivateFileDownloadURL(fileID int64)** - Get a temporary download URL for a private file

**GetUsage()** - Get account usage statistics

**GetLinkVisitStat(domain, slug, period string)** - Get visit statistics for a short URL

**GetDomains()** - List available domains

**GetFileDomains()** - List available domains for file sharing

**GetTextDomains()** - List available domains for text sharing

**GetTags()** - List available tags

**CreateBioPage(req CreateBioPageRequest)** - Create a bio page

**UpdateBioPage(req UpdateBioPageRequest)** - Update a bio page

**DeleteBioPage(id int64)** - Delete a bio page by ID

**GetBioPageHistory(page int)** - Get paginated bio page list

**CreateQRCode(req CreateQRCodeRequest)** - Create a dynamic QR code

**DeleteQRCode(req DeleteQRCodeRequest)** - Delete a QR code

**GetQRCodeHistory(page int)** - Get paginated QR code list

**CheckToken(token string)** - Validate an API token

### Request Models

**CreateShortURLRequest**

| Field                 | Type    | Required | Description               |
| --------------------- | ------- | -------- | ------------------------- |
| TargetURL             | string  | Yes      | Destination URL           |
| Domain                | string  | Yes      | Short domain name         |
| CustomSlug            | string  | No       | Custom URL slug           |
| ExpireAt              | int64   | No       | Unix timestamp (seconds)  |
| Password              | string  | No       | Access password           |
| TagIDs                | []int64 | No       | Associated tag IDs        |
| Title                 | string  | No       | Link description          |
| ExpirationRedirectURL | string  | No       | Redirect after expiration |

**UpdateShortURLRequest**

| Field     | Type   | Required |
| --------- | ------ | -------- |
| Domain    | string | Yes      |
| Slug      | string | Yes      |
| TargetURL | string | Yes      |
| Title     | string | No       |

**DeleteURLRequest**

| Field  | Type   | Required |
| ------ | ------ | -------- |
| Domain | string | Yes      |
| Slug   | string | Yes      |

**CreateTextRequest**

| Field      | Type    | Required | Description              |
| ---------- | ------- | -------- | ------------------------ |
| Content    | string  | Yes      | Text content             |
| Domain     | string  | No       | Short domain name        |
| CustomSlug | string  | No       | Custom URL slug          |
| TextType   | string  | No       | plain_text, source_code, or markdown |
| Title      | string  | No       | Text title               |
| Password   | string  | No       | Access password          |
| ExpireAt   | int64   | No       | Unix timestamp (seconds) |
| TagIDs     | []int64 | No       | Associated tag IDs       |

**UpdateTextRequest**

| Field   | Type   | Required |
| ------- | ------ | -------- |
| Domain  | string | Yes      |
| Slug    | string | Yes      |
| Content | string | Yes      |
| Title   | string | No       |

**DeleteTextRequest**

| Field  | Type   | Required |
| ------ | ------ | -------- |
| Domain | string | Yes      |
| Slug   | string | Yes      |

**UploadFileRequest**

| Field      | Type      | Required | Description                          |
| ---------- | --------- | -------- | ------------------------------------ |
| Filename   | string    | Yes      | Name of the file                     |
| File       | io.Reader | Yes      | File content reader                  |
| Domain     | string    | No       | Domain for the short link            |
| CustomSlug | string    | No       | Custom slug for the file URL         |
| IsPrivate  | bool      | No       | Set to true for private file upload  |

**CreateLargeFileUploadRequest**

| Field       | Type   | Required | Description                                  |
| ----------- | ------ | -------- | -------------------------------------------- |
| FileName    | string | Yes      | Original filename                            |
| FileSize    | int64  | Yes      | File size in bytes (max 5GB)                 |
| FileHash    | string | No       | SHA256 hash for instant upload deduplication |
| Domain      | string | No       | Domain for the short link                    |
| Alias       | string | No       | Custom slug for the short link               |
| IsPrivate   | int    | No       | 0 = public (default), 1 = private            |
| Password    | string | No       | Access password (3-32 chars)                 |
| ExpireAt    | int64  | No       | Unix timestamp for link expiry               |
| Title       | string | No       | Title for the file                           |
| Description | string | No       | Description                                  |
| MimeType    | string | No       | MIME type hint                               |

**CreateBioPageRequest**

| Field       | Type            | Required | Description                            |
| ----------- | --------------- | -------- | -------------------------------------- |
| Title       | string          | Yes      | Bio page title                         |
| Description | string          | No       | Bio page description                   |
| Domain      | string          | No       | Domain for the bio page URL            |
| CustomSlug  | string          | No       | Custom URL slug                        |
| MastodonURL | string          | No       | Mastodon profile URL                   |
| RSSURL      | string          | No       | RSS feed URL                           |
| CustomLinks | []BioCustomLink | No       | Custom links (Title, URL, Description) |

**CreateQRCodeRequest**

| Field      | Type   | Required | Description               |
| ---------- | ------ | -------- | ------------------------- |
| TargetURL  | string | Yes      | Destination URL           |
| Title      | string | Yes      | QR code title             |
| Domain     | string | No       | Domain for the short link |
| CustomSlug | string | No       | Custom URL slug           |

## Error Handling

All methods return standard Go errors. Always check for errors:

```go
resp, err := client.CreateShortURL(req)
if err != nil {
    log.Printf("Failed: %v", err)
    return
}
```

## Example

See [examples/main.go](examples/main.go) for complete working examples.

```bash
cd examples && go run main.go
```

## Contributing

Issues and Pull Requests are welcome!

## License

MIT License
