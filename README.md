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
| `builder auth github` | Authenticate with GitHub via device flow |
| `builder auth status` | Check GitHub authentication status |
| `builder auth logout` | Remove GitHub authentication |
| `builder repo connect` | Connect repository from git remote |
| `builder repo info` | Display repository metadata |
| `builder repo validate` | Validate repository access |
| `builder build run` | Trigger an iOS build workflow |
| `builder build status` | Check build status |
| `builder build list` | List recent builds |
| `builder build log` | View build logs |
| `builder build artifacts` | List build artifacts |
| `builder init run` | Initialize project and generate config/workflow |
| `builder ios build` | Build iOS app remotely via GitHub Actions |

## Roadmap

- **Phase 1** — Foundation & CLI framework (completed)
- **Phase 2** — GitHub authentication & repository management (completed)
- **Phase 3** — iOS build pipeline integration (completed)
- **Phase 4** — Initialization engine & workflow generator (completed)
- **Phase 5** — Remote iOS build engine (current)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## License

MIT - see [LICENSE](LICENSE) for details.
