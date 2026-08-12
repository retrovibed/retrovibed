import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library/dropdown.upload.dart';
import 'package:retrovibed/library/search.mimetype.dropdown.dart';
import 'package:retrovibed/library/search.minimal.dart';
import 'api.dart' as remote;
import 'player.control.fastforward.dart';
import 'player.control.next.dart';
import 'player.control.playpause.dart';
import 'player.control.previous.dart';
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

class _State extends State<Connect> {
  ValueNotifier<meta.Daemon> _library = ValueNotifier(meta.Daemon());
  remote.RemoteControlSocket _socket = remote.RemoteControlSocket.noop;
  Stream<remote.Stream> _messages = Stream.empty();

  void refresh() {
    setState(() {}); // force rebuild
  }

  void _reconnect() {
    if (!mounted) return;
    print("reconnecting");
    setState(() => _socket = remote.RemoteControlSocket.noop);
    Future.delayed(const Duration(seconds: 2), _connect);
  }

  void _connect() {
    remote.remotecontrol
        .connect(options: [authn.request(authn.AuthzCache.meta(context))])
        .then((socket) {
          final c = Completer();
          setState(() {
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
        .catchError((cause) {
          debugPrint("remote control connect socket failed: ${cause}");
        })
        .whenComplete(_reconnect);
  }

  @override
  void initState() {
    super.initState();
    _library = meta.EndpointAuto.of(context)?.changed ?? _library;
    _library.addListener(refresh);
    ds.postframe(_connect);
  }

  @override
  void dispose() {
    super.dispose();
    _library.removeListener(refresh);
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
            library: _library,
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
                    const Divider(),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        PlayerControlPrevious(socket: _socket),
                        PlayerControlPlayPause(socket: _socket, status: _messages),
                        PlayerControlFastForward(socket: _socket),
                        PlayerControlNext(socket: _socket),
                      ],
                    ),
                    PlaylistCurrent(_messages),
                    SearchMinimal(
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
