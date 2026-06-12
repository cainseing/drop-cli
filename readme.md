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

### Create a secret

Pass a value directly:

```sh
drop "my_api_key"
```

Or pipe from stdin:

```sh
cat .env | drop
```

Both print a `drop_…` token to stdout.

### Retrieve a secret

```sh
drop get <token>
```

Pipe directly into a file:

```sh
drop get <token> > .env
```

As a shorthand, you can pass the full token (with the `drop_` prefix) as the first argument and `get` is implied:

```sh
drop drop_<token>
```

### Purge a secret

Destroy a drop before it is read:

```sh
drop purge <token>
```

### Signed drops

Sign a drop with your SSH key so recipients can verify who created it:

```sh
drop identity set github <username>   # or gitlab
drop -s "secret value"
```

The CLI reads your local SSH private key (`id_ed25519`, `id_rsa`, or `id_ecdsa`) and verifies it matches your public keys on GitHub/GitLab. When the recipient fetches a signed drop, the signature is verified automatically.

### Version

```sh
drop version
```

Prints the running version and checks for updates. Updates are also checked automatically once per day when running any command.

---

## Flags

These flags apply to the create command (`drop [secret]`):

| Flag | Long form | Default | Description |
|------|-----------|---------|-------------|
| `-t` | `--ttl`   | `5`     | Expiry in minutes (max 10080 = 7 days) |
| `-r` | `--reads` | `1`     | Maximum number of reads before the secret is destroyed |
| `-s` | `--signed`| `false` | Sign the drop with your local SSH key |
| `-c` | `--copy`  | `false` | Copy the token to clipboard after creation |

**Example — share a secret that expires in 2 hours and allows 3 reads:**

```sh
drop -t 120 -r 3 "temporary_password"
```

---

## Limits

| Parameter | Limit |
|-----------|-------|
| Payload size | 1 MB |
| TTL | 1 minute – 7 days |
| Reads | 1 – unlimited |

---

## How it works

1. The secret is encrypted locally with **AES-256-GCM** using a randomly generated key.
2. The plaintext is wrapped in an envelope (`[4-byte length][plaintext][random padding]`) before encryption to obfuscate payload size.
3. The encrypted blob is uploaded to the server. The server has no access to the key.
4. A token is generated locally:

   ```
   drop_<base64(protocol_version.blob_id.hex_key)>
   ```

   The encryption key is embedded in the token and never stored server-side.
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
