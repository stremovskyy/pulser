# Pulser

[![Go Report Card](https://goreportcard.com/badge/github.com/stremovskyy/pulser)](https://goreportcard.com/report/github.com/stremovskyy/pulser)
[![GoDoc](https://godoc.org/github.com/stremovskyy/pulser?status.svg)](https://godoc.org/github.com/stremovskyy/pulser)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Pulser is a Go library for tracking user events and analytics with support for service codes and metadata. It provides a simple and efficient interface for event tracking in Go applications.

## Features

- Simple and intuitive API
- Support for user event tracking with metadata
- Service code support for multi-tenant applications
- Context-aware operations
- Graceful shutdown handling

## Installation

```bash
go get github.com/stremovskyy/pulser
```

## Quick Start

```go
package main

import (
    "context"
    "time"
    "github.com/stremovskyy/pulser"
)

func main() {
    client := pulser.NewClient()
    defer client.Shutdown(5 * time.Second)

    ctx := context.Background()
    
    // Track a simple event
    err := client.Track(ctx, "user123", pulser.EventTypeUserAction, pulser.EventSubTypeLogin, map[string]interface{}{
        "ip": "192.168.1.1",
        "device": "mobile",
    })
    if err != nil {
        // Handle error
    }

    // Track an event with service code
    err = client.TrackWithServiceCode(ctx, "service1", "user123", pulser.EventTypeUserAction, pulser.EventSubTypePurchase, map[string]interface{}{
        "amount": 99.99,
        "currency": "USD",
    })
    if err != nil {
        // Handle error
    }
}
```

## Documentation

For detailed documentation, visit [GoDoc](https://godoc.org/github.com/stremovskyy/pulser).

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 