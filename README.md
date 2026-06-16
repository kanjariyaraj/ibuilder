# Builder

Build, Test, Sign and Release iOS Apps From Anywhere.

Builder is an open-source Go CLI that allows developers to build iOS apps from Windows, Linux and WSL using GitHub-hosted macOS runners and real iPhone devices.

## Installation

```bash
go install github.com/kanjariyaraj/Builder/cmd/builder@latest
```

Or build from source:

```bash
git clone https://github.com/kanjariyaraj/Builder.git
cd Builder
go build -o builder ./cmd/builder/
```

## Commands

| Command          | Description                          |
|------------------|--------------------------------------|
| `builder`        | Display general information          |
| `builder version`| Print version information            |
| `builder doctor` | Check system dependencies            |
| `builder config` | Manage Builder configuration         |
| `builder config init` | Create default configuration    |
| `builder config show` | Display configuration          |
| `builder config validate` | Validate configuration     |
| `builder --help` | Display help information             |

## Roadmap

- **Phase 1** — Foundation & CLI framework (current)
- **Phase 2** — iOS build pipeline integration
- **Phase 3** — Signing & code signing management
- **Phase 4** — Release automation & GitHub integration
- **Phase 5** — Multi-platform support & testing

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## License

MIT - see [LICENSE](LICENSE) for details.
