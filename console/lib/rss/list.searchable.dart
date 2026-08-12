import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'list.dart';
import 'feed.new.dart';
import 'api.dart' as api;
import 'package:retrovibed/design.kit/stateful.dart';

class ListSearchable extends StatefulWidget {
  final api.FnSearch search;

  ListSearchable({super.key, this.search = api.search});

  @override
  State<ListSearchable> createState() => SearchableView();
}

class SearchableView extends State<ListSearchable> with LoadingState {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  Widget _overlay = ds.Empty;
  api.Feed _created = api.Feed();
  api.FeedSearchResponse _res = api.FeedSearchResponse(
    next: api.FeedSearchRequest(
      query: '',
      offset: ds.Int64(0),
      limit: ds.Int64(10),
    ),
    items: [],
  );

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<api.FeedSearchResponse> refresh(api.FeedSearchRequest next) {
    return widget
        .search(next, options: [authn.request(authn.AuthzCache.meta(context))])
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
    ds.postframe(() {
      refresh(_res.next).catchError((e) {
        setState(() {
          _cause = ds.Error.unknown(e, onTap: reseterr);
        });
        return _res;
      });
    });
  }

  void resetleading() => setState(() {
    _overlay = ds.Empty;
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
        .create(
          api.FeedCreateRequest(feed: n),
          options: [authn.request(authn.AuthzCache.meta(context))],
        )
        .then((v) {
          refresh(_res.next);
          return v;
        })
        .then((v) => resetleading())
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e, onTap: reseterr);
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final feedproto = _FeedCreate(
      current: _created,
      onCancel: resetleading,
      onSubmit: submitfeed,
      onChange: updatefeed,
    );

    return ds.Table(
      loading: _loading,
      cause: _cause,
      leading: ds.SearchTray(
        autofocus: defaults.desktop,
        decoration: InputDecoration(hintText: "search feeds"),
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
        empty: _res.items.length < _res.next.limit.toInt(),
        leading: [
          IconButton(
            onPressed: () {
              setState(() {
                _overlay = _overlay == ds.Empty ? feedproto : ds.Empty;
              });
            },
            icon: Icon(_overlay == ds.Empty ? Icons.add : Icons.remove),
          ),
        ],
      ),
      children: _res.items,
      ds.Table.expanded<api.Feed>(
        (w) => Item(
          key: ValueKey(w.id),
          current: w,
          onChange: (v) {
            final upd = ds.fnOnChange(_res.items, v, (old) => old.id == w.id);
            setState(() {
              _res = api.FeedSearchResponse(
                next: _res.next.deepCopy(),
                items: upd,
              );
            });
          },
        ),
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
    final defaults = ds.Defaults.of(context);
    return ds.Card(
      alignment: null,
      padding: EdgeInsets.zero,
      forms.Container(
        Column(
          spacing: defaults.spacing,
          mainAxisSize: MainAxisSize.min,
          children: [
            FeedNew(current: current, onChange: onChange),
            Row(
              spacing: defaults.spacing,
              children: [
                Spacer(),
                ds.LoadingButton(
                  Text("cancel"),
                  onPressed: () async => onCancel?.call(),
                ),
                ds.LoadingButton(
                  Text("create"),
                  onPressed: () async => onSubmit?.call(current),
                ),
                Spacer(),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
