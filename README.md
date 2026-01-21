# URL Shortener API Documentation

A lightweight, in-memory URL shortening service that converts long URLs into short, easy-to-share codes and provides redirect functionality.

## Overview

This API provides a simple way to create shortened URLs and redirect to the original destination. All data is stored in-memory, making it fast but non-persistent across server restarts.


## Features

- Create shortened URLs with unique short codes
- Redirect from short codes to original URLs
- In-memory storage for fast access
- Unique ID generation for short codes
- URL encoding/decoding support

## API Endpoints

### 1. Create Short URL

**Endpoint:** `POST /shorten`

Creates a new short URL and returns the generated short code.

**Request Body:**
```json
{
  "url": "https://www.example.com/very/long/path/to/some/resource?param=value&other=123"
}
```

**Response (201 Created):**
```json
{
  "url": "http://localhost:8080/abc123"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/golang/go"
  }'
```

**Response:**
```json
{
  "url": "http://localhost:8080/x7k9m2"
}
```

### 2. Redirect to Original URL

**Endpoint:** `GET /:shortcode`

Redirects to the original URL associated with the given short code.

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| shortcode | string | The unique short code generated from POST /shorten |

**Response:**
- `302 Found` - Redirect to original URL (if shortcode exists)
- `404 Not Found` - If shortcode doesn't exist

**Example:**
```bash
curl -L http://localhost:8080/x7k9m2
# Will redirect to: https://github.com/golang/go
```

## Usage Examples

### Creating Multiple Short URLs

```bash
# First URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com/page1"}'

# Second URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com/page2"}'
```

### Following Redirects

```bash
# Without -L flag (shows redirect headers)
curl -i http://localhost:8080/x7k9m2

# With -L flag (follows redirect automatically)
curl -L http://localhost:8080/x7k9m2
```

## Request/Response Examples

### Success Case

**Request:**
```bash
POST /shorten HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "url": "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status"
}
```

**Response:**
```bash
HTTP/1.1 201 Created
Content-Type: application/json

{
  "url": "http://localhost:8080/mdn9x3"
}
```

### Redirect Case

**Request:**
```bash
GET /mdn9x3 HTTP/1.1
Host: localhost:8080
```

**Response:**
```bash
HTTP/1.1 302 Found
Location: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
```

### Error Case (Not Found)

**Request:**
```bash
GET /invalid HTTP/1.1
Host: localhost:8080
```

**Response:**
```bash
HTTP/1.1 404 Not Found
Content-Type: application/json

{
  "error": "Link not found"
}
```

## Error Handling

| Status Code | Error | Description |
|-------------|-------|-------------|
| 400 | Invalid Request | Missing or invalid URL in request body |
| 404 | Not Found | Short code doesn't exist in the system |
| 500 | Server Error | Internal server error |

## Technical Details

### Short Code Generation

- **Length:** Typically 6-8 characters
- **Characters:** Alphanumeric (0-9, a-z, A-Z)
- **Uniqueness:** Generated using secure random or sequential ID generation
- **Collision Handling:** Retry generation if code already exists

### URL Encoding

- Original URLs are stored as-is
- URL validation ensures proper format
- Special characters are preserved

### In-Memory Storage

- Uses Go maps for O(1) lookup and insertion
- Non-persistent (data lost on server restart)
- No database required
- Fast performance for typical use cases

## Constraints & Limitations

- **Data Persistence:** All URLs are lost when server restarts
- **Scalability:** Memory usage grows with number of shortened URLs
- **Concurrent Access:** May need mutex for thread-safe operations in high-concurrency scenarios
- **URL Validation:** Basic URL format validation recommended
- **Rate Limiting:** Not implemented (consider adding for production)

## Implementation Guide

### Key Components

1. **URL Storage** - In-memory map to store short code → original URL mapping
2. **ID Generator** - Unique random/sequential ID generator for short codes
3. **HTTP Handlers** - POST handler for shortening, GET handler for redirects
4. **URL Validation** - Ensure valid URL format before storing

### Dependencies

```go
import (
	"encoding/json"
	"net/http"
	"crypto/rand"
)
```