# Security Policy

## Security Philosophy
"Drop" is built on a **Zero-Knowledge** architecture. The core principle is that the server should never possess the ability to decrypt user data. Security is enforced at the edge (the CLI) rather than the infrastructure.

## Supported Versions
Only the latest version of the Drop CLI is supported for security updates. Please ensure you are running the most recent release before reporting an issue.

| Version | Supported          |
| ------- | ------------------ |
| v0.1.x  | ✅ Supported       |
| < v0.1  | ❌ Not Supported   |

## Security Architecture

### 1. Encryption (AES-256-GCM)
Secrets are encrypted locally using AES-256 in Galois/Counter Mode (GCM). 
* **Key Generation:** A new key is generated for every single drop.
* **Nonce:** A unique nonce is used for every encryption operation.
* **Integrity:** GCM provides authenticated encryption, ensuring that if a Drop is tampered with on the server, decryption will fail locally.

### 2. Envelope Protocol
Data is wrapped in a versioned envelope before encryption:
`[Data Length (4 bytes)] [Plaintext Payload] [Random Padding]`

### 3. Zero-Knowledge Proof
The decryption key is appended to the `drop_` token locally. This token (containing the Protocol Version, ID and the Key) is never stored in its entirety on the server. The server only stores the encrypted blob and the ID.

## Reporting a Vulnerability
If you discover a security vulnerability, please do not open a public issue. Instead, follow these steps:

1. **Email:** Send a detailed report to security@getdrop.dev.
2. **Details:** Include a proof-of-concept, the version of the CLI, and your OS.
3. **Response:** Acknowledgement of your report will typically occur within 48 hours.

We follow a coordinated disclosure model. We ask that you do not share details about the vulnerability until a fix has been published.