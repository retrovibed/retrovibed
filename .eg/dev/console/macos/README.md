# macOS Local Development Setup

## Prerequisites

### Xcode

Install Xcode from the App Store, then configure the command line tools:

```bash
sudo xcode-select --switch /Applications/Xcode.app/Contents/Developer
sudo xcodebuild -license accept
```

### EG Build Tool

See https://egdaemon.com/docs/installation/index.html

## Running the Setup

```bash
eg compute baremetal dev/macos
```

This command:

1. Installs Homebrew dependencies (go, duckdb, gpgme, flutter, ffmpeg@7, cocoapods)
2. Builds the native Go library (`retrovibed.dylib`) with the `localdev` tag
3. Generates Dart FFI bindings via `ffigen`
4. Runs `flutter pub get`

## Running the App

```bash
cd console
flutter run -d macos
```

## Rebuilding

Re-run the setup command after making changes to Go code in `shallows/` or `console/retrovibedbind/`:

```bash
eg compute baremetal dev/console/macos
```

Hot reload works for Dart changes but not for native library changes.

## Troubleshooting

### DuckDB Extension Loading

If you see "library load disallowed by system policy":

```bash
xattr -d com.apple.quarantine \
  ~/Library/Containers/com.example.retrovibed/Data/.duckdb/extensions/v1.4.1/osx_arm64/inet.duckdb_extension
```

### Downloads Folder Access

If the app crashes with "unable to watch: /Users/.../Downloads: operation not permitted", ensure the Downloads entitlement is set in both `macos/Runner/DebugProfile.entitlements` and `macos/Runner/Release.entitlements`:

```xml
<key>com.apple.security.files.downloads.read-write</key>
<true />
```

Rebuild after adding the entitlement.

Sometimes, the .dylib file may not be copied to the correct location. Can fix with:

```bash
cp [Path_to_retrovibed/console]/build/nativelib/retrovibed.dylib [Path_to_retrovibed/console]/build/macos/Build/Products/Debug/retrovibed.app/Contents/MacOS/retrovibed.dylib
```
