# Contributing

## How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Development Setup

See [docs/development.md](docs/development.md).

## Code Style

- Follow Go standard formatting (`gofmt`)
- Run `go vet` before committing
- Ensure all tests pass (`go test ./...`)
- No dead code, unused imports, or debug leftovers

## Pull Request Checklist

- [ ] Tests pass
- [ ] Code is formatted
- [ ] Documentation updated
- [ ] No breaking changes without discussion
