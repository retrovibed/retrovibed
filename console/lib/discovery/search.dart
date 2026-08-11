import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/library/dropdown.upload.dart';
import 'package:retrovibed/library/search.mimetype.dropdown.dart';
import 'grid.dart';
import 'search.button.dart';

class Search extends StatefulWidget {
  final media.FnUploadRequest apiupload;
  final TextEditingController? controller;
  final FocusNode? focus;
  final ValueNotifier<media.MediaSearchState> search;
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;
  final Widget downloading;
  final void Function(Widget) onDownloadingChanged;

  const Search({
    super.key,
    this.apiupload = media.media.upload,
    this.controller,
    this.focus,
    required this.search,
    required this.mode,
    required this.onModeChanged,
    required this.downloading,
    required this.onDownloadingChanged,
  });

  @override
  State<Search> createState() => _SearchState();
}

class _SearchState extends State<Search> {
  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return Column(
      verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
      children: [
        ValueListenableBuilder<media.MediaSearchState>(
          valueListenable: widget.search,
          builder: (context, state, _) => ds.SearchTray(
            autoscroll: true,
            focus: widget.focus,
            autofocus: defaults.desktop,
            decoration: InputDecoration(hintText: "discover content across the network"),
            controller: widget.controller,
            trailing: [
              SearchButton(search: state, label: ds.Empty),
            ],
            onSubmitted: (v) {
              final freshNext = widget.search.value.next.clone()
                ..query = v
                ..offset = ds.Grid.int64(0);
              widget.search.value = media.MediaSearchState(
                next: freshNext,
                count: widget.search.value.count,
              );
              widget.focus?.requestFocus();
              ds.textediting.refocus(widget.controller);
              return Future.value();
            },
            next: (i) {
              final freshNext = widget.search.value.next.clone()..offset = i;
              widget.search.value = media.MediaSearchState(
                next: freshNext,
                count: widget.search.value.count,
              );
            },
            current: state.next.offset,
            empty: ds.Grid.int64(state.count) < state.next.limit,
            leading: [
              ds.CompactingMenu.pinned(
                DropdownUpload(
                  icon: SearchMimetypeDropdown.icon(mimex.checksum(state.next.mimetypes)),
                  help: ds.Hint(
                    const Text(
                      "filter by mimetype, upload files, torrents, magnet links, or switch to library search",
                    ),
                  ),
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
            ],
            help: ds.Hint(const Text("discover content over the network, use @ to access advanced filtering")),
          ),
        ),
        Expanded(
          child: DiscoveryGrid(
            search: widget.search,
            leading: [
              widget.downloading,
            ],
          ),
        ),
      ],
    );
  }
}
