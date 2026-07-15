import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/discovery.dart' as disc;
import 'known.media.display.dart';
import 'media.settings.dart';
import 'known.media.dropdown.dart';

class Grid extends StatefulWidget {
  final media.FnMediaSearch apisearch;
  final ValueNotifier<media.MediaSearchState> search;
  final String highlighted;

  const Grid({
    super.key,
    this.apisearch = media.media.search,
    required this.search,
    required this.highlighted,
  });

  @override
  State<Grid> createState() => _GridState();
}

class _GridState extends State<Grid> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  List<media.Media> _items = [];
  media.MediaSearchRequest? _lastFetchedNext;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  media.Media _replace(media.Media v) {
    setState(() {
      _items = _items.map((o) => o.id == v.id ? v : o).toList();
    });
    return v;
  }

  Future<void> refresh() {
    final req = widget.search.value.next;
    return httpx
        .withRetry(
          () => widget.apisearch(req, options: [authn.request(authn.AuthzCache.meta(context))]),
        )
        .then((v) {
          setState(() {
            _items = v.items;
            _loading = false;
          });
          widget.search.value = media.MediaSearchState(next: req, count: v.items.length);
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: reseterr);
            _loading = false;
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e, onTap: reseterr);
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<media.MediaSearchState>(
      valueListenable: widget.search,
      builder: (context, state, _) {
        if (!identical(state.next, _lastFetchedNext)) {
          _lastFetchedNext = state.next;
          ds.postframe(() => refresh());
        }

        final defaults = ds.Defaults.of(context);
        final category = mimex.category(state.next.mimetypes);

        return ds.Grid<media.Media>(
          children: _items,
          loading: _loading,
          cause: _cause,
          leading: [
            Visibility(
              key: ValueKey('library.query.disc.home'),
              visible: state.next.query.isEmpty,
              replacement: ds.Empty,
              child: disc.Home(
                category,
                key: ValueKey('library.disc.home'),
                padding: defaults.padding.copyWith(
                  top: 0.0,
                  bottom: 0.0,
                ),
              ),
            ),
          ],
          empty: Center(
            child: Padding(
              padding: defaults.padding,
              child: Text(
                "no results in library",
                style: TextStyle(
                  color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
                ),
              ),
            ),
          ),
          (context, _media) {
            final onSettings = () {
              ds.modals.asyncfn<media.Media>(
                context,
                (completion) => MediaSettings(
                  current: _media,
                  onChange: (pending, {bool forced = false, bool autoclose = false}) {
                    pending
                        .then(_replace)
                        .then((v) {
                          if (forced) refresh();
                          if (autoclose) completion.complete(v);
                        })
                        .catchError((cause) {
                          setState(() {
                            _cause = ds.Error.unknown(cause, onTap: reseterr);
                          });
                        });
                  },
                ),
              );
            };
            final trailing = [
              ds.LoadingIconButton.info(
                tooltip: "manually identify the media",
                help: ds.Hint(
                  Text("search for and select the correct media identity from the known library"),
                ),
                onPressed: KnownMediaDropdown.modal(
                  context,
                  _media,
                  onChange: _replace,
                  mimetype: category,
                ),
              ),
            ];

            return KnownMediaDisplay.auto(
              context,
              _media,
              onTap: media.PlayAction(
                context,
                _media,
                media.MediaSearchResponse(next: state.next, items: _items),
              ),
              onSettings: onSettings,
              onChange: _replace,
              highlighted: _media.id == widget.highlighted,
              help: KnownMediaDisplay.hintPlayMedia,
              trailing: trailing,
            );
          },
        );
      },
    );
  }
}
