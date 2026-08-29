import 'package:flutter/material.dart';
import 'package:flutter/widget_previews.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:fixnum/fixnum.dart' as fixnum;
import 'recent.edit.dart';

class RecentList extends StatefulWidget {
  const RecentList(
    this.mimetype, {
    super.key,
    this.latest = lib.recent.latest,
    this.tombstone = lib.recent.delete,
  });

  final String mimetype;
  final lib.FnRecent latest;
  final lib.FnRecentTombstone tombstone;

  @override
  State<RecentList> createState() => _RecentListState();
}

class _RecentListState extends State<RecentList> {
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
    ds.postframe(() => _load(context));
  }

  @override
  void didUpdateWidget(RecentList oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.mimetype != widget.mimetype) {
      _load(context);
    }
  }

  Future<void> _load(BuildContext context) async {
    setState(() => _loading = true);
    return httpx
        .withRetry(
          () => widget.latest(
            lib.recent.request(mimetype: widget.mimetype),
            options: [authn.request(authn.AuthzCache.meta(context))],
          ),
        )
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
        }, test: httpx.ErrorsTest.notimplemented)
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
    return ds.Container(
      ds.Loading(
        loading: _loading,
        cause: _cause,
        Column(
          children: _result.items.map((item) {
            final deletion = () {
              return httpx.withRetry(
                () => widget
                    .tombstone(item.id, options: [authn.request(authn.AuthzCache.meta(context))])
                    .then((_) => _load(context)),
              );
            };
            return lib.KnownMediaRowDisplay.future(
              lib.known.autodetect(item.media, options: [authn.request(authn.AuthzCache.meta(context))]),
              onTap: () async {
                final pos = Duration(milliseconds: item.position.toInt());
                final dur = Duration(milliseconds: item.duration.toInt());
                final delta = (dur - pos).compareTo(Duration(seconds: 1));
                final playlist = media.Playlist.of(context);

                playlist?.setPlaylist(
                  item.query,
                  item.media,
                  playlist.autoqueue,
                  pos: delta < 0 ? Duration(milliseconds: 0) : Duration(milliseconds: item.position.toInt()),
                );
              },
              trailing: [
                // ds.LoadingIconButton.edit(
                //   onPressed: () async => ds.modals.push(
                //     context,
                //     RecentEdit(item, constraints: const BoxConstraints(maxWidth: 400)),
                //   ),
                // ),
                ds.LoadingIconButton.remove(onPressed: deletion),
              ],
            );
          }).toList(),
        ),
      ),
    );
  }
}

Future<media.RecentSearchResponse> recentListPreviewLatest(
  media.RecentSearchRequest req, {
  String? host,
  List<httpx.Option> options = const [],
}) {
  return Future.value(
    media.RecentSearchResponse(
      items: [
        media.RecentRecordRequest(
          id: '1',
          media: media.Media(description: 'Sample Movie Title'),
          duration: fixnum.Int64(5400000),
          position: fixnum.Int64(1800000),
        ),
        media.RecentRecordRequest(
          id: '2',
          media: media.Media(description: 'Another Show Episode'),
          duration: fixnum.Int64(2700000),
          position: fixnum.Int64(900000),
        ),
      ],
    ),
  );
}

Widget recentListPreviewWrapper(Widget child) {
  return MaterialApp(home: Material(child: child));
}

@Preview(name: 'Recent List', wrapper: recentListPreviewWrapper)
Widget recentListPreview() {
  return authn.AuthzCache(
    RecentList(mimex.video, latest: recentListPreviewLatest),
    current: authn.AuthzCache.fake,
  );
}
