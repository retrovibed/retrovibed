qq### Notice

this is public alpha software. its under active development and testing. expect issues.

#### download

[android](https://play.google.com/store/apps/details?id=space.retrovibe.retrovibed)  
[macos](https://github.com/retrovibed/retrovibed/releases/latest/download/retrovibed.dmg)

#### retrovibe

retrovibe is a personal digital archiving and distribution platform built designed to make digital distribution
user friendly and easy to manage. It provides the ability to manage and share content within a personal library
with the world and allows users to sign up for _at cost cloud storage functionality_.

It allow supporting communities via subscribing to content, like linux ISOs, AI models, and public interest archives that you can subscribe
to and donate your storage towards enabling distribution and resist censorship, all while remaining anonymous (not yet audited, best effort).

#### features

see the [site](https://retrovibe.space) for more details

- [x] a builtin userspace wireguard vpn. bring any wireguard vpn provider. allowing you to access your library from anywhere.
- [x] builtin media player, watch your personal music, video, images.
- [x] builtin bittorrent, share your personal media with whoever you want.
- [x] rss feeds for subscribing to content.
- [x] communities to reduce costs for sharable content.
- [x] integrated at cost archival service available that encrypts and offloads data from your devices.
- [x] monetarily support your favorite artists and communities directly by archiving their content.
- [x] distributed indexing/search for content discovery and exchange.
- [x] autoplay for music.
- [x] PVR functionality.
- [ ] recommendation engine.
- [ ] multi device support.

#### community sharing

build a community around content. each member reduces the cost for everyone.

#### install via appimage (recommended)

while you can just download the appimage file from [releases](https://github.com/retrovibed/retrovibed/releases/latest)

usage is simplified and auto-updates are provided via [AM](https://github.com/ivan-hc/AM)

```bash
am extra --user https://github.com/retrovibed/retrovibed retrovibed
# fix am management of the desktop integration.
am icons retrovibed
# update retrovibed
am update retrovibed
```

#### install via flatpak

generally not recommend at this time, requires flatpak-builder 1.4.2 or later to be installed.

```bash
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
curl -L -o space.retrovibe.Console.yml https://github.com/retrovibed/retrovibed/releases/latest/download/space.retrovibe.Console.yml
flatpak-builder --user --install-deps-from=flathub --install --ccache --force-clean retrovibe space.retrovibe.Console.yml
```

#### install deb daemon

```bash
add-apt-repository ppa:jljatone/retrovibed
apt-get update && apt-get install retrovibed

# /etc/retrovibed/config.env has documentation
# around the available settings.
cat | tee -a /etc/retrovibed/config.env << EOF
RETROVIBED_MDNS_DISABLED=true
RETROVIBED_TORRENT_AUTO_DISCOVERY=false
RETROVIBED_TORRENT_AUTO_BOOTSTRAP=true
RETROVIBED_TORRENT_PORT=10000
RETROVIBED_TORRENT_PUBLIC_IP4=""
RETROVIBED_TORRENT_PUBLIC_IP6=""
RETROVIBED_SELF_SIGNED_HOSTS=127.0.0.1
EOF

# generate an account. essentially used to create a static id for your account.
retrovibed identity generate {secret}

# authorize initial users using ssh keys. can be located by using `retrovibed identity show`
retrovibed identity bootstrap public-key "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBEdpDo/fUPKK7OUuZ4VM6JeBJmyZ882tQYPBN6nQwIk"
retrovibed identity bootstrap authorized-file /root/.ssh/authorized_keys

systemctl enable --now retrovibed.service
```

### determine ssh public key for client side

```bash
retrovibed identity show
```

### install daemon from source

```bash
go install github.com/retrovibed/retrovibed/shallows/cmd/retrovibe/...
```

### general commands

These are more reference commands for development / debugging.

#### generating media metadata archive

```bash
retrovibed media known tmdb --apikey="..." --start=1800-01-01 --end=2030-12-31 | retrovibed media known archive --directory="." --pattern="retrovibed.media.archive.d"
```

#### publishing torrents to an rss feed

```bash
# export the jwt secret so your client can connect. TODO: cleanup
RETROVIBED_JWT_SECRET="00000000-0000-0000-0000-000000000000"
# community names are globally unique. we reserve the right for change owners if someone is found squatting on a well known entity.
# we wont do it without informing the current owner 3 months in advance.
retrovibed community create --name="foo" --description="my special feed"

retrovibed torrent import directory --peer="localhost:9998" {directory} | retrovibed community publish --dry-run foo
# future work will allow using the exporting functionality to publish. either torrents or media.
# retrovibed torrent export "query" | retrovibed community publish --dry-run foo
```

#### generating a search plugin with an llm

search plugins are `wasip1`/`wasm` binaries that `retroapi/searchplugin.Registry` loads and runs sandboxed against the real network. They are significantly safer than other search systems in that the plugins are not given any access to the local network or the filesystem.

the API contract is small enough that an LLM can write a working plugin from a single reference file plus the name of the site to search:

give an LLM `https://github.com/retrovibed/retrovibed/blob/main/retroapi/examples/searchplugin-noop/main.go` as context — its comments walk through
the argv handling, the stdout contract, and why autohijack is required — along with the name of
the website you want it to search. that's enough to prompt something like:

```
Read https://github.com/retrovibed/retrovibed/blob/main/retroapi/examples/searchplugin-noop/main.go — that's the protocol a retrovibed search
plugin must satisfy. Write a new Go program in the same shape that searches <site> for
<query>, parsing real results and emitting one ddiscapi.Import per result to stdout instead
of the noop's fabricated one.
```

build and drop it in the well-known plugin directory, picked up live with no daemon restart:

```bash
# verify it natively first
go run . plugin --category all --query ubuntu

# then cross-compile for the registry
GOOS=wasip1 GOARCH=wasm go build -o myplugin.wasm .
cp myplugin.wasm ~/.config/retrovibed/search.d/
```

#### run flatpak from cli

```bash
flatpak run --user space.retrovibe.Console
flatpak run --command=sh --user space.retrovibe.Console # for debugging the runtime
```

#### manually moving a storage device.

sometimes you'll want to move what device your service is running on and if you dont have all your data in the archive you'll have to copy it.

here are the commands to do it:

```bash
# bootstrap your identity on the new device.
retrovibed identity generate {secret}
```

on the device you're exporting from:

```bash
retrovibed library export --no-torrent | ssh user@newdevicehost "~/go/bin/retrovibe library import jsonl"

retrovibed torrent export | ssh user@newdevicehost "retrovibed torrent import peer --peer='olddevicehost:port'"
```
