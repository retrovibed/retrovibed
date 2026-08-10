import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/media/play.list.history.dart';
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library/dropdown.upload.dart';
import 'package:retrovibed/library/search.mimetype.dropdown.dart';
import 'player.control.fastforward.dart';
import 'player.control.next.dart';
import 'player.control.playpause.dart';
import 'player.control.previous.dart';
import 'playlist.current.dart';

class Display extends StatefulWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final media.FnUploadRequest apiupload;
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;
  final Widget downloading;
  final void Function(Widget) onDownloadingChanged;

  const Display({
    super.key,
    required this.search,
    this.apiupload = media.media.upload,
    required this.mode,
    required this.onModeChanged,
    required this.downloading,
    required this.onDownloadingChanged,
  });

  @override
  State<Display> createState() => _State();
}

class _State extends State<Display> {
  ValueNotifier<meta.Daemon> _library = ValueNotifier(meta.Daemon());

  void refresh() {
    setState(() {}); // force rebuild
  }

  @override
  void initState() {
    super.initState();
    _library = meta.EndpointAuto.of(context)?.changed ?? _library;
    _library.addListener(refresh);
  }

  @override
  void dispose() {
    super.dispose();
    _library.removeListener(refresh);
  }

  @override
  Widget build(BuildContext context) {
    final playlist = media.Playlist.of(context);
    final defaults = ds.Defaults.of(context);
    return ds.Container(
      Column(
        verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
        mainAxisSize: MainAxisSize.min,
        children: [
          meta.DaemonDropdown(
            library: _library,
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
          ds.Container(
            padding: defaults.padding,
            Column(
              verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
              mainAxisSize: MainAxisSize.min,
              children: [
                if (playlist != null) ...[
                  const Divider(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      PlayerControlPrevious(),
                      PlayerControlPlayPause(),
                      PlayerControlFastForward(),
                      PlayerControlNext(),
                    ],
                  ),
                  PlaylistCurrent(playlist.queue.current),
                  PlayListHistory(playlist.queue.recent.take(5).toList()),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }
}
