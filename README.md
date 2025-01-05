# Temp Mail Go Client
[![Go Reference](https://pkg.go.dev/badge/github.com/temp-mail-io/temp-mail-go.svg)](https://pkg.go.dev/github.com/temp-mail-io/temp-mail-go)
[![Go](https://github.com/temp-mail-io/temp-mail-go/actions/workflows/test.yml/badge.svg)](https://github.com/temp-mail-io/temp-mail-go/actions)

The **official Go Client** for [Temp Mail](https://temp-mail.io). This library provides developers a straightforward way to create and manage temporary email addresses, retrieve and delete messages, all via the Temp Mail API.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
    - [Listing Domains](#listing-domains)
    - [Getting Rate Limits](#getting-rate-limits)
    - [Creating Temporary Email](#creating-temporary-email)
    - [Fetching and Deleting Messages](#fetching-and-deleting-messages)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Features
- **Create** temporary email addresses with optional domain specifications
- **Get** current rate limits after API request
- **Delete** a temporary email along with all its messages
- **Retrieve** all messages for a specified email
- **Get** a specific message or download its attachment

## Installation
To install this Go package, run:
```bash
go get github.com/temp-mail-io/temp-mail-go
```

## Quick Start
Below is a simple example to get started:
```go
package main

import (
	"context"
	"log"

	"github.com/temp-mail-io/temp-mail-go"
)

func main() {
	// Replace with your real API key
	client := tempmail.NewClient("YOUR_API_KEY", nil)

	email, _, err := client.CreateEmail(context.Background(), tempmail.CreateEmailOptions{})
	if err != nil {
		log.Fatalf("Failed to create temporary email: %v", err)
	}

	// Use the created temporary email on the website, service, etc...
	...

	// Fetch messages for the created email
	data, _, err := client.ListEmailMessages(context.Background(), email.Email)
	if err != nil {
		log.Fatalf("Failed to fetch messages: %v", err)
	}

	for _, m := range data.Messages {
		// Iterate over messages
	}
}
```

## Usage Examples
### Listing Domains
```go
domains, _, err := client.ListDomains(context.Background())
if err != nil {
    // handle error
}
for _, d := range domains.Domains {
    fmt.Printf("Domain: %s, Type: %s\n", d.Name, d.Type)
}
```

### Getting Rate Limits
```go
rate, _, err := client.RateLimit(context.Background())
if err != nil {
    // handle error
}
fmt.Printf("limit: %d, used: %d, remaining: %d, reset: %s\n", rate.Limit, rate.Used, rate.Remaining, rate.Reset)
// Output: limit: 1000, used: 0, remaining: 1000, reset: 2025-01-31 23:59:59 +0000 UTC
```

You can also get the rate limits from the response headers:
```go
email, resp, err := client.CreateEmail(context.Background(), tempmail.CreateEmailOptions{})
if err != nil {
    // handle error
}
fmt.Printf("Rate limit: %d\n", resp.Rate.Limit)
```

### Creating Temporary Email
```go
email, _, err := client.CreateEmail(context.Background(), tempmail.CreateEmailOptions{
	Domain: "example.com",
})
if err != nil {
    // handle error
}
fmt.Printf("Created temporary email: %s (TTL: %s)\n", email.Email, email.TTL)
```

### Fetching and Deleting Messages
```go
messages, _, err := client.ListEmailMessages(context.Background(), "your_email@example.com")
if err != nil {
    // handle error
}
fmt.Printf("Fetched %d messages.\n", len(messages))

// Deleting a specific message
_, err = client.DeleteMessage(context.Background(), messages[0].ID)
if err != nil {
    // handle error
}
```

## Testing
We use the Go testing framework with both unit tests and optional integration tests.

Run tests locally:
```bash
go test ./... -v
```

In CI, the tests are automatically executed via [GitHub Actions](https://github.com/temp-mail-io/temp-mail-go/actions).

## Contributing
We welcome and appreciate contributions! Please see our CONTRIBUTING.md for guidelines on how to open issues, submit pull requests, and follow our coding standards.

## License
This project is licensed under the MIT License.

## Support
If you encounter any issues, please open [an issue](https://github.com/temp-mail-io/temp-mail-go/issues) on GitHub. We are happy to help you!
