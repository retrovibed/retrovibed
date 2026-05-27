import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as _media;
import 'known.media.card.dart';
import 'known.media.typography.dart';
import './api.dart' as api;

class KnownMediaDropdown extends StatefulWidget {
  final api.FnKnownSearch search;
  final TextEditingController? controller;
  final FocusNode? focus;
  final String current;
  final String mimetype;
  final void Function(api.Known? k)? onChange;
  const KnownMediaDropdown({
    super.key,
    this.search = api.known.search,
    this.controller,
    this.focus,
    this.current = "",
    this.onChange,
    this.mimetype = "",
  });

  // Applies [known] to [current] and fires the appropriate metadatasync
  // endpoint, returning the server-updated [Media].  When [known] is null
  // and [current] has no known-media ID to clear (deactivation with nothing
  // ever selected), returns the unmodified [current].
  // [authOptions] must be pre-captured by the caller while the context is
  // still valid — _sync itself has no BuildContext dependency.
  // Both sync functions default to the real API and can be replaced in tests.
  static Future<_media.Media> _sync(
    List<httpx.Option> authOptions,
    _media.Media current,
    api.Known? known, {
    api.FnLibraryMetadataSync libraryMetadataSync = _media.media.metadatasync,
    api.FnDiscoveredMetadataSync discoveredMetadataSync = _media.discovered.metadatasync,
  }) {
    if (known == null && uuidx.isMin(uuidx.fromString(current.knownMediaId))) {
      return Future.value(current);
    }
    final updated = current..knownMediaId = known?.id ?? uuidx.min();
    if (uuidx.isMin(uuidx.fromString(current.torrentId))) {
      return libraryMetadataSync(updated.id, updated, options: authOptions).then((v) => v.media);
    }
    return discoveredMetadataSync(updated.torrentId, updated, options: authOptions).then((v) => v.media);
  }

  static Future<void> Function() modal(
    BuildContext context,
    _media.Media current, {
    String mimetype = "",
    void Function(_media.Media)? onChange,
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
            onChange: (known) {
              _sync(authOptions, current, known)
                  .then<void>((v) {
                    onChange?.call(v);
                    completion.complete();
                  })
                  .catchError(completion.completeError);
            },
          ),
        ),
      );
    };
  }

  /// Returns a [KnownMediaDropdown] widget configured to synchronise metadata
  /// when the user selects a known-media entry.  Routing mirrors [modal]:
  /// - no torrent → calls [apiLibraryMetadataSync] (`/m/:id/metadatasync`)
  /// - with torrent → calls [apiDiscoveredMetadataSync] (`/d/:id/metadatasync`)
  ///
  /// Both sync functions are injectable so they can be replaced in tests.
  static Widget inline(
    BuildContext context,
    _media.Media current, {
    String mimetype = "",
    void Function(_media.Media)? onChange,
    api.FnKnownSearch search = api.known.search,
    api.FnLibraryMetadataSync apiLibraryMetadataSync = _media.media.metadatasync,
    api.FnDiscoveredMetadataSync apiDiscoveredMetadataSync = _media.discovered.metadatasync,
  }) {
    // Capture auth while the caller's context is still valid (widget build time),
    // mirroring the modal approach.  The onChange closure may fire during
    // deactivate() when dependOnInheritedWidgetOfExactType is no longer safe.
    final authOptions = [authn.request(authn.AuthzCache.meta(context))];
    return KnownMediaDropdown(
      current: current.knownMediaId,
      mimetype: mimetype,
      search: search,
      onChange: (known) {
        _sync(
          authOptions,
          current,
          known,
          libraryMetadataSync: apiLibraryMetadataSync,
          discoveredMetadataSync: apiDiscoveredMetadataSync,
        ).then<void>((v) => onChange?.call(v));
      },
    );
  }

  @override
  State<StatefulWidget> createState() => _KnownMediaDropdown();
}

class _KnownMediaDropdown extends State<KnownMediaDropdown> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  api.KnownSearchResponse _res = api.known.response(
    next: api.known.request(limit: 4),
  );
  api.Known? current = null;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh(api.KnownSearchRequest req) {
    return widget
        .search(req..mimetype = widget.mimetype, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });
          widget.focus?.requestFocus();
          ds.textediting.refocus(widget.controller);
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unauthorized(cause, onTap: reseterr);
            _loading = false;
          });
        }, test: httpx.ErrorsTest.unauthorized)
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

    if (uuidx.isMinMax(uuidx.fromString(widget.current))) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        refresh(_res.next);
      });
      return;
    }

    WidgetsBinding.instance.addPostFrameCallback((_) {
      api.known
          .cached(
            widget.current,
            () => api.known.get(widget.current, options: [authn.request(authn.AuthzCache.meta(context))]),
          )
          .then(
            (w) => setState(() {
              current = w.known;
            }),
          )
          .whenComplete(() => refresh(_res.next));
    });
  }

  @override
  void deactivate() {
    if (current == null) {
      widget.onChange?.call(current);
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
            onPressed:
                () => setState(() {
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
          ),
        ),
        ds.Loading(
          loading: _res.items.isEmpty,
          overlay: ds.Empty,
          forms.Container(
            ds.Grid(
              children: _res.items,
              loading: _loading,
              cause: _cause,
              leading: [],
              (context, v) {
                return KnownMediaCard(
                  v,
                  icon: Icons.search,
                  onDoubleTap:
                      widget.onChange == null
                          ? null
                          : () {
                            setState(() {
                              current = v;
                            });
                            widget.onChange!(v);
                          },
                );
              },
            ),
          ),
        ),
      ],
    );
  }
}
