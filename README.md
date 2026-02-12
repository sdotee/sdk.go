# See Go SDK

Official Golang SDK for [S.EE](https://s.ee) URL shortener service. Create, manage, and track short URLs with ease.

## Features

- 🔗 Create short URLs with custom slugs
- 📝 Create text/paste with syntax highlighting
- 📂 File upload and sharing (public or private)
- 📜 File upload history with pagination
- 🔒 Password-protected links
- ⏰ Expiration time support
- 🏷️ Tag management for organization
- 🌐 Multiple domain support
- 📈 View account usage statistics

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

// Delete file using hash
client.DeleteFile(uploadResp.Data.Hash)
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

**GetFileHistory(page int)** - Get paginated file upload history (30 per page)

**DeleteFile(deleteKey string)** - Delete a file using the delete key

**GetUsage()** - Get account usage statistics

**GetDomains()** - List available domains

**GetFileDomains()** - List available domains for file sharing

**GetTextDomains()** - List available domains for text sharing

**GetTags()** - List available tags

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
