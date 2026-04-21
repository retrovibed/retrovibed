import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/media.dart' as media;
import 'recent.edit.dart';

class Recent extends StatefulWidget {
  const Recent({
    super.key,
    this.latest = lib.recent.latest,
    this.tombstone = lib.recent.delete,
  });

  final lib.FnRecent latest;
  final lib.FnRecentTombstone tombstone;

  @override
  State<Recent> createState() => _RecentState();
}

class _RecentState extends State<Recent> {
  Widget _cause = ds.Error.zero;
  bool _loading = true;

  media.RecentSearchResponse _result = media.RecentSearchResponse();

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  void initState() {
    super.initState();
    _load(context);
  }

  Future<void> _load(BuildContext context) async {
    setState(() => _loading = true);
    final auth = authn.AuthzCache.bearer(context);
    return httpx
        .withRetry(() => widget.latest(lib.recent.request(), options: [auth]))
        .then(
          (resp) => setState(() {
            _result = resp;
            _loading = false;
          }),
        )
        .catchError((cause) {
          setState(() {
            _loading = false;
          });
        }, test: httpx.ErrorsTest.httpnotimplemented)
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: reseterr);
            _loading = false;
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: reseterr);
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.CarouselRow(
      title: const Text('Continue Watching'),
      constraints: BoxConstraints.tightForFinite(height: 256),
      background: ds.Repeat(() => lib.KnownMediaCard(lib.Known(), icon: null)),
        empty: Text(
          'Media you watch will appear here',
          style: TextStyle(color: Colors.grey),
        ),
      items:
          _result.items.map((item) {
            final deletion = () {
              return httpx.withRetry(
                () =>
                    widget.tombstone(item.id, options: [authn.AuthzCache.bearer(context)]).then((_) => _load(context)),
              );
            };
            return lib.KnownMediaCard.future(
              lib.known.autodetect(item.media, options: [authn.AuthzCache.bearer(context)]),
              onTap: () {
                final pos = Duration(milliseconds: item.position.toInt());
                final dur = Duration(milliseconds: item.duration.toInt());
                final delta = (dur - pos).compareTo(Duration(seconds: 1));
                final playlist = media.Playlist.of(context);

                playlist?.setPlaylist(
                  item.query,
                  media.range(
                    item.query,
                    item.media,
                    pos: delta < 0 ? Duration(milliseconds: 0) : Duration(milliseconds: item.position.toInt()),
                    options: () => [authn.AuthzCache.bearer(context)],
                  ),
                );
              },
              onLongPress: deletion,
              onSecondaryTap:
                  () => ds.modals.push(
                    context,
                    RecentEdit(item, constraints: const BoxConstraints(maxWidth: 400)),
                  ),
            );
          }).toList(),
      loading: _loading,
      cause: _cause,
    );
  }
}
