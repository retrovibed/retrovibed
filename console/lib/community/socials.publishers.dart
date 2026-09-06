import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'api.dart' as api;
import 'publisher.dropdown.dart';
import 'socials.publisher.row.dart';

// Manage the set of publishers attached to a community.
//
// The table is what this community publishes through - the selections
// themselves, each row resolving its own plugin - and the dropdown is the
// catalog minus those, so a plugin never appears in both at once.
class SocialsPublishers extends StatefulWidget {
  final api.Community community;
  final api.FnSocialsSearch search;
  final api.FnPublishersSearch publishers;
  final api.FnPublishersFind find;
  final api.FnSocialsEnable enable;
  final api.FnSocialsDisable disable;

  const SocialsPublishers(
    this.community, {
    super.key,
    this.search = api.socials.search,
    this.publishers = api.publishers.search,
    this.find = api.publishers.find,
    this.enable = api.socials.enable,
    this.disable = api.socials.disable,
  });

  @override
  State<StatefulWidget> createState() => _ListDisplay();
}

class _ListDisplay extends State<SocialsPublishers> with ds.LoadingState {
  final ValueNotifier<api.PluginPublisher> _dropdown = ValueNotifier<api.PluginPublisher>(api.PluginPublisher());
  List<api.CommunityPublisher> _attached = [];
  Widget _overlay = ds.Empty;

  @override
  void dispose() {
    _dropdown.dispose();
    super.dispose();
  }

  Future<void> _refresh() {
    setState(() => loading = true);
    return widget
        .search(
          api.SocialsSearchRequest()..communities.add(widget.community.id),
          options: [authn.request(authn.AuthzCache.meta(context))],
        )
        .then((v) {
          setState(() {
            // one community was asked for, so at most one social comes back;
            // none at all means it has attached nothing.
            _attached = v.items.isEmpty ? [] : v.items.first.publishers;
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

  // enabling is what makes the plugin disappear from the dropdown's results, so
  // the row is only added once the server has recorded the selection.
  void _attach(api.PluginPublisher selected) {
    setState(() => loading = true);
    widget
        .enable(
          widget.community.id,
          selected.id,
          options: [authn.request(authn.AuthzCache.meta(context))],
        )
        .then((v) {
          setState(() {
            _overlay = ds.Empty;
            _attached.add(v.enabled);
            cause = ds.Error.zero;
          });
        })
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
      children: _attached,
      overlay: _overlay,
      empty: const Center(child: Text('No publishers attached to this community')),
      leading: Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          IconButton(
            tooltip: 'attach a publish plugin to this community',
            onPressed: () {
              setState(() {
                _overlay = _overlay == ds.Empty
                    ? PublisherDropdown(
                        current: _dropdown,
                        readonly: true,
                        help: const ds.Hint(Text("search the plugins this community does not publish through yet")),
                        // the exclusion is the whole point: what is already
                        // attached is in the table below, not in the picker.
                        search: (req) => widget.publishers(
                          req..excluded.addAll(_attached.map((cp) => cp.publisherId)),
                          options: [authn.request(authn.AuthzCache.meta(context))],
                        ),
                        onSelected: _attach,
                      )
                    : ds.Empty;
              });
            },
            icon: Icon(_overlay == ds.Empty ? Icons.add : Icons.remove),
          ),
        ],
      ),
      ds.Table.expanded<api.CommunityPublisher>((v) {
        return SocialsPublisherRow(
          v,
          key: ValueKey(v.id),
          find: widget.find,
          disable: widget.disable,
          onDetach: (detached) {
            setState(() {
              _attached.removeWhere((cp) => cp.id == detached.id);
            });
          },
        );
      }),
    );
  }
}
