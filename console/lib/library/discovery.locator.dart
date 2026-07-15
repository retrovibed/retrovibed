import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/discovery.dart' as disc;
import './api.dart' as api;

// DiscoveryLocator offers to search the wider p2p network for query/mimetype
// when neither the local library nor the catalog (library.Known) have a
// match - the free-text analog of KnownMediaLocator, which locates a
// specific already-catalogued item.
class DiscoveryLocator extends StatefulWidget {
  final String query;
  final String mimetype;
  final Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate;
  final Widget help;

  const DiscoveryLocator({
    super.key,
    required this.query,
    required this.mimetype,
    this.locate = api.locate.create,
    this.help = const ds.Hint(
      Text(
        "searches the peer-to-peer network and your search plugins for this title. "
        "when a match is found its added to recommendations.",
      ),
    ),
  });

  @override
  State<StatefulWidget> createState() => _DiscoveryLocator();
}

class _DiscoveryLocator extends State<DiscoveryLocator> {
  bool _queued = false;
  bool _loading = false;
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  void _onTap() async {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    final options = [authn.request(authn.AuthzCache.meta(context))];

    disc
        .ensureP2P(context, options: options)
        .then((proceed) {
          if (!proceed) {
            setState(() {
              _loading = false;
            });
            return null;
          }

          return widget
              .locate(
                api.Locate.create()
                  ..autodownload = false
                  ..query = widget.query
                  ..mimetype = widget.mimetype,
                options: options,
              )
              .then((v) {
                setState(() {
                  _queued = true;
                  _loading = false;
                });
              });
        })
        .catchError((e) {
          setState(() {
            _loading = false;
            _cause = ds.Error.unknown(e, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    if (widget.query.trim().isEmpty || widget.mimetype.isEmpty) {
      return ds.Empty;
    }

    return ds.Loading(
      loading: _loading,
      cause: _cause,
      Column(
        children: [
          ds.Card(
            margin: defaults.padding.copyWith(bottom: 0, top: 0),
            Center(
              child: Icon(
                _queued ? Icons.query_builder_rounded : Icons.travel_explore_rounded,
                size: 48,
              ),
            ),
            onTap: _loading || _queued ? null : _onTap,
            help: widget.help,
            trailing: [
              Center(
                child: Text(
                  'search the network for less known media',
                  textAlign: TextAlign.center,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
