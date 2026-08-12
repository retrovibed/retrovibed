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
import 'api.dart' as remote;
import 'player.control.seek.dart';
import 'player.control.playpause.dart';
import 'playlist.current.dart';

class Connect extends StatefulWidget {
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
  State<Connect> createState() => _State();
}

class _State extends State<Connect> with LoadingState {
  final ValueNotifier<meta.Daemon> _endpoint = ValueNotifier(meta.Daemon());
  remote.RemoteControlSocket _socket = remote.RemoteControlSocket.noop;
  Stream<remote.Stream> _messages = Stream.empty();

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
    if (_endpoint.value.hostname.isEmpty) return;

    setState(() => loading = true);

    remote.remotecontrol
        .connect(host: _endpoint.value.hostname, options: [authn.request(authn.AuthzCache.meta(context))])
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

          return c.future;
        })
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
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((error) {
          setState(() {
            loading = false;
            _socket = remote.RemoteControlSocket.noop;
            cause = ds.Error.unknown(error, onTap: resetCause);
          });
        })
        .whenComplete(_reconnect);
  }

  @override
  void initState() {
    super.initState();
    _endpoint.value = meta.EndpointAuto.of(context)?.changed.value ?? _endpoint.value;
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
          SingleChildScrollView(
            child: ds.Container(
              padding: defaults.padding,
              ds.Loading(
                loading: _socket == remote.RemoteControlSocket.noop,
                Column(
                  verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
                  mainAxisSize: MainAxisSize.min,
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
                    SearchMinimal(
                      apisearch: (req, {host, options = const []}) =>
                          media.media.search(req, host: _endpoint.value.hostname, options: options),
                      empty: ds.Empty,
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
