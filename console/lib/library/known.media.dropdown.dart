import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/lucene.dart' as lucene;
import 'package:retrovibed/media.dart' as _media;
import 'package:retrovibed/meta/meta.search.pb.dart' as meta;
import 'package:retrovibed/timex.dart' as timex;
import 'known.media.card.dart';
import 'known.media.typography.dart';
import './api.dart' as api;

class KnownMediaDropdown extends StatefulWidget {
  final api.FnKnownSearch search;
  final TextEditingController? controller;
  final FocusNode? focus;
  final String current;
  final String mimetype;
  final Future<api.Known?> Function(api.Known? k) onChange;
  const KnownMediaDropdown({
    super.key,
    this.search = api.known.search,
    this.controller,
    this.focus,
    this.current = "",
    this.onChange = ds.fnAsyncPassthrough,
    this.mimetype = "",
  });

  // Applies [known] to [current] and fires the library metadatasync
  // endpoint, returning the server-updated [Media].  When [known] is null
  // and [current] has no known-media ID to clear (deactivation with nothing
  // ever selected), returns the unmodified [current].
  static Future<_media.Media> _sync(
    List<httpx.Option> authOptions,
    _media.Media current,
    api.Known? known, {
    api.FnLibraryMetadataSync libraryMetadataSync = _media.media.metadatasync,
  }) {
    if (known == null && uuidx.isMinMax(uuidx.fromString(current.knownMediaId))) {
      return Future.value(current);
    }
    final updated = current..knownMediaId = known?.uid ?? uuidx.min();
    return libraryMetadataSync(updated.id, updated, options: authOptions).then((v) => v.media);
  }

  static Future<void> Function() modal(
    BuildContext context,
    _media.Media current, {
    String mimetype = "",
    void Function(_media.Media) onChange = ds.fnNoop,
    api.FnKnownSearch search = api.known.search,
    api.FnLibraryMetadataSync libraryMetadataSync = _media.media.metadatasync,
  }) {
    return () {
      // Capture auth while the caller's context is still valid (modal opening).
      final authOptions = [authn.request(authn.AuthzCache.meta(context))];
      return ds.modals.asyncfn<void>(
        context,
        (completion) => ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 512.0),
          child: KnownMediaDropdown(
            current: current.knownMediaId,
            mimetype: mimetype,
            search: search,
            onChange: (known) {
              return _sync(
                    authOptions,
                    current,
                    known,
                    libraryMetadataSync: libraryMetadataSync,
                  )
                  .then<void>(onChange)
                  .then(
                    // dismissing the modal completes it before the deactivate triggered sync resolves.
                    (_) => completion.isCompleted ? null : completion.complete(),
                    onError: (e) => completion.isCompleted ? null : completion.completeError(e),
                  )
                  .then((_) => known);
            },
          ),
        ),
      );
    };
  }

  @override
  State<StatefulWidget> createState() => _KnownMediaDropdown();
}

class _KnownMediaDropdown extends State<KnownMediaDropdown> with ds.LoadingState {
  api.KnownSearchResponse _res = api.known.response(
    next: api.known.request(limit: 4),
  );
  api.Known? current = null;
  bool loaded = false;

  Future<void> refresh(api.KnownSearchRequest req) {
    return widget
        .search(req..mimetype = widget.mimetype, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            loading = false;
          });
          widget.focus?.requestFocus();
          ds.textediting.refocus(widget.controller);
        })
        .catchError((cause) {
          setState(() {
            loading = false;
          });
        }, test: httpx.ErrorsTest.err404)
        .catchError((cause) {
          setState(() {
            this.cause = ds.Errors.httpauto(cause, onTap: reseterr);
            loading = false;
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((e) {
          setState(() {
            cause = ds.Error.unknown(e, onTap: reseterr);
            loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();

    if (uuidx.isMinMax(uuidx.fromString(widget.current))) {
      ds.postframe(() {
        refresh(_res.next);
      });
      return;
    }

    ds.postframe(() {
      api.known
          .cached(
            widget.current,
            () => api.known.get(widget.current, options: [authn.request(authn.AuthzCache.meta(context))]),
          )
          .then(
            (w) => setState(() {
              loaded = true;
              current = w.known;
            }),
          )
          .catchError((cause) {
            setState(() {
              loading = false;
            });
          }, test: httpx.ErrorsTest.err404)
          .then((_) => refresh(_res.next))
          .catchError((cause) {
            setState(() {
              this.cause = ds.Errors.httpauto(cause, onTap: reseterr);
              loading = false;
            });
          }, test: httpx.ErrorsTest.httpauto)
          .catchError((e) {
            setState(() {
              cause = ds.Error.unknown(e, onTap: reseterr);
              loading = false;
            });
          });
    });
  }

  @override
  void deactivate() {
    if (loaded && widget.current != current?.uid) {
      widget.onChange(current);
    }
    super.deactivate();
  }

  @override
  Widget build(BuildContext context) {
    if (current != null) {
      return KnownMediaTypography(
        current!,
        trailing: [
          Spacer(),
          KnownMediaTypography.removebtn(
            context,
            widget.current,
            onPressed: () => setState(() {
              current = null;
            }),
          ),
        ],
      );
    }

    final defaults = ds.Defaults.of(context);

    return Column(
      spacing: defaults.spacing,
      children: [
        ds.Container(
          padding: defaults.padding,
          ds.SearchTray(
            decoration: InputDecoration(hintText: "search known media"),
            controller: widget.controller,
            focus: widget.focus,
            filters: [
              lucene.DateRange.auto('released', timex.Range.everything(), (r) {
                setState(() {
                  _res.next.released = meta.DateRange(
                    oldest: timex.formatISO8601(r.begin),
                    newest: timex.formatISO8601(r.end),
                  );
                  _res.next.offset = ds.Grid.int64(0);
                });
                refresh(_res.next);
              }),
              lucene.Boolean.auto('adult', false, (v) {
                setState(() {
                  _res.next.adult = v;
                  _res.next.offset = ds.Grid.int64(0);
                });
                refresh(_res.next);
              }),
            ],
            onSubmitted: (v) {
              setState(() {
                _res.next.query = v;
                _res.next.offset = ds.Grid.int64(0);
              });
              return refresh(_res.next);
            },
            next: (i) {
              setState(() {
                _res.next.offset = i;
              });
              refresh(_res.next);
            },
            current: _res.next.offset,
            empty: ds.Grid.int64(_res.items.length) < _res.next.limit,
            autofocus: defaults.desktop,
            tuning: ds.LoadingIconButton.close(
              onPressed: () async {
                widget.onChange(current);
              },
            ),
          ),
        ),
        ds.Grid(
          padding: EdgeInsets.zero,
          children: _res.items,
          loading: loading,
          cause: cause,
          leading: [],
          (context, v) {
            return KnownMediaCard(
              v,
              icon: Icons.search,
              onTap: () {
                setState(() {
                  current = v;
                });
                widget.onChange(v);
              },
            );
          },
        ),
      ],
    );
  }
}
