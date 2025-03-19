#### retrovibe

retrovibe is a personal digital archiving and distribution platform built designed to make digital distribution
user friendly and easy to manage. It provides the ability to manage and share content within a personal library
with the world and has an optional at cost cloud backup functionality.

#### install flatpak gui

```bash
mkdir retrovibe
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak-builder --user --install-deps-from=flathub --install --ccache --force-clean retrovibe .eg.cache/flatpak.client.yml
flatpak run --user space.retrovibe.Client
```

#### install flatpak daemon

```bash
mkdir retrovibe
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak-builder --user --install-deps-from=flathub --install --ccache --force-clean retrovibe .eg.cache/flatpak.daemon.yml
flatpak run --user space.retrovibe.Daemon
```

### install daemon from source

```bash

```
