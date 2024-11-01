## Go (Golang) Wabot API Client

### README.md

```markdown
# Wabot API Client for Go (Golang)
```

This is a Go client library for interacting with the Wabot API. It handles authentication, token management, and provides methods to interact with Wabot API endpoints.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Initialization](#initialization)
  - [Authentication](#authentication)
  - [Getting Templates](#getting-templates)
  - [Sending Messages](#sending-messages)
  - [Logout](#logout)
- [Example](#example)
- [Notes](#notes)
- [License](#license)

## Prerequisites

- Go 1.13 or higher
- Third-party packages:
  - `github.com/dgrijalva/jwt-go`
  - `github.com/pkg/errors`

## Installation

1. **Install Dependencies**

   ```bash
   go get github.com/dgrijalva/jwt-go
   go get github.com/pkg/errors
```
1. **Include the Client**

- Save the wabot_api_client.go file in a package directory (e.g., wabot) within your project and import it in your code:

   ```bash
   import "path/to/your/package/wabot"
  ```
- Replace "path/to/your/package/wabot" with the actual import path where you placed the wabot package.


