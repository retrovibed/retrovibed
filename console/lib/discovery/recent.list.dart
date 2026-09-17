import 'package:flutter/material.dart';
import 'package:flutter/widget_previews.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:fixnum/fixnum.dart' as fixnum;

typedef FnRecentTap = Future<void> Function(BuildContext context, media.RecentRecordRequest item);

class RecentList extends StatefulWidget {
  const RecentList(
    this.mimetype, {
    super.key,
    this.latest = lib.recent.latest,
    this.tombstone = lib.recent.delete,
    this.onTap = defaultOnTap,
    this.padding,
    this.margin,
    this.host,
    this.authz,
  });

  final String mimetype;
  final lib.FnRecent latest;
  final lib.FnRecentTombstone tombstone;
  final FnRecentTap onTap;
  final EdgeInsets? padding;
  final EdgeInsets? margin;
  final String? host;
  final httpx.Option? authz;

  static Future<void> defaultOnTap(BuildContext context, media.RecentRecordRequest item) async {
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
  }

  @override
  State<RecentList> createState() => _RecentListState();
}

class _RecentListState extends State<RecentList> with ds.LoadingState {
  media.RecentSearchResponse _result = media.RecentSearchResponse();

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
    setState(() => loading = true);
    return httpx
        .withRetry(
          () => widget.latest(
            lib.recent.request(mimetype: widget.mimetype),
            host: widget.host,
            options: [widget.authz ?? authn.request(authn.AuthzCache.meta(context))],
          ),
        )
        .then(
          (resp) => setState(() {
            _result = resp;
            loading = false;
          }),
        )
        .catchError((cause) {
          setState(() {
            loading = false;
          });
        }, test: httpx.ErrorsTest.notimplemented)
        .catchError((cause) {
          setState(() {
            this.cause = ds.Errors.httpauto(cause, onTap: reseterr);
            loading = false;
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
            loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.Container(
      padding: widget.padding,
      margin: widget.margin,
      ds.Loading(
        loading: loading,
        cause: cause,
        Column(
          children: _result.items.map((item) {
            final deletion = () {
              return httpx.withRetry(
                () => widget
                    .tombstone(
                      item.id,
                      host: widget.host,
                      options: [widget.authz ?? authn.request(authn.AuthzCache.meta(context))],
                    )
                    .then((_) => _load(context)),
              );
            };
            return lib.KnownMediaRowDisplay.future(
              lib.known.autodetect(
                item.media,
                host: widget.host,
                options: [widget.authz ?? authn.request(authn.AuthzCache.meta(context))],
              ),
              onTap: () async => widget.onTap(context, item),
              trailing: [
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
