import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:fractal/rss.dart';
import 'package:fixnum/fixnum.dart';
import 'package:fractal/rss/list.dart';
import 'package:protobuf/protobuf.dart';
import './api.dart' as api;

class ListSearchable extends StatefulWidget {
  final api.FnSearch search;

  ListSearchable({super.key, this.search = api.search});

  @override
  State<ListSearchable> createState() => SearchableView();
}

class SearchableView extends State<ListSearchable> {
  bool _loading = true;
  ds.Error? _cause = null;
  Widget? _overlay = null;
  Feed _created = Feed();
  FeedSearchResponse _res = FeedSearchResponse(
    next: FeedSearchRequest(query: '', offset: Int64(0), limit: Int64(10)),
    items: [],
  );

  Future<FeedSearchResponse> refresh() {
    return widget
        .search(_res.next)
        .then((r) {
          setState(() {
            _res = r;
          });
          return r;
        })
        .whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    refresh().catchError((e) {
      setState(() {
        _cause = ds.Error.unknown(e);
      });
      return _res;
    });
  }

  @override
  Widget build(BuildContext context) {
    final resetleading =
        () => setState(() {
          _overlay = null;
          _loading = false;
          _created = Feed();
        });
    final createfeed = (Feed n) {
      setState(() => _loading = true);
      api
          .create(FeedCreateRequest(feed: n))
          .then((v) => resetleading())
          .then((v) {
            refresh();
          })
          .catchError((e) {
            setState(() {
              _cause = ds.Error.unknown(e);
              _loading = false;
            });
          });
    };

    final feedproto = Center(
      child: Column(
        children: [
          Edit(
            feed: _created,
            onChange: (upd) {
              setState(() {
                _created = upd;
              });
            },
          ),
          Row(
            children: [
              Spacer(),
              TextButton(onPressed: resetleading, child: Text("cancel")),
              SizedBox(width: 10),
              TextButton(
                onPressed: () {
                  createfeed(_created);
                },
                child: Text("save"),
              ),
              Spacer(),
            ],
          ),
        ],
      ),
    );

    return ds.Table(
      loading: _loading,
      cause: _cause,
      leading: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          IconButton(
            onPressed: () {
              setState(() {
                _overlay = ds.Debug(feedproto);
              });
            },
            icon: Icon(Icons.add),
          ),
          Expanded(
            child: TextField(
              decoration: InputDecoration(hintText: "search feeds"),
              onChanged:
                  (v) => setState(() {
                    _res.next.query = v;
                  }),
              onSubmitted: (v) => refresh(),
            ),
          ),
          IconButton(
            onPressed: () {
              setState(() {
                _res.next.offset -= 1;
              });
              refresh();
            },
            icon: Icon(Icons.arrow_left),
          ),
          IconButton(
            onPressed: () {
              setState(() {
                _res.next.offset += 1;
              });
              refresh();
            },
            icon: Icon(Icons.arrow_right),
          ),
        ],
      ),
      children: _res.items,
      (w) => Item(
        current: w,
        onChange: (v) {
          final upd =
              _res.items.map((old) => old.id == v.id ? v : old).toList();
          setState(() {
            _res = FeedSearchResponse(next: _res.next.deepCopy(), items: upd);
          });
        },
      ),
      empty: feedproto,
      overlay: _overlay,
    );
  }
}
