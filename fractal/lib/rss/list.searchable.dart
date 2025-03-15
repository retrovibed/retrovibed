import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:fixnum/fixnum.dart';
import 'package:protobuf/protobuf.dart';
import './list.dart';
import './feed.new.dart';
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
  api.Feed _created = api.Feed();
  api.FeedSearchResponse _res = api.FeedSearchResponse(
    next: api.FeedSearchRequest(query: '', offset: Int64(0), limit: Int64(10)),
    items: [],
  );

  Future<api.FeedSearchResponse> refresh() {
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

  void resetleading() => setState(() {
    print("UPDATING FEED CREATION");
    _overlay = null;
    _loading = false;
    _created = api.Feed();
  });

  void updatefeed(api.Feed upd) => setState(() {
    _created = upd;
    _overlay = _FeedCreate(
      current: upd,
      onCancel: resetleading,
      onSubmit: submitfeed,
      onChange: updatefeed,
    );
  });

  void submitfeed(api.Feed n) {
    setState(() => _loading = true);
    api
        .create(api.FeedCreateRequest(feed: n))
        .then((v) {
          refresh();
          return v;
        })
        .then((v) => resetleading())
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e);
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final feedproto = _FeedCreate(
      current: _created,
      onCancel: resetleading,
      onSubmit: submitfeed,
      onChange: updatefeed,
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
                _overlay = _overlay == null ? feedproto : null;
              });
            },
            icon: Icon(_overlay == null ? Icons.add : Icons.remove),
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
            onPressed:
                _res.next.offset.toInt() > 0
                    ? () {
                      setState(() {
                        _res.next.offset -= 1;
                      });
                      refresh();
                    }
                    : null,
            icon: Icon(Icons.arrow_left),
          ),
          IconButton(
            onPressed:
                _res.items.length == _res.next.limit
                    ? () {
                      setState(() {
                        _res.next.offset += 1;
                      });
                      refresh();
                    }
                    : null,
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
            _res = api.FeedSearchResponse(
              next: _res.next.deepCopy(),
              items: upd,
            );
          });
        },
      ),
      empty: feedproto,
      overlay: _overlay,
    );
  }
}

class _FeedCreate extends StatelessWidget {
  final api.Feed current;
  final Function(api.Feed)? onChange;
  final Function(api.Feed)? onSubmit;
  final Function()? onCancel;

  _FeedCreate({
    required this.current,
    this.onChange,
    this.onCancel,
    this.onSubmit,
  });

  @override
  Widget build(BuildContext context) {
    final themex = ds.Defaults.of(context);
    return Center(
      child: Column(
        spacing: themex.spacing ?? 0.0,
        children: [
          FeedNew(current: current, onChange: onChange),
          Row(
            spacing: themex.spacing ?? 0.0,
            children: [
              Spacer(),
              TextButton(onPressed: onCancel, child: Text("cancel")),
              TextButton(
                onPressed: () {
                  onSubmit?.call(current);
                },
                child: Text("create"),
              ),
              Spacer(),
            ],
          ),
        ],
      ),
    );
  }
}
