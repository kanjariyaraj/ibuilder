# 🚀 Builder: Step-by-Step Execution Guide

This guide provides a comprehensive, "one-by-one" walkthrough for setting up and running the **Builder** CLI for iOS delivery.

---

## 🛠 Step 1: Clone and Build from Source
First, get the code and build the executable binary.

```bash
# 1. Clone the repository
git clone https://github.com/kanjariyaraj/iBuilder.git
cd Builder

# 2. Build the binary (requires Go 1.25+)
make build

# 3. Verify the installation
./builder version
```

## 🔐 Step 2: Authenticate with GitHub
Builder uses GitHub Actions as its remote build engine. You need to link your GitHub account.

```bash
# Start the authentication flow
./builder auth github
```
*Follow the on-screen instructions to authorize via the browser. This saves a secure token locally.*

## 📁 Step 3: Initialize Your Project
Create the necessary configuration file (`builder.json`) in your iOS project root.

```bash
# Navigate to your iOS/Flutter/React Native project
cd /path/to/your/app

# Initialize the config
/path/to/Builder/builder config init
```

## 🔗 Step 4: Connect Your Repository
Link your local project directory to its GitHub remote repository.

```bash
# Detect and save repository owner/name
/path/to/Builder/builder repo connect

# Validate the connection and permissions
/path/to/Builder/builder repo validate
```

## 🩺 Step 5: Run the AI Doctor
Before building, ensure your environment is healthy and all dependencies are met.

```bash
./builder doctor
```
*The AI Doctor will analyze your system and provide suggestions if anything is missing.*

## 🏗 Step 6: Trigger a Remote Build
Kick off the iOS build process on a remote macOS runner.

```bash
# Run the build (uses settings from builder.json)
./builder build run --wait --logs
```
- `--wait`: Keeps the CLI open until the build finishes.
- `--logs`: Streams real-time logs from the macOS runner to your terminal.

## 📦 Step 7: Retrieve Artifacts
Once the build is successful, download the `.ipa` file.

```bash
# List available artifacts
./builder build artifacts

# Download the latest build
./builder build artifacts download
```

---

## 🤖 Automation & CI/CD
To run Builder automatically in a script or CI pipeline, you can use flags to bypass interactive prompts.

### Example Automation Script (`deploy.sh`)
```bash
#!/bin/bash
set -e

echo "Starting automated iOS build..."

# 1. Ensure config is valid
./builder config validate

# 2. Run build and wait for result
# Note: Ensure GITHUB_TOKEN is set in environment for non-interactive auth
./builder build run --wait --json > build_result.json

# 3. Check result
STATUS=$(jq -r '.status' build_result.json)
if [ "$STATUS" == "completed" ]; then
    echo "Build Success! Downloading artifact..."
    ./builder build artifacts download
else
    echo "Build Failed. Running AI diagnostics..."
    ./builder ai fix
fi
```

---

## 💡 Troubleshooting
If a build fails:
1. Run `./builder ai fix` to get an automated repair suggestion.
2. Run `./builder build log` to review the full output.
3. Check `docs/troubleshooting.md` for common issues.
