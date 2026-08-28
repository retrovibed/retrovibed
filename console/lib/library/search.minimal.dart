import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'known.media.icon.dart';

class SearchMinimal extends StatefulWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final int capacity;
  final media.FnMediaSearch apisearch;
  final Widget empty;
  final Widget help;
  final Future<void> Function()? Function(BuildContext, media.Media, media.MediaSearchResponse) onPlay;

  const SearchMinimal({
    super.key,
    this.apisearch = media.media.search,
    this.empty = const Text("no results"),
    this.help = ds.Empty,
    this.onPlay = media.PlayAction,
    this.capacity = 32,
    required this.search,
  });

  @override
  State<StatefulWidget> createState() => _SearchMinimal();
}

class _SearchMinimal extends State<SearchMinimal> with ds.LoadingState {
  media.MediaSearchResponse _res = media.media.response(
    next: media.media.request(limit: 32),
  );

  void _searchChanged() {
    final freshNext = widget.search.value.next.clone()..offset = ds.Grid.int64(0);
    setState(() {
      _res.next = freshNext;
    });
    refresh(_res.next);
  }

  Future<void> refresh(media.MediaSearchRequest req) {
    return widget
        .apisearch(req, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            loading = false;
          });
        })
        .catchError((cause) {
          setState(() {
            cause = ds.Error.unauthorized(cause, onTap: reseterr);
            loading = false;
          });
        }, test: httpx.ErrorsTest.unauthorized)
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
    _res.next.limit = ds.Int64(widget.capacity);
    ds.postframe(() => refresh(_res.next));
    widget.search.addListener(_searchChanged);
  }

  @override
  void dispose() {
    widget.search.removeListener(_searchChanged);
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ds.Table<media.Media>(
      loading: loading,
      cause: cause,
      leading: ds.SearchTray(
        decoration: InputDecoration(hintText: "search ${httpx.host()}'s library"),
        onSubmitted: (v) {
          final freshNext = widget.search.value.next.clone()
            ..query = v
            ..offset = ds.Grid.int64(0);
          setState(() {
            _res.next = freshNext;
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
        empty: ds.Int64(_res.items.length) < _res.next.limit,
        leading: [],
      ),
      children: _res.items,
      ds.Table.expanded<media.Media>(
        (v) => media.RowDisplay(
          media: v,
          leading: [KnownMediaIcon(v)],
          onTap: widget.onPlay(context, v, _res),
        ),
      ),
      empty: Center(child: widget.empty),
    );
  }
}
