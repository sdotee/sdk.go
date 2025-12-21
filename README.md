# See Go SDK

Official Golang SDK for [S.EE](https://s.ee) URL shortener service. Create, manage, and track short URLs with ease.

## Features

- 🔗 Create short URLs with custom slugs
- 📝 Create text/paste with syntax highlighting
- �🔒 Password-protected links
- ⏰ Expiration time support
- 🏷️ Tag management for organization
- 🌐 Multiple domain support
- 📊 Track and analyze link performance

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
    BaseCreateRequest: seesdk.BaseCreateRequest{
        Title: "My Link",
    },
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
    fmt.Printf("%s (ID: %d)\n", tag.Name, tag.Id)
}
```

### Advanced Short URL Creation

Create a custom branded link with expiration and password protection:

```go
expireAt := time.Now().Add(30 * 24 * time.Hour).Unix()

resp, err := client.CreateShortURL(seesdk.CreateShortURLRequest{
    TargetURL:  "https://www.example.com/campaign",
    Domain:     "s.ee",
    BaseCreateRequest: seesdk.BaseCreateRequest{
        CustomSlug: "summer-sale",
        ExpireAt:   expireAt,
        Password:   "secret123",
        Title:      "Summer Sale Campaign",
        TagIDs:     []int64{1, 2},
    },
})
```

### Update and Delete

```go
// Update existing short URL
client.UpdateShortURL(seesdk.UpdateShortURLRequest{
    BaseSlugRequest: seesdk.BaseSlugRequest{
        Domain: "s.ee",
        Slug:   "summer-sale",
    },
    TargetURL: "https://www.example.com/new-campaign",
    Title:     "Updated Campaign",
})

// Delete short URL
client.DeleteShortURL(seesdk.DeleteURLRequest{
    BaseSlugRequest: seesdk.BaseSlugRequest{
        Domain: "s.ee",
        Slug:   "summer-sale",
    },
})

### Text Management

```go
// Create a new text/paste
textResp, err := client.CreateText(seesdk.CreateTextRequest{
    Content:    "fmt.Println(\"Hello World\")",
    Domain:     "s.ee",
    Title:      "Go Hello World",
    TextType:   "go", // Syntax highlighting
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

**GetDomains()** - List available domains

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
| TargetUrl | string | Yes      |
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
| TextType   | string  | No       | Syntax highlighting type |
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
