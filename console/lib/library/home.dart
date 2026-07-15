import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/discovery.dart' as disc;
import 'grid.display.dart';
import 'search.mimetype.dropdown.dart';

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
                    lucene.Mode.auto(
                      'discover',
                      false,
                      (enabled) {
                        setState(() {
                          _mode = enabled ? _Mode.discovery : _Mode.library;
                        });
                      },
                      help: ds.Hint(
                        const Text(
                          "switch to discover mode to find and download new media, instead of searching your existing library",
                        ),
                      ),
                    ),
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
                      SearchMimetypeDropdown(
                        state.next,
                        onChange: (upd) {
                          widget.search.value = media.MediaSearchState(
                            next: upd.clone(),
                            count: widget.search.value.count,
                          );
                        },
                      ),
                    ),
                  ],
                  help: ds.Hint(const Text("search your library, use @ to access advanced filtering")),
                ),
              ),
              switch (_mode) {
                _Mode.library => Grid(
                  apisearch: widget.apisearch,
                  search: widget.search,
                  highlighted: widget.highlighted,
                ),
                _Mode.discovery => disc.DiscoveryGrid(search: widget.search),
              },
            ],
          ),
        ),
      ),
    );
  }
}
