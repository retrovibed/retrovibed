import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/stateful.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library/dropdown.upload.dart';
import 'package:retrovibed/library/search.mimetype.dropdown.dart';
import 'package:retrovibed/library/search.minimal.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'api.dart' as remote;
import 'player.control.seek.dart';
import 'player.control.playpause.dart';
import 'playlist.current.dart';

// Mirrors media.PlayAction's shape but queues the media on the connected
// remote daemon's playlist instead of this device's local Playlist, since
// SearchMinimal here searches the app's own library to feed *that* daemon.
Future<void> Function()? Function(BuildContext, media.Media, media.MediaSearchResponse) RemotePlayAction(
  remote.RemoteControlSocket socket,
) {
  return (context, current, s) {
    switch (mimex.icon(current.mimetype)) {
      case mimex.icomovie:
      case mimex.icoaudio:
        return () async {
          socket.send(
            remote.Stream(
              sid: uuidx.v7(),
              queue: remote.Queue(media: current),
            ),
          );
        };
      default:
        return null;
    }
  };
}

// Public entrypoint: wraps _Connect in an authn.AuthedEndpoint so it can target a
// user-selected remote daemon (via its DaemonDropdown) with a matching
// scoped auth token, independent of the app-root EndpointAuto/AuthzCache.
class Connect extends StatelessWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final media.FnUploadRequest apiupload;
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;
  final Widget downloading;
  final void Function(Widget) onDownloadingChanged;

  const Connect({
    super.key,
    required this.search,
    this.apiupload = media.media.upload,
    required this.mode,
    required this.onModeChanged,
    required this.downloading,
    required this.onDownloadingChanged,
  });

  @override
  Widget build(BuildContext context) {
    return authn.AuthedEndpoint(
      _Connect(
        search: search,
        apiupload: apiupload,
        mode: mode,
        onModeChanged: onModeChanged,
        downloading: downloading,
        onDownloadingChanged: onDownloadingChanged,
      ),
    );
  }
}

class _Connect extends StatefulWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final media.FnUploadRequest apiupload;
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;
  final Widget downloading;
  final void Function(Widget) onDownloadingChanged;

  const _Connect({
    required this.search,
    this.apiupload = media.media.upload,
    required this.mode,
    required this.onModeChanged,
    required this.downloading,
    required this.onDownloadingChanged,
  });

  @override
  State<_Connect> createState() => _State();
}

class _State extends State<_Connect> with LoadingState {
  remote.RemoteControlSocket _socket = remote.RemoteControlSocket.noop;
  Stream<remote.Stream> _messages = Stream.empty();

  ValueNotifier<meta.Daemon> get _endpoint => authn.AuthedEndpoint.daemon(context);

  void _onEndpointChanged() {
    _socket.close();
    setState(() => _socket = remote.RemoteControlSocket.noop);
    _connect();
  }

  void _reconnect() {
    if (!mounted) return;
    print("reconnecting");
    setState(() => _socket = remote.RemoteControlSocket.noop);
    Future.delayed(const Duration(seconds: 2), _connect);
  }

  void _connect() {
    if (!mounted) return;
    if (_endpoint.value.hostname.isEmpty) return;

    setState(() => loading = true);

    remote.remotecontrol
        .connect(host: _endpoint.value.hostname, options: [authn.request(authn.AuthedEndpoint.token(context))])
        .then((socket) {
          final c = Completer();
          setState(() {
            loading = false;
            cause = ds.Error.zero;
            _socket = socket;
            _messages = socket.messages.asBroadcastStream();
            _messages.listen(
              (_) {},
              cancelOnError: true,
              onError: c.completeError,
              onDone: c.complete,
            );
          });

          socket.send(remote.Stream(sid: uuidx.v7(), sync: remote.Sync()));

          return c.future;
        })
        .then((_) => _reconnect())
        .catchError((error) {
          setState(() {
            loading = false;
            _socket = remote.RemoteControlSocket.noop;
            cause = ds.Error.unauthorized(
              error,
              message: const Text("remote control is disabled on this daemon"),
              onTap: resetCause,
            );
          });
        }, test: httpx.ErrorsTest.forbidden)
        .catchError((error) {
          setState(() {
            loading = false;
            _socket = remote.RemoteControlSocket.noop;
            cause = ds.Errors.httpauto(error, onTap: resetCause);
          });
          _reconnect();
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((error) {
          setState(() {
            loading = false;
            _socket = remote.RemoteControlSocket.noop;
            cause = ds.Error.unknown(error, onTap: resetCause);
          });
          _reconnect();
        });
  }

  @override
  void initState() {
    super.initState();
    _endpoint.addListener(_onEndpointChanged);
    ds.postframe(_connect);
  }

  @override
  void dispose() {
    super.dispose();
    _endpoint.removeListener(_onEndpointChanged);
    _socket.close();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Container(
      Column(
        verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
        children: [
          meta.DaemonDropdown(
            library: _endpoint,
            onSelect: meta.DaemonDropdown.local,
            remoteonly: true,
            readonly: true,
            leading: [
              ValueListenableBuilder<media.MediaSearchState>(
                valueListenable: widget.search,
                builder: (context, state, _) => ds.CompactingMenu.pinned(
                  DropdownUpload(
                    icon: SearchMimetypeDropdown.icon(mimex.checksum(state.next.mimetypes)),
                    items: [
                      ...SearchMimetypeDropdown.menuItems(widget.search),
                      media.SearchModeToggle(
                        mode: media.SearchMode.discovery,
                        current: widget.mode,
                        icon: Icons.travel_explore,
                        label: "Discover",
                        onSelect: widget.onModeChanged,
                      ),
                      media.SearchModeToggle(
                        mode: media.SearchMode.remote,
                        current: widget.mode,
                        icon: Icons.settings_remote,
                        label: "Remote",
                        onSelect: widget.onModeChanged,
                      ),
                      const PopupMenuDivider(),
                      PopupMenuItem<String>(
                        enabled: false,
                        child: ValueListenableBuilder<media.MediaSearchState>(
                          valueListenable: widget.search,
                          builder: (context, s, _) => mimex.CategoryOptionsLabel(s.next.mimetypes),
                        ),
                      ),
                      media.MenuItemUploadFiles(
                        context,
                        widget.search,
                        apiupload: widget.apiupload,
                      ),
                      downloads.MenuItemDownloadTorrent(context, (downloads) {
                        widget.onDownloadingChanged(
                          media.DownloadQueue(
                            downloads,
                            onQueueComplete: () => widget.onDownloadingChanged(ds.Empty),
                          ),
                        );
                        print("downloading torrents ${downloads}");
                      }),
                      downloads.MenuItemDownloadMagnet(context, (downloads) {
                        widget.onDownloadingChanged(
                          media.DownloadQueue(
                            downloads,
                            onQueueComplete: () => widget.onDownloadingChanged(ds.Empty),
                          ),
                        );
                        print("downloading magnets ${downloads}");
                      }),
                    ],
                  ),
                ),
              ),
            ],
          ),
          Expanded(
            child: ds.Container(
              padding: defaults.padding,
              ds.Loading(
                loading: _socket == remote.RemoteControlSocket.noop,
                Column(
                  verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        PlayerControlSeek.prev(socket: _socket),
                        PlayerControlSeek.backward(socket: _socket),
                        PlayerControlPlayPause(socket: _socket, status: _messages),
                        PlayerControlSeek.forward(socket: _socket),
                        PlayerControlSeek.next(socket: _socket),
                      ],
                    ),
                    PlaylistCurrent(_messages),
                    Expanded(
                      child: SearchMinimal(
                        empty: ds.Empty,
                        onPlay: RemotePlayAction(_socket),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
