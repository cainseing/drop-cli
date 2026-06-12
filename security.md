# Security Policy

## Supported versions

Only the latest release is supported. Update before reporting a vulnerability.

| Version | Supported |
|---------|-----------|
| v0.5.x  | Yes       |
| < v0.5  | No        |

---

## Security model

Drop is built on a zero-knowledge architecture. The server stores only an opaque encrypted blob—it has no access to the plaintext or the decryption key at any point.

### Encryption — AES-256-GCM

Secrets are encrypted on the client before upload using AES-256 in Galois/Counter Mode.

- A fresh 256-bit key is generated per drop using `crypto/rand`.
- A unique random nonce is generated for every encryption operation.
- GCM provides authenticated encryption. Any server-side tampering causes decryption to fail locally with an authentication error.

### Envelope and padding

Before encryption, the plaintext is wrapped in a versioned envelope:

```
[ 4-byte length (big-endian) ][ plaintext ][ random padding ]
```

The minimum ciphertext size is 128 bytes. Padding is filled with random bytes from `crypto/rand`, ensuring payload sizes do not leak information about secret length.

### Token structure

The decryption key is embedded in the token and never transmitted to or stored by the server:

```
drop_<base64url( protocol_version "." blob_id "." hex_key ":" signed_flag )>
```

The server generates and returns the `blob_id` after the blob is uploaded. The client embeds it in the token locally. The server stores the `blob_id` and the ciphertext; the decryption key exists only in the token. Whoever holds the token holds the key.

The `signed_flag` (`1` or `0`) records, at creation time, whether the sender signed the drop. Because the token is shared out-of-band and never seen by the server, this flag cannot be tampered with server-side — it sets the recipient's expectation independently of whatever the server reports.

### SSH signing (optional)

Drops can be cryptographically signed using your local SSH private key (`id_ed25519`, `id_rsa`, or `id_ecdsa`). The signature is verified against public keys fetched from your GitHub or GitLab profile at retrieval time.

- Signing uses the raw ciphertext as the payload, not the plaintext.
- Verification confirms both the sender's identity and that the ciphertext has not been modified since it was signed.
- If signature verification fails, the retrieved secret is discarded and an error is shown to the recipient.
- If the token's `signed_flag` indicates the drop was signed but the server response is missing the signature, sender, or provider, the secret is discarded and an error is shown — a server cannot downgrade a signed drop to unsigned without detection.

### What the server holds

| Data | Stored server-side |
|------|--------------------|
| Encrypted blob | Yes |
| Blob ID | Yes |
| Signature (if signed) | Yes |
| Sender profile (if signed) | Yes |
| Decryption key | **No** |
| Plaintext | **No** |

---

## Reporting a vulnerability

Do not open a public GitHub issue for security vulnerabilities.

1. **Email** `security@getdrop.dev` with a detailed report.
2. **Include** a proof-of-concept, the CLI version (`drop version`), and your OS.
3. **Response** — acknowledgement within 48 hours.

We follow coordinated disclosure. Please do not share vulnerability details publicly until a fix has been released.
