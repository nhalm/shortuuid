# ShortUUID Go

[![CI/CD Pipeline](https://github.com/nhalm/shortuuid/actions/workflows/ci.yml/badge.svg)](https://github.com/nhalm/shortuuid/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/nhalm/shortuuid.svg)](https://pkg.go.dev/github.com/nhalm/shortuuid)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhalm/shortuuid)](https://goreportcard.com/report/github.com/nhalm/shortuuid)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go implementation of the Ruby [shortuuid](https://github.com/sudhirj/shortuuid.rb) library that produces **identical** short IDs for the same UUIDs.

## Features

- **100% Ruby Compatible**: Produces identical short IDs for the same UUIDs
- **Simple API**: Just 4 public functions for all use cases
- **UUID Type Support**: Works with both strings and `uuid.UUID` types
- **UUID Version Preservation**: Maintains UUID version (v4, v7, etc.) through encode/decode
- **High Performance**: Optimized for speed with minimal allocations
- **Typed Error Handling**: Separate error types for encoding and decoding operations

## Installation

```bash
go get github.com/nhalm/shortuuid
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/nhalm/shortuuid"
)

func main() {
    // Shorten a UUID string (Ruby compatible)
    uuid := "4890586e-32a5-4f9c-a000-2a2bb68eb1ce"
    short, err := shortuuid.Shorten(uuid)
    if err != nil {
        panic(err)
    }
    fmt.Println(short) // "2CvPdpytrcURpSLoPxYb30"
    
    // Expand back to UUID
    expanded, err := shortuuid.Expand(short)
    if err != nil {
        panic(err)
    }
    fmt.Println(expanded) // "4890586e-32a5-4f9c-a000-2a2bb68eb1ce"
}
```

### Working with UUID Types

```go
import (
    "github.com/google/uuid"
    "github.com/nhalm/shortuuid"
)

// Generate a UUIDv4 and shorten it
uuidv4 := uuid.New()
short, err := shortuuid.ShortenUUID(uuidv4)
if err != nil {
    panic(err)
}

// Expand back to UUID type
expanded, err := shortuuid.ExpandUUID(short)
if err != nil {
    panic(err)
}
fmt.Printf("Original: %s\nExpanded: %s\n", uuidv4, expanded)
```

## Error Handling

ShortUUID uses typed errors for better error handling:

### Error Types

- `EncodeError`: Errors during UUID shortening (invalid UUID format, etc.)
- `DecodeError`: Errors during short ID expansion (invalid characters, etc.)

### Using Errors

```go
import (
    "errors"
    "github.com/nhalm/shortuuid"
)

// Handle encoding errors
_, err := shortuuid.Shorten("invalid-uuid")
if err != nil {
    var encodeErr *shortuuid.EncodeError
    if errors.As(err, &encodeErr) {
        fmt.Printf("UUID: %s\n", encodeErr.UUID)
        fmt.Printf("Reason: %s\n", encodeErr.Reason)
        fmt.Printf("Error: %s\n", encodeErr.Error())
    }
}

// Handle decoding errors
_, err = shortuuid.Expand("invalid@chars")
if err != nil {
    var decodeErr *shortuuid.DecodeError
    if errors.As(err, &decodeErr) {
        fmt.Printf("Short ID: %s\n", decodeErr.ShortID)
        fmt.Printf("Reason: %s\n", decodeErr.Reason)
        fmt.Printf("Error: %s\n", decodeErr.Error())
    }
}
```

## API Reference

### Functions

```go
// String-based functions (Ruby compatible)
func Shorten(uuidStr string) (string, error)
func Expand(shortID string) (string, error)

// UUID type-based functions
func ShortenUUID(uuid uuid.UUID) (string, error)
func ExpandUUID(shortID string) (uuid.UUID, error)
```

### Error Types

```go
type EncodeError struct {
    UUID   string // The UUID that caused the error
    Reason string // Description of what went wrong
}

func (e *EncodeError) Error() string

type DecodeError struct {
    ShortID string // The short ID that caused the error
    Reason  string // Description of what went wrong
}

func (e *DecodeError) Error() string
```

## Ruby Compatibility

This library is designed to be 100% compatible with the Ruby shortuuid library:

- Uses the same default alphabet: `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`
- Produces identical short IDs for the same UUIDs
- Maintains the same encoding/decoding behavior

## Performance

The library is optimized for performance:

- Encode: ~2300ns per operation
- Decode: ~762ns per operation
- Minimal memory allocations
- Efficient big integer arithmetic

## License

MIT License 