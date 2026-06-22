# drop-cli Development Guide

## Commit Messages

Commit messages should follow this format:

**First line:** A high-level sentence that outlines the changes in lowercase

**Then:** Each change itemized with a `-` prefix

Example:
```
Refactor output styling and update version command display

- Capitalize output color constants for export
- Add RenderText helper for consistent styling
- Update version command with branded logo
```

## Code Style

### Output Styling
- Use `fatih/color` for terminal output styling
- Support 24-bit ANSI colors
- Brand color: `#10b981` (emerald green)

### General Go Style
- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and concise
- Add comments only for non-obvious logic

## Project Structure

- `cmd/` - CLI entry points
- `internal/api/` - API client and request handling
- `internal/app/` - Application logic and CLI commands
- `internal/config/` - Configuration management
- `bin/` - Build artifacts

## Building and Testing

See `Makefile` for build and test commands.
