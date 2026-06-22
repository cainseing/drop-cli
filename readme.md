# drop

[![CI](https://github.com/cainseing/drop-cli/actions/workflows/ci.yml/badge.svg?branch=development)](https://github.com/cainseing/drop-cli/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.26-blue)](https://golang.org/dl/)

A zero-knowledge CLI for securely sharing secrets. Secrets are encrypted locally before leaving your machine, ensuring the server only stores an unreadable ciphertext blob—along with minimal metadata if you choose to sign the drop.

---

## Install

**Homebrew (macOS)**

```sh
brew tap cainseing/tap && brew install drop
```

**Install script (macOS & Linux)**

```sh
curl -sL getdrop.dev/install.sh | bash
```

**From source**

```sh
git clone https://github.com/cainseing/drop-cli.git
cd drop-cli
make install   # builds and copies to /usr/local/bin
```

---

## Usage

For detailed usage instructions and examples, see the [documentation](https://getdrop.dev/#docs).

---

## How it works

1. The secret is encrypted locally with **AES-256-GCM** using a randomly generated key.
2. The plaintext is wrapped in an envelope (`[4-byte length][plaintext][random padding]`) before encryption to obfuscate payload size.
3. The encrypted blob is uploaded to the server. The server has no access to the key.
4. A token is generated locally:

   ```
   drop_<base64(protocol_version.blob_id.hex_key:signed_flag)>
   ```

   The encryption key is embedded in the token and never stored server-side. The
   `signed_flag` records whether the sender signed the drop, so a server cannot
   silently strip a signature without the recipient noticing. It is packed into
   the same segment as the key (rather than as its own dot-separated segment)
   so older CLI versions still recognize the token as having 3 parts and show
   a "please update" message instead of "Token provided is not valid".
5. When retrieved, the token is split, the blob is fetched, and the secret is decrypted locally.

See [SECURITY.md](SECURITY.md) for a full breakdown of the security model.

---

## Building

```sh
make build          # build for current platform → ./drop
make release        # cross-compile for Linux and macOS (amd64 + arm64) → ./bin/
make test           # run tests
make test-coverage  # run tests with coverage report
make lint           # run golangci-lint
make sec            # run gosec
```

---

## Contributing

Pull requests are welcome. Open an issue first for significant changes.

---

## Security

Report vulnerabilities privately — see [SECURITY.md](SECURITY.md) for the disclosure process.
