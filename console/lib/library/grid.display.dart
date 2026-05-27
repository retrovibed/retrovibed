import 'package:flutter/material.dart';
import 'package:language_code/language_code.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/discovery.dart' as disc;
import 'api.dart' as api;
import 'known.media.download.dart';
import 'known.media.display.dart';
import 'media.settings.dart';
import 'search.mimetype.dropdown.dart';
import 'grid.setting.dart';
import 'known.media.dropdown.dart';

class AvailableGridDisplay extends StatefulWidget {
  final media.FnMediaSearch apisearch;
  final media.FnUploadRequest apiupload;
  final TextEditingController? controller;
  final FocusNode? focus;
  final String highlighted;
  final ValueNotifier<media.MediaSearchResponse> search;

  const AvailableGridDisplay({
    super.key,
    this.apisearch = media.media.search,
    this.apiupload = media.media.upload,
    this.controller,
    this.focus,
    required this.highlighted,
    required this.search,
  });

  @override
  State<StatefulWidget> createState() => _AvailableGridDisplay();
}

class _AvailableGridDisplay extends State<AvailableGridDisplay> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh(media.MediaSearchRequest req, {bool refocus = false}) {
    return httpx
        .withRetry(
          () => widget.apisearch(req, options: [authn.request(authn.AuthzCache.meta(context))]),
        )
        .then((v) {
          setState(() {
            widget.search.value = v;
            _loading = false;
          });

          widget.focus?.requestFocus();
          if (refocus) ds.textediting.refocus(widget.controller);
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
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => refresh(widget.search.value.next));
  }

  @override
  Widget build(BuildContext context) {
    final upload = (
      ds.FilesEvent v, {
      ValueNotifier<int>? progress,
    }) {
      final multiparts = v.files.map((c) {
        return media.media.uploadable(c.path, c.name, c.mimeType!, progress: progress);
      });

      return Future.microtask(() {
        return Future.wait(
          multiparts.map((fv) {
            return fv.then((v) {
              return widget
                  .apiupload((req) {
                    req..files.add(v);
                    return req;
                  })
                  .then((uploaded) {
                    setState(() {
                      widget.search.value..items.add(uploaded.media);
                    });
                  })
                  .catchError((cause) {
                    setState(() {
                      _cause = ds.Error.unknown(cause, onTap: reseterr);
                    });
                  });
            });
          }),
        ).then((v) => ds.NullWidget).catchError((cause) {
          return ds.Error.unknown(cause, onTap: reseterr);
        });
      });
    };

    final replace = (media.Media v) {
      final replaced = widget.search.value.items.map((o) => o.id == v.id ? v : o);

      setState(() {
        widget.search.value = media.MediaSearchResponse(items: replaced, next: widget.search.value.next);
      });

      return v;
    };

    final category = mimex.category(widget.search.value.next.mimetypes);

    return RefreshIndicator(
      onRefresh: () => refresh(widget.search.value.next),
      child: LayoutBuilder(
        builder: (context, constraints) {
          final defaults = ds.Defaults.of(context);
          final compact = defaults.isCompact;
          return SingleChildScrollView(
            physics: AlwaysScrollableScrollPhysics(),
            child: ConstrainedBox(
              constraints: BoxConstraints(minHeight: constraints.maxHeight),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                mainAxisAlignment: MainAxisAlignment.start,
                verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
                children: [
                  ds.Grid<media.Media>(
                    children: widget.search.value.items,
                    loading: _loading,
                    cause: _cause,
                    leading: [
                      ds.SearchTray(
                        ensureVisible: true,
                        decoration: InputDecoration(hintText: "search library"),
                        filters: [
                          lucene.Boolean.auto('hidden', false, (v) {
                            setState(() => widget.search.value.next.hidden = v);
                            refresh(widget.search.value.next);
                          }),
                        ],
                        controller: widget.controller,
                        focus: widget.focus,
                        autofocus: defaults.desktop,
                        padding: defaults.padding.copyWith(bottom: 0.0),
                        onSubmitted: (v) {
                          setState(() {
                            widget.search.value.next
                              ..query = v
                              ..offset = ds.Grid.int64(0);
                            widget.search.value = widget.search.value;
                          });
                          return refresh(widget.search.value.next, refocus: true);
                        },
                        next: (i) {
                          setState(() {
                            widget.search.value.next.offset = i;
                          });
                          refresh(widget.search.value.next);
                        },
                        current: widget.search.value.next.offset,
                        empty: ds.Grid.int64(widget.search.value.items.length) < widget.search.value.next.limit,
                        leading: [
                          ds.CompactingMenu.pinned(
                            SearchMimetypeDropdown(
                              widget.search.value.next,
                              onChange: (upd) {
                                setState(() {
                                  widget.search.value.next = upd;
                                  widget.search.value = widget.search.value;
                                });
                                refresh(widget.search.value.next);
                              },
                            ),
                          ),
                          ds.FileDropWell.icon(
                            upload,
                            mimetypes: widget.search.value.next.mimetypes,
                            help: ds.Hint(const Text("drag and drop files onto the grid to add media to your library")),
                          ),
                        ],
                        tuning: defaults.debug ? GridSettings() : ds.SearchTray.zerobox,
                        help: ds.Hint(const Text("search your library, use @ to access advanced filtering")),
                      ),
                      (widget.search.value.next.query.isEmpty)
                          ? disc.Home(
                            category,
                            padding: defaults.padding.copyWith(
                              top: 0.0,
                              bottom: 0.0,
                            ),
                          )
                          : ds.Empty,
                    ],
                    empty:
                        widget.search.value.items.isEmpty
                            ? ds.Empty
                            : KnownMediaDownload.query(
                              () {
                                final search = widget.search.value;
                                if (_loading) return Future.value([]);
                                if (search.items.isNotEmpty) return Future.value([]);
                                if (search.next.mimetypes.isNotEmpty && !mimex.isVideo(category))
                                  return Future.value([]);

                                return httpx.withRetry(
                                  () => api.known
                                      .search(
                                        api.known.request(
                                          language: LanguageCode.code.locale.languageCode,
                                          adult: search.next.adult,
                                          query: search.next.query,
                                          limit: search.next.limit.toInt(),
                                        ),
                                        options: [authn.request(authn.AuthzCache.meta(context))],
                                      )
                                      .then((v) => v.items),
                                );
                              },
                              leading: Center(
                                child: Padding(
                                  padding: defaults.padding,
                                  child: Text(
                                    "no results in library. select below to automatically locate and download",
                                    style: TextStyle(
                                      color: Theme.of(
                                        context,
                                      ).colorScheme.onSurface.withValues(alpha: 0.6),
                                    ),
                                  ),
                                ),
                              ),
                              onTap: (v) {
                                return api.locate
                                    .create(
                                      api.Locate.create()..knownMediaId = v.id,
                                      options: [authn.request(authn.AuthzCache.meta(context))],
                                    )
                                    .then((_) {
                                      ScaffoldMessenger.of(context).showSnackBar(
                                        SnackBar(
                                          content: Text(
                                            'enqueued to be found and downloaded',
                                          ),
                                        ),
                                      );
                                      return v;
                                    })
                                    .catchError((e) {
                                      print("failed to initiate media location ${e}");
                                      ScaffoldMessenger.of(context).showSnackBar(
                                        SnackBar(
                                          content: Text(
                                            'failed to enqueue media locate and download',
                                          ),
                                        ),
                                      );
                                      throw e;
                                    });
                              },
                            ),
                    (context, _media) {
                      var onSettings = () {
                        ds.modals.asyncfn<media.Media>(
                          context,
                          (completion) => MediaSettings(
                            current: _media,
                            onChange: (pending, {bool forced = false, bool autoclose = false}) {
                              pending
                                  .then(replace)
                                  .then((v) {
                                    if (forced) refresh(widget.search.value.next, refocus: false);
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
                            onChange: replace,
                            mimetype: category,
                          ),
                        ),
                      ];

                      final key = ValueKey(uuidx.md5x("${_media.id}.${_media.updatedAt}"));

                      if (uuidx.isMinMax(
                        uuidx.fromString(_media.knownMediaId),
                      )) {
                        return KnownMediaDisplay.missing(
                          key: key,
                          _media,
                          onTap: media.PlayAction(context, _media, widget.search.value),
                          onSettings: onSettings,
                          onChange: replace,
                          highlighted: _media.id == widget.highlighted,
                          help: KnownMediaDisplay.hintPlayMedia,
                          trailing: trailing,
                        );
                      }

                      return KnownMediaDisplay(
                        key: key,
                        api.known
                            .cached(
                              _media.knownMediaId,
                              () => api.known.get(
                                _media.knownMediaId,
                                options: [authn.request(authn.AuthzCache.meta(context))],
                              ),
                            )
                            .then(
                              (w) => (w.known..description = _media.description),
                            ),
                        onTap: media.PlayAction(context, _media, widget.search.value),
                        onSettings: onSettings,
                        onChange: replace,
                        media: _media,
                        highlighted: _media.id == widget.highlighted,
                        help: KnownMediaDisplay.hintPlayMedia,
                        trailing: trailing,
                      );
                    },
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}
