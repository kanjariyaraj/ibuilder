# 🚀 Builder: Industrial-Grade iOS Delivery from Anywhere

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](go.mod)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20WSL-blue)](docs/DESIGN.md)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)](.github/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kanjariyaraj/Builder)](https://goreportcard.com/report/github.com/kanjariyaraj/Builder)

**Builder** is an industrial-grade open-source Go CLI that democratizes iOS development. It eliminates the "Mac Tax" by enabling developers on **Windows, Linux, and WSL** to build, test, sign, and release native iOS applications using remote macOS infrastructure and local AI-driven diagnostics.

---

## 📖 Table of Contents

- [🏛️ Why Builder?](#️-why-builder)
- [🧠 Core Technical Pillars](#-core-technical-pillars)
- [🏗️ Architectural Overview](#️-architectural-overview)
- [📦 Installation](#-installation)
- [🚦 Quick Start](#-quick-start)
- [🛠️ Command Reference](#️-command-reference)
- [📂 Documentation Deep-Dive](#-documentation-deep-dive)
- [🗺️ Roadmap](#️-roadmap)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

---

## 🏛️ Why Builder?

iOS development has traditionally been locked to macOS hardware. **Builder** breaks this barrier by orchestrating professional-grade pipelines:

- **Zero macOS Requirement**: Build native `.ipa` files using GitHub-hosted macOS runners.
- **AI-Powered Diagnostics**: Integrated "AI Doctor" analyzes thousands of lines of build logs to find and fix errors automatically.
- **Unified Workflow**: One tool for Xcode, Flutter, and React Native projects.
- **Release Ready**: Direct integration with TestFlight and App Store Connect.

---

## 🧠 Core Technical Pillars

### 1. Remote Build Engine (`internal/build`)
Orchestrates complex iOS build pipelines on remote macOS environments. It handles workflow dispatch, real-time log streaming, and secure artifact retrieval (IPA, dSYM) back to your local machine.
> [Read more about the Build Engine](docs/build-engine.md)

### 2. AI Doctor & Diagnostics (`internal/ai`)
More than a CLI, Builder is an intelligent assistant. It uses a specialized knowledge base to:
- **Analyze Failures**: Identifies root causes in build logs (provisioning, code-signing, dependency conflicts).
- **Suggest Fixes**: Provides actionable steps or applies "Auto-Fixes" to project configurations.
> [Read more about AI Doctor](docs/ai-doctor.md)

### 3. MobAI Remote Protocol (`internal/mobai`)
A high-performance TCP-based interaction layer that allows you to communicate with remote devices or simulator agents, enabling remote hot-reload and debugging sessions.
> [Read more about MobAI](docs/mobai.md)

---

## 🏗️ Architectural Overview

Builder follows a **Clean Architecture** pattern, ensuring high maintainability and testability.

```text
       USER (CLI)
          │
    ┌─────▼────────────────┐
    │  cmd/builder/        │◄── Entry Point & Routing
    └─────┬────────────────┘
          │
    ┌─────▼────────────────┐
    │  internal/domain/    │◄── Business Logic & Core Engines
    │  (build, ai, mobai)  │
    └─────┬────────────────┘
          │
    ┌─────▼────────────────┐
    │  External Services   │◄── GitHub Actions, Apple Connect,
    │                      │    Real Devices via MobAI
    └──────────────────────┘
```

---

## 📦 Installation

### Via Go (Recommended)
```bash
go install github.com/kanjariyaraj/Builder/cmd/builder@latest
```

### From Source
```bash
git clone https://github.com/kanjariyaraj/Builder.git
cd Builder
make build
./builder version
```

---

## 🚦 Quick Start

Experience the power of Builder in 5 steps:

1. **Initialize**: `builder config init` (Creates `builder.json`)
2. **Authenticate**: `builder auth github` (Connect your remote build account)
3. **Connect**: `builder repo connect` (Link your project repository)
4. **Doctor Check**: `builder doctor` (Verify your environment is ready)
5. **Build**: `builder build run` (Trigger your first remote iOS build)

---

## 🖥️ CLI Experience

Builder provides a rich, color-coded CLI experience. For example, `builder doctor` provides an instant health report:

```text
$ builder doctor
✓ Git installed (2.43.0)
✓ Go installed (1.25.0)
✓ builder.json found and valid
⚠ Flutter not found (Optional for Xcode-only projects)
✗ Not authenticated with GitHub

Summary: 3 Healthy, 1 Warning, 1 Failure
Run 'builder auth github' to fix the failure.
```

---

## 🛠️ Command Reference

| Category | Command | Description |
|:---|:---|:---|
| **System** | `builder doctor` | Comprehensive system & dependency audit |
| | `builder version` | Print detailed version information |
| **Auth** | `builder auth github` | Authenticate with GitHub (Device Flow) |
| | `builder auth status` | Check authentication & token health |
| **Config** | `builder config init` | Generate an optimized `builder.json` |
| | `builder config validate` | Deep-check configuration integrity |
| **Repo** | `builder repo connect` | Connect local project to remote builder |
| | `builder repo info` | Inspect remote repository metadata |
| **Build** | `builder build run` | Trigger a remote iOS build workflow |
| | `builder build log` | Stream real-time logs from the macOS runner |
| | `builder build artifacts` | Download IPAs and build reports |
| **AI** | `builder ai fix` | Attempt to auto-repair project issues |

---

## 📂 Documentation Deep-Dive

- [**Getting Started**](docs/getting-started.md) — Comprehensive setup guide.
- [**Configuration**](docs/configuration.md) — Detailed field-by-field breakdown of `builder.json`.
- [**Architecture**](docs/DESIGN.md) — Deep dive into the project's internal design.
- [**Development**](docs/development.md) — How to contribute and set up your dev environment.
- [**Troubleshooting**](docs/troubleshooting.md) — Common issues and their AI-suggested fixes.

---

## 🗺️ Roadmap

- **Phase 1-4**: Foundation, Auth, and Workflow Generation (✅ Completed)
- **Phase 5-10**: Build Engine, Artifacts, and AI Diagnostics (✅ Completed)
- **Phase 11-12**: TestFlight Integration & Release Management (✅ Completed)
- **Phase 13+**: App Store Connect Metadata & Automated Screenshotting (🚧 In Progress)

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for our "Research-First" development process.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feat/amazing-feature`)
3. Commit your Changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the Branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

---
<p align="center">
  Built with ❤️ by <a href="https://github.com/kanjariyaraj">Kanjariya Raj</a>
</p>
