import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;
import './create.inlined.dart';
import './filter.status.dart';
import './list.row.dart';

typedef FnSearch =
    Future<meta.ProfileSearchResponse> Function(
      meta.ProfileSearchRequest req, {
      List<httpx.Option> options,
    });

class ListDisplay extends StatefulWidget {
  final FnSearch search;
  final TextEditingController? controller;
  final InputDecoration? decoration;
  final ValueNotifier<int>? events;
  const ListDisplay({
    super.key,
    this.search = meta.profiles.search,
    this.controller,
    this.decoration,
    this.events,
  });

  @override
  State<StatefulWidget> createState() => _ListDisplay();
}

class _ListDisplay extends State<ListDisplay> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  Widget _create = ds.Empty;
  meta.ProfileSearchResponse _res = meta.profiles.response(
    next: meta.profiles.request(limit: 32),
  );
  VoidCallback? _eventsListener;

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void resetcause() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh(meta.ProfileSearchRequest req) {
    return widget
        .search(req, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
            _cause = ds.Error.zero;
          });
        })
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e, onTap: resetcause);
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(() => refresh(_res.next));
    _eventsListener = () {
      if (mounted) refresh(_res.next);
    };
    widget.events?.addListener(_eventsListener!);
  }

  @override
  void dispose() {
    if (_eventsListener != null) {
      widget.events?.removeListener(_eventsListener!);
    }
    super.dispose();
  }

  void _replace(meta.Profile v) {
    setState(() {
      _res = meta.ProfileSearchResponse(
        items: ds.fnOnChange(_res.items, v, (o) => o.id == v.id),
        next: _res.next,
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    final togglecreate = () {
      final reset = () => setState(() {
        _create = ds.Empty;
      });
      setState(() {
        _create = _create == ds.Empty ? CreateInlined(onClose: () => Future.microtask(reset)) : ds.Empty;
      });
    };

    return ds.Table(
      padding: defaults.padding.copyWith(top: 0, bottom: 0),
      loading: _loading,
      cause: _cause,
      children: _res.items,
      overlay: _create,
      leading: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ds.SearchTray(
            autofocus: defaults.desktop,
            decoration: widget.decoration ?? InputDecoration(hintText: "search users"),
            controller: widget.controller,
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
            leading: [
              IconButton(
                onPressed: togglecreate,
                icon: Icon(_create == ds.Empty ? Icons.person_add : Icons.close),
              ),
            ],
            trailing: [
              FilterStatus(meta.ProfileStatus.valueOf(_res.next.status)!, (
                upd,
              ) {
                setState(() {
                  _res.next.status = upd?.value ?? _res.next.status;
                });
                refresh(_res.next);
              }),
            ],
          ),
          Builder(
            builder: (context) {
              final compact = defaults.isCompact;
              return ds.TableHeader(
                [
                  SizedBox.square(dimension: 12),
                  if (!compact) Expanded(child: Text("id")),
                  Expanded(
                    child: Text(
                      "username",
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  if (!compact) Expanded(child: Text("updated")),
                ],
                padding: defaults.padding / 2,
              );
            },
          ),
        ],
      ),
      ds.Table.expanded<meta.Profile>((v) => ListRow(key: ValueKey(v.id), v, onChange: _replace)),
    );
  }
}
