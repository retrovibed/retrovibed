#### install flatpak daemon

```bash
mkdir retrovibe
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak-builder --user --install-deps-from=flathub --install --ccache --force-clean retrovibe .eg.cache/flatpak.daemon.yml
flatpak run --user space.retrovibe.Daemon
```

#### install flatpak gui

```bash
mkdir retrovibe
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak-builder --user --install-deps-from=flathub --install --ccache --force-clean retrovibe .eg.cache/flatpak.client.yml
flatpak run --user space.retrovibe.Client
```

### install daemon from source

```bash

```
