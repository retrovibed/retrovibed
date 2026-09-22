import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'api.dart';
import 'list.display.item.dart';
import 'community.create.dart';
import 'qr.scanner.dart';
import 'link.content.dart';

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
    this.search = communities.search,
    this.subscribe = communities.subscribe,
  });

  @override
  State<ListDisplay> createState() => _ListDisplayState();
}

class _ListDisplayState extends State<ListDisplay> with ds.LoadingState {
  CommunitySearchResponse _resp = CommunitySearchResponse(
    next: CommunitySearchRequest(
      offset: ds.Int64(0),
      limit: ds.Int64(20),
    ),
  );
  Widget _overlay = ds.Empty;

  Future<void> _refresh(CommunitySearchRequest req) {
    setState(() => loading = true);
    return httpx
        .withRetry(() => widget.search(req, options: [authn.request(authn.AuthzCache.meta(context))]))
        .then((response) {
          setState(() {
            _resp = response;
            cause = ds.Error.zero;
          });
        })
        .catchError((cause) {
          setState(() {
            this.cause = ds.Errors.httpauto(cause, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        })
        .whenComplete(() {
          setState(() => loading = false);
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
    ds.postframe(() => _refresh(_resp.next));
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
              () => communities.create(
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
      loading: loading,
      cause: cause,
      children: _resp.items,
      overlay: _overlay,
      empty: Center(child: Text('No communities found')),
      leading: ds.SearchTray(
        autoscroll: true,
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
              help: ds.Hint(const Text("create a new community with a url and description")),
            ),
          ),
        ],
        trailing: [
          if (!defaults.desktop)
            ds.CompactingMenu.pinned(
              ds.LoadingIconButton(
                icon: Icon(Icons.qr_code_scanner),
                onPressed: () {
                  setState(() {
                    _overlay = QRScannerModal(
                      onScanned: (community, attribution) =>
                          handleSubscribeAction(context, community, attribution).then(
                            (_) {
                              _resetoverlay();
                              return _refresh(_resp.next);
                            },
                          ),
                      onCancel: _resetoverlay,
                    );
                  });
                  return Future.value();
                },
                help: ds.Hint(const Text("scan a QR code to subscribe to content")),
              ),
            ),
        ],
      ),
      ds.Table.expanded<Community>(
        (v) => ListDisplayItem(
          community: v,
          onChanged: (c) {
            final upd = ds.fnOnChange(_resp.items, c, (o) => o.id == v.id);
            setState(() {
              _resp = CommunitySearchResponse(
                items: upd,
                next: _resp.next,
              );
            });
          },
          subscribe: widget.subscribe,
        ),
      ),
    );
  }
}
