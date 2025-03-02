#### release command (local system)

```bash
GH_TOKEN="$(gh auth token)" eg compute local --hotswap -vv -e GH_TOKEN
```

#### install flatpak

```bash
mkdir derp; cd derp
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak-builder --user --install-deps-from=flathub --install --force-clean derp .eg.cache/flatpak.manifest.yml
flatpak run --user space.retrovibe.Daemon
```
