# iOS Local Development Setup

## Prerequisites

### Xcode

Install Xcode from the App Store, then configure the command line tools:

```bash
sudo xcode-select --switch /Applications/Xcode.app/Contents/Developer
sudo xcodebuild -license accept
```

### EG Build Tool

See https://egdaemon.com/docs/installation/index.html

## Running the App

```bash
eg compute baremetal --no-podman dev/console/ios --invalidate-cache
```

This command:

1. Installs Homebrew dependencies (go, duckdb, gpgme, flutter, ffmpeg@7, cocoapods, cmake, ninja)
2. Clones and builds DuckDB as a static library for the iOS simulator (x86_64)
3. Builds the native Go library (`libretrovibed.a`) with `c-archive` buildmode and `localdev` tag
4. Generates Dart FFI bindings via `ffigen`
5. Patches Mach-O platform metadata via `vtool` for simulator compatibility
6. Merges all DuckDB static libraries via `libtool`
7. Installs CocoaPods and launches the app on the iOS simulator
