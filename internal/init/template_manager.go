package init

import (
	"fmt"
)

func getTemplate(pt ProjectType, projectName string) string {
	switch pt {
	case ProjectFlutter:
		return flutterTemplate(projectName)
	case ProjectReactNative:
		return reactNativeTemplate(projectName)
	case ProjectExpo:
		return expoTemplate(projectName)
	case ProjectCapacitor:
		return capacitorTemplate(projectName)
	case ProjectNativeiOS:
		return nativeIOSTemplate(projectName)
	case ProjectUnity:
		return unityTemplate(projectName)
	case ProjectUnreal:
		return unrealTemplate(projectName)
	default:
		return nativeIOSTemplate(projectName)
	}
}

func flutterTemplate(name string) string {
	return fmt.Sprintf(`name: %s iOS Build

on:
  workflow_dispatch:
    inputs:
      build_mode:
        description: "Build mode"
        default: "release"
        type: choice
        options:
          - debug
          - release

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: subosito/flutter-action@v2
        with:
          channel: stable
      - run: flutter pub get
      - run: flutter build ios --no-codesign --${{ github.event.inputs.build_mode }}
      - uses: actions/upload-artifact@v4
        with:
          name: ios-build
          path: build/ios/
`, name)
}

func reactNativeTemplate(name string) string {
	return fmt.Sprintf(`name: %s iOS Build

on:
  workflow_dispatch:
    inputs:
      build_mode:
        description: "Build mode"
        default: "release"
        type: choice
        options:
          - debug
          - release

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm install
      - run: npx react-native bundle --platform ios --dev false --entry-file index.js --bundle-output ios/main.jsbundle
      - uses: actions/upload-artifact@v4
        with:
          name: ios-build
          path: ios/
`, name)
}

func expoTemplate(name string) string {
	return fmt.Sprintf(`name: %s iOS Build

on:
  workflow_dispatch:
    inputs:
      profile:
        description: "Build profile"
        default: "production"
        type: string

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm install
      - run: npx expo build:ios
      - uses: actions/upload-artifact@v4
        with:
          name: ios-build
          path: dist/
`, name)
}

func capacitorTemplate(name string) string {
	return fmt.Sprintf(`name: %s iOS Build

on:
  workflow_dispatch:
    inputs:
      build_mode:
        description: "Build mode"
        default: "release"
        type: choice
        options:
          - debug
          - release

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npm install
      - run: npx cap sync ios
      - run: npx cap build ios --${{ github.event.inputs.build_mode }}
      - uses: actions/upload-artifact@v4
        with:
          name: ios-build
          path: ios/
`, name)
}

func nativeIOSTemplate(name string) string {
	return fmt.Sprintf(`name: %s iOS Build

on:
  workflow_dispatch:
    inputs:
      scheme:
        description: "Xcode scheme"
        required: true
        type: string
      configuration:
        description: "Build configuration"
        default: "Release"
        type: choice
        options:
          - Debug
          - Release

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build
        run: |
          xcodebuild -scheme "${{ github.event.inputs.scheme }}" \
                     -configuration "${{ github.event.inputs.configuration }}" \
                     -destination "generic/platform=iOS" \
                     clean build
      - uses: actions/upload-artifact@v4
        with:
          name: ios-build
          path: build/
`, name)
}

func unityTemplate(name string) string {
	return fmt.Sprintf(`name: %s iOS Build

on:
  workflow_dispatch: {}

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - uses: game-ci/unity-builder@v4
        with:
          targetPlatform: iOS
      - uses: actions/upload-artifact@v4
        with:
          name: ios-build
          path: build/
`, name)
}

func unrealTemplate(name string) string {
	return fmt.Sprintf(`name: %s iOS Build

on:
  workflow_dispatch: {}

jobs:
  build:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build iOS
        run: |
          echo "Unreal Engine iOS build requires manual setup"
          echo "See: https://docs.unrealengine.com/SharingAndReleasing/iOS
      - uses: actions/upload-artifact@v4
        with:
          name: ios-build
          path: build/
`, name)
}
