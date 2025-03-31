#### retrovibe

retrovibe is a personal digital archiving and distribution platform built designed to make digital distribution
user friendly and easy to manage. It provides the ability to manage and share content within a personal library
with the world and has allows users to sign up for at cost cloud backup functionality.

#### install deb daemon

```bash
add-apt-repository ppa:jljatone/retrovibed
apt-get update && apt-get install retrovibed

# /etc/retrovibed/config.env has documentation
# around the available settings.
cat | tee -a /etc/retrovibed/config.env << EOF
RETROVIBED_MDNS_DISABLED=true
RETROVIBED_TORRENT_AUTO_DISCOVERY=false
RETROVIBED_TORRENT_AUTO_BOOTSTRAP=false
RETROVIBED_SELF_SIGNED_HOSTS=127.0.0.1
EOF

systemctl enable --now retrovibed.service
```

#### install flatpak daemon

```bash
mkdir retrovibe
curl -L -o retrovibed.daemon.yml https://github.com/retrovibed/retrovibed/releases/latest/download/flatpak.daemon.yml
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak-builder --user --install-deps-from=flathub --install --ccache --force-clean retrovibe retrovibed.daemon.yml
flatpak run --user space.retrovibe.Daemon
```

### install daemon from source

```bash
go install github.com/retrovibed/retrovibed/shallows/cmd/retrovibed/...
```

#### install flatpak gui

```bash
mkdir retrovibe
curl -L -o retrovibed.client.yml https://github.com/retrovibed/retrovibed/releases/latest/download/flatpak.client.yml
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak-builder --user --install-deps-from=flathub --install --ccache --force-clean retrovibe retrovibed.client.yml
flatpak run --user space.retrovibe.Client
flatpak run --command=sh --user space.retrovibe.Client
```
