import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/discovery.dart' as disc;
import 'grid.display.dart';
import 'grid.setting.dart';
import 'menu.upload.files.dart';
import 'search.mimetype.dropdown.dart';
import 'dropdown.upload.dart';

class Home extends StatefulWidget {
  final media.FnMediaSearch apisearch;
  final media.FnUploadRequest apiupload;
  final TextEditingController? controller;
  final FocusNode? focus;
  final String highlighted;
  final ValueNotifier<media.MediaSearchState> search;

  const Home({
    super.key,
    this.apisearch = media.media.search,
    this.apiupload = media.media.upload,
    this.controller,
    this.focus,
    required this.highlighted,
    required this.search,
  });

  @override
  State<StatefulWidget> createState() => _HomeState();
}

enum _Mode { library, discovery }

class _HomeState extends State<Home> {
  _Mode _mode = _Mode.library;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  Widget _cause = ds.Error.zero;

  void _reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final compact = defaults.isCompact;

    return LayoutBuilder(
      builder: (context, constraints) => SingleChildScrollView(
        reverse: compact,
        physics: AlwaysScrollableScrollPhysics(),
        child: ConstrainedBox(
          constraints: BoxConstraints(minHeight: constraints.maxHeight),
          child: Column(
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
                  padding: defaults.padding.copyWith(bottom: 0.0),
                  tuning: GridSettings(),
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
                          ...SearchMimetypeDropdown.menuItems(state.next, (upd) {
                            widget.search.value = media.MediaSearchState(
                              next: upd.clone(),
                              count: widget.search.value.count,
                            );
                          }),
                          const PopupMenuDivider(),
                          PopupMenuItem<String>(
                            enabled: false,
                            child: mimex.CategoryOptionsLabel(state.next.mimetypes),
                          ),
                          MenuItemUploadFiles(
                            context,
                            widget.search,
                            apiupload: widget.apiupload,
                            mimetypes: state.next.mimetypes,
                          ),
                          downloads.MenuItemDownloadTorrent(context),
                          downloads.MenuItemDownloadMagnet(context),
                          PopupMenuItem<String>(
                            child: ListTile(
                              leading: Icon(_mode == _Mode.discovery ? Icons.check : Icons.travel_explore),
                              title: const Text("Search / Discover"),
                              onTap: () {
                                final next = _mode == _Mode.discovery ? _Mode.library : _Mode.discovery;
                                setState(() {
                                  _mode = next;
                                });
                                if (next == _Mode.discovery) {
                                  widget.focus?.requestFocus();
                                }
                              },
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                  help: ds.Hint(const Text("search your library, use @ to access advanced filtering")),
                ),
              ),
              ds.ErrorScreen(
                cause: _cause,
                switch (_mode) {
                  _Mode.library => Grid(
                    apisearch: widget.apisearch,
                    search: widget.search,
                    highlighted: widget.highlighted,
                  ),
                  _Mode.discovery => disc.DiscoveryGrid(search: widget.search),
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
