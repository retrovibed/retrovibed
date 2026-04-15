import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'api.dart';
import 'list.display.item.dart';
import 'community.create.dart';

typedef FnSearch =
    Future<CommunitySearchResponse> Function(
      CommunitySearchRequest req, {
      List<httpx.Option> options,
    });

class ListDisplay extends StatefulWidget {
  final FnSearch search;
  final FnSubscribe subscribe;

  const ListDisplay({
    super.key,
    this.search = API.search,
    this.subscribe = API.subscribe,
  });

  @override
  State<ListDisplay> createState() => _ListDisplayState();
}

class _ListDisplayState extends State<ListDisplay> {
  CommunitySearchResponse _resp = CommunitySearchResponse(
    next: CommunitySearchRequest(
      offset: ds.Int64(0),
      limit: ds.Int64(20),
    ),
  );
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  Widget _overlay = ds.Empty;

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _resetCause() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> _refresh(CommunitySearchRequest req) {
    setState(() => _loading = true);
    return httpx
        .withRetry(() => widget.search(req, options: [authn.AuthzCache.bearer(context)]))
        .then((response) {
          setState(() {
            _resp = response;
            _cause = ds.Error.zero;
          });
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: _resetCause);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _resetCause);
          });
        })
        .whenComplete(() {
          setState(() => _loading = false);
        });
  }

  void _resetoverlay() {
    setState(() {
      _overlay = ds.Empty;
    });
  }

  @override
  void initState() {
    super.initState();
    _refresh(_resp.next);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final communitycreation = ds.Masked(
      alignment: Alignment.center,
      SingleChildScrollView(
        child: CommunityCreate(
          constraints: BoxConstraints(maxWidth: 512.0),
          create: (c) {
            return httpx.withRetry(
              () => API.create(
                CommunityCreateRequest(community: c),
                options: [authn.DeeppoolAuthzCache.bearer(context)],
              ),
            );
          },
          onCreate: (_) {
            _refresh(_resp.next).whenComplete(_resetoverlay);
          },
          onCancel: _resetoverlay,
        ),
      ),
      reset: _resetoverlay,
    );
    return ds.Table(
      loading: _loading,
      cause: _cause,
      children: _resp.items,
      overlay: _overlay,
      empty: Center(child: Text('No communities found')),
      leading: ds.SearchTray(
        decoration: InputDecoration(hintText: "search communities"),
        autofocus: defaults.desktop,
        onSubmitted: (v) {
          setState(() {
            _resp.next
              ..query = v
              ..offset = ds.Int64(0);
          });
          return _refresh(_resp.next);
        },
        next: (i) {
          setState(() {
            _resp.next..offset = i;
          });
          _refresh(_resp.next);
        },
        current: _resp.next.offset,
        empty: ds.Int64(_resp.items.length) < _resp.next.limit,
        help: ds.Hint(const Text("find communities by name or description")),
        leading: [
          ds.CompactingMenu.pinned(
            ds.LoadingIconButton(
              icon: Icon(_overlay == ds.Empty ? Icons.add : Icons.remove),
              onPressed: () async {
                setState(() {
                  _overlay = _overlay == ds.Empty ? communitycreation : ds.Empty;
                });
              },
              help: ds.Hint(const Text("create a new community with a domain and description")),
            ),
          ),
        ],
      ),
      ds.Table.expanded<Community>(
        (v) => ListDisplayItem(
          community: v,
          onChanged: (c) => _refresh(_resp.next),
          subscribe: widget.subscribe,
        ),
      ),
    );
  }
}
