import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;

class SearchMinimal extends StatefulWidget {
  final media.FnMediaSearch apisearch;

  const SearchMinimal({
    super.key,
    this.apisearch = media.media.search,
  });

  @override
  State<StatefulWidget> createState() => _SearchMinimal();
}

class _SearchMinimal extends State<SearchMinimal> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  media.MediaSearchResponse _res = media.media.response(
    next: media.media.request(limit: 32),
  );

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh(media.MediaSearchRequest req) {
    return widget
        .apisearch(req, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });
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
    ds.postframe(() => refresh(_res.next));
  }

  @override
  Widget build(BuildContext context) {
    return ds.Table<media.Media>(
      loading: _loading,
      cause: _cause,
      leading: ds.SearchTray(
        decoration: InputDecoration(hintText: "search library"),
        onSubmitted: (v) {
          setState(() {
            _res.next.query = v;
            _res.next.offset = ds.Int64(0);
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
      ),
      children: _res.items,
      ds.Table.expanded<media.Media>(
        (v) => media.RowDisplay(
          media: v,
          leading: [Icon(mimex.icon(v.mimetype))],
          onTap: media.PlayAction(context, v, _res),
        ),
      ),
      empty: Center(child: Text("no results")),
    );
  }
}
