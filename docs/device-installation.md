# Device Installation

## Commands

### `builder device install`

Install an app on a connected device:

```bash
builder device install --ipa path/to/app.ipa
builder device install --artifact path/to/artifact.zip
```

### `builder device launch`

Launch an installed app:

```bash
builder device launch --bundle-id com.example.app
```

### Auto Install Flow

Build, download, install, and launch in one step:

```bash
builder ios build --install
```
