import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'api.dart' as api;
import 'publisher.dropdown.dart';
import 'publisher.list.row.dart';

// Manage the set of publishers attached to a community.
class SocialsPublishers extends StatefulWidget {
  final api.Community community;
  SocialsPublishers(this.community, {super.key});

  @override
  State<StatefulWidget> createState() => _ListDisplay();
}

class _ListDisplay extends State<SocialsPublishers> with ds.LoadingState {
  ValueNotifier<api.PluginPublisher> _dropdown = ValueNotifier<api.PluginPublisher>(api.PluginPublisher());
  api.PluginPublisherSearchResponse _res = api.PluginPublisherSearchResponse();
  Widget _overlay = ds.Empty;

  Future<void> _refresh() {
    setState(() => loading = true);
    return api.publishers
        .search(_res.next, options: [authn.request(authn.AuthzCache.meta(context))])
        .then((v) {
          setState(() {
            _res = v;
            cause = ds.Error.zero;
          });
        })
        .catchError((cause) {}, test: httpx.ErrorsTest.err404)
        .catchError((cause) {
          setState(() => this.cause = ds.Error.unauthorized(cause, onTap: reseterr));
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((cause) {
          setState(() => this.cause = ds.Error.unknown(cause, onTap: reseterr));
        })
        .whenComplete(() => setState(() => loading = false));
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(_refresh);
  }

  @override
  Widget build(BuildContext context) {
    return ds.Table(
      loading: loading,
      cause: cause,
      children: _res.items,
      overlay: _overlay,
      leading: ds.SearchTray(
        leading: [
          IconButton(
            onPressed: () {
              setState(() {
                _overlay = _overlay == ds.Empty ? PublisherDropdown(current: _dropdown, readonly: true) : ds.Empty;
              });
            },
            icon: Icon(_overlay == ds.Empty ? Icons.add : Icons.remove),
          ),
        ],
        onSubmitted: (v) {
          setState(() {
            _res.next
              ..query = v
              ..offset = ds.Int64(0);
          });
          return _refresh();
        },
        next: (i) {
          setState(() {
            _res.next.offset = i;
          });
          _refresh();
        },
        current: _res.next.offset,
        empty: ds.Int64(_res.items.length) < _res.next.limit,
      ),
      ds.Table.expanded<api.PluginPublisher>((v) {
        return PublisherRow(
          v,
          key: ValueKey(v.id),
          onDelete: (deleted) {
            setState(() {
              _res.items.removeWhere((p) => p.id == deleted.id);
            });
          },
        );
      }),
    );
  }
}
