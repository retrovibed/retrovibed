import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/lucene.dart' as lucene;
import 'grid.display.dart';
import 'grid.setting.dart';
import 'search.mimetype.dropdown.dart';
import 'dropdown.upload.dart';
import 'empty.results.dart';

class LibrarySearch extends StatefulWidget {
  final media.FnMediaSearch apisearch;
  final media.FnUploadRequest apiupload;
  final TextEditingController? controller;
  final FocusNode? focus;
  final String highlighted;
  final ValueNotifier<media.MediaSearchState> search;
  final bool discovering;
  final VoidCallback onToggleMode;
  final Widget downloading;
  final void Function(Widget) onDownloadingChanged;

  const LibrarySearch({
    super.key,
    this.apisearch = media.media.search,
    this.apiupload = media.media.upload,
    this.controller,
    this.focus,
    required this.highlighted,
    required this.search,
    required this.discovering,
    required this.onToggleMode,
    required this.downloading,
    required this.onDownloadingChanged,
  });

  @override
  State<LibrarySearch> createState() => _LibrarySearchState();
}

class _LibrarySearchState extends State<LibrarySearch> {
  Widget _tuning = ds.Empty;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

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
            decoration: InputDecoration(hintText: "search library, @... for filters"),
            filters: [
              lucene.Boolean.auto('hidden', false, (v) {
                final freshNext = widget.search.value.next.clone()..hidden = v;
                widget.search.value = media.MediaSearchState(
                  next: freshNext,
                  count: widget.search.value.count,
                );
              }),
            ],
            controller: widget.controller,
            tuning: ds.buttons.settings(
              onPressed: () => setState(() {
                _tuning = _tuning == ds.Empty ? GridSettings() : ds.Empty;
              }),
              help: ds.Hint(Text("display advance settings")),
            ),
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
                      "filter by mimetype, upload files, torrents, magnet links, or switch to discover mode",
                    ),
                  ),
                  items: [
                    ...SearchMimetypeDropdown.menuItems(widget.search),
                    media.SearchModeToggle(discovering: widget.discovering, onToggle: widget.onToggleMode),
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
            help: ds.Hint(const Text("search your library, use @ to access advanced filtering")),
          ),
        ),
        Expanded(
          child: Grid(
            apisearch: widget.apisearch,
            search: widget.search,
            highlighted: widget.highlighted,
            empty: EmptyResults(onDiscover: widget.onToggleMode),
            leading: [
              _tuning,
              widget.downloading,
            ],
          ),
        ),
      ],
    );
  }
}
