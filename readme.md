# Drop CLI - Secure, zero-knowledge, secret sharing CLI

[![CI](https://github.com/cainseing/drop-cli/actions/workflows/ci.yml/badge.svg?branch=development)](https://github.com/cainseing/drop-cli/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.26-blue)](https://golang.org/dl/)

**Drop CLI** is a command-line tool for securely sharing sensitive data—such as API keys, tokens, and credentials—through the [Drop API](https://github.com/cainseing/drop-api). It uses end-to-end encryption to ensure secrets remain private and accessible only to intended recipients.

Ideal for development teams that need a fast, secure way to exchange confidential information.

---

## Security & Technical Workflow

Drop follows a zero-trust, client-side encryption model:

1. **Client-Side Encryption**  
   Uses **AES-256-GCM** for authenticated encryption before data leaves your machine.

2. **Size Obfuscation**  
   Applies binary padding to prevent traffic analysis based on payload length.

3. **Secure Encapsulation**  
   Prepends a 4-byte header inside the encrypted payload to enable accurate decoding and decryption.

4. **Zero-Knowledge Transport**  
   Generates a composite token containing:
    - Protocol Version
    - Identifier
    - Encryption Key

   The encryption key is never transmitted or stored on the server.

---

## Installation

Install the `drop` binary to your system path.

### Homebrew (macOS)

    brew tap cainseing/tap &&
    brew install drop

### Install Script (macOS & Linux)

Automatically detects your platform and installs the correct binary:

    curl -sL getdrop.dev/install.sh | bash

---

## Usage

### Create a Secret

#### Standard Input

    drop "api_key"

#### From a Pipe

    cat .env | drop

---

### Retrieve a Secret

    drop get <token>

    or

    drop get <token> > .env

You can also pass the full token directly (with the `drop_` prefix) and it will be treated as a `get`:

    drop drop_<token>

---

### Purge a Secret

    drop purge <token>

---

### Manage Identity

Set a GitHub or GitLab profile to associate with signed drops:

    drop identity set github <username>
    drop identity set gitlab <username>

---

### Check Version

    drop version

---

## Command-Line Options

### Drop Flags

| Flag | Long Form    | Description                          | Default |
|------|--------------|--------------------------------------|---------|
| `-t` | `--ttl`      | Expiry time in minutes               | `5`     |
| `-r` | `--reads`    | Maximum number of allowed reads      | `1`     |
| `-s` | `--signed`   | Sign the drop with SSH key           | `false` |
| `-c` | `--copy`     | Copy the token to clipboard          | `false` |

Example:

    drop -t 120 -r 3 -c "Temporary secret"

---

## Build & Release

The project supports cross-compilation for major platforms.

### Build All Targets

    make release

This generates binaries for:

- Linux
- macOS

---

## Contributing

Contributions are welcome via a pull request

---

## Security & Support

- Report security issues **privately**.
- Open issues on GitHub for bugs or feature requests. 
