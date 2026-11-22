# Linux Local Development Setup

## Prerequisites

### EG Build Tool

See https://egdaemon.com/docs/installation/index.html

## Running the application

```bash
eg compute local --wayland --ttl=infinity dev/console/linux
```

The console runs as systemd user services (see `.eg/dev/.systemd`) inside the workload:

- saving a Dart file under `console/lib`, `console/hook` or `console/pubspec.yaml` runs `dart analyze` and `flutter analyze`, then hot reloads (`SIGUSR1`). Analyzer errors abort the reload.
- a change to a `.so` in `.eg.cache/dev.native.libs` restarts the app, a native library cannot be hot reloaded.

Logs stream in the eg session, or `journalctl --user -fu retro-console` inside the workload.

## Troubleshooting
