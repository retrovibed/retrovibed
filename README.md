#### retrovibe

retrovibe is a personal digital archiving and distribution platform built designed to make digital distribution
user friendly and easy to manage. It provides the ability to manage and share content within a personal library
with the world and has an optional at cost cloud backup functionality.

#### install flatpak daemon

```bash
mkdir retrovibe
curl -L -o retrovibed.daemon.yml https://github.com/retrovibed/retrovibed/releases/download/${version}/flatpak.daemon.yml
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak-builder --user --install-deps-from=flathub --install --ccache --force-clean retrovibe retrovibed.daemon.yml
flatpak run --user space.retrovibe.Daemon
```

#### install flatpak gui

```bash
mkdir retrovibe
curl -L -o retrovibed.client.yml https://github.com/retrovibed/retrovibed/releases/download/${version}/flatpak.client.yml
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak-builder --user --install-deps-from=flathub --install --ccache --force-clean retrovibe retrovibed.client.yml
flatpak run --user space.retrovibe.Client
```

### install daemon from source

```bash
go install -tags no_duckdb_arrow github.com/retrovibed/retrovibed/shallows/cmd/retrovibed/...
```
