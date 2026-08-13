import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/library/search.dropdown.dart';
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
                SearchUploadDropdown(
                  search: widget.search,
                  mode: widget.mode,
                  onModeChanged: widget.onModeChanged,
                  apiupload: widget.apiupload,
                  onDownloadingChanged: widget.onDownloadingChanged,
                  help: ds.Hint(
                    const Text(
                      "filter by mimetype, upload files, torrents, magnet links, or switch to library search",
                    ),
                  ),
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
