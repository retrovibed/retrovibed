import 'package:flutter/material.dart';
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;
import './pending.typography.dart' as local;
import 'package:retrovibed/design.kit/stateful.dart';

typedef FnPendingSearch =
    Future<meta.ProfileSearchResponse> Function(
      meta.ProfileSearchRequest req, {
      List<httpx.Option> options,
    });

class PendingListDisplay extends StatefulWidget {
  final FnPendingSearch search;
  final TextEditingController? controller;
  final InputDecoration? decoration;
  const PendingListDisplay({
    super.key,
    this.search = meta.profiles.search,
    this.controller,
    this.decoration,
  });

  @override
  State<StatefulWidget> createState() => _PendingListDisplay();
}

class _PendingListDisplay extends State<PendingListDisplay> with LoadingState {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  meta.ProfileSearchResponse _res = meta.profiles.pending();

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
          });
        })
        .catchError((e) {
          setState(() {
            setState(() {
              _cause = ds.Error.unknown(e, onTap: resetcause);
            });
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
    final defaults = ds.Defaults.of(context);
    return _res.items.length > 0
        ? ds.Table(
          loading: _loading,
          cause: _cause,
          children: _res.items,
          leading: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              ds.SearchTray(
                autofocus: defaults.desktop,
                decoration: InputDecoration(
                  hintText: "search access requests",
                ),
                controller: widget.controller,
                onSubmitted: (v) {
                  setState(() {
                    _res.next.query = v;
                    _res.next.offset = ds.Table.offset(0);
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
                empty: ds.Table.offset(_res.items.length) < _res.next.limit,
              ),
              Row(
                mainAxisSize: MainAxisSize.min,
                spacing: defaults.spacing,
                children: [
                  SizedBox.fromSize(size: Size(90.0, 1.0)),
                  Expanded(child: Text("id")),
                  Expanded(child: Text("description")),
                  Expanded(child: Text("updated")),
                ],
              ),
            ],
          ),
          ds.Table.expanded<meta.Profile>(
            (v) => local.Typography(
              leading: [
                local.Typography.approvebtn(
                  context,
                  v.id,
                  onPressed: () {
                    meta.profiles
                        .update(
                          meta.ProfileUpdateRequest(
                            profile: v..disabledPendingApprovalAt = timex.inf.toIso8601String(),
                          ),
                          options: [authn.request(authn.AuthzCache.meta(context))],
                        )
                        .then((v) => refresh(_res.next..offset = ds.Table.offset(0)));
                  },
                ),
                local.Typography.removebtn(
                  context,
                  v.id,
                  onPressed: () {
                    meta.profiles
                        .disable(v.id, options: [authn.request(authn.AuthzCache.meta(context))])
                        .then((v) => refresh(_res.next..offset = ds.Table.offset(0)));
                  },
                ),
              ],
              v,
            ),
          ),
        )
        : const SizedBox();
  }
}
