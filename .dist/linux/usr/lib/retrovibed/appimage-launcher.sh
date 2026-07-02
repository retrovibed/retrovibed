#!/bin/sh
# Launches the real retrovibed binary (retrovibed.bin, renamed from retrovibed by
# AppImageBuilder.yml's after_runtime hook) under bubblewrap when available. This
# script is installed in retrovibed's place at usr/lib/retrovibed/retrovibed by
# that same hook -- app_info.exec has to point at the real ELF binary for
# appimage-builder's own runtime setup to work, so the swap happens afterward.
#
# Why: GTK's icon-theme rendering routes SVG loading through glycin, a sandboxed
# helper spawned via a hardcoded absolute path (e.g.
# /usr/libexec/glycin-loaders/2+/glycin-svg) baked in at compile time. AppImage's
# mount doesn't remap /usr the way Flatpak's sandbox does, so that absolute lookup
# always hits the real host filesystem, which usually doesn't have glycin at that
# exact path (or at all) -- causing a hard abort the first time GTK needs any
# fallback icon.
#
# Fix: relaunch under bwrap with a synthetic /usr: every real top-level entry is
# rebound in read-only, except usr/libexec/glycin-loaders, which is grafted in
# from the bundle. bwrap is what Ubuntu's AppArmor unprivileged-userns hardening
# already trusts for this (Flatpak itself relies on it), unlike a raw
# unshare/mount --bind. bwrap itself is bundled (see apt.include in
# AppImageBuilder.yml), so this doesn't depend on the host having it installed.
HERE="$APPDIR"
REAL="$HERE/usr/lib/retrovibed/retrovibed.bin"
BWRAP="$HERE/usr/bin/bwrap"

USR_BINDS=""
for d in /usr/*/; do
  name=$(basename "$d")
  USR_BINDS="$USR_BINDS --ro-bind /usr/$name /usr/$name"
done

exec "$BWRAP" --dev-bind / / --tmpfs /usr $USR_BINDS \
  --bind "$HERE/usr/libexec/glycin-loaders" /usr/libexec/glycin-loaders \
  --chdir "$HERE/runtime/compat" \
  -- "$REAL" "$@"
