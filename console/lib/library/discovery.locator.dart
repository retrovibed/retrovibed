import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/discovery.dart' as disc;
import 'package:retrovibed/uuidx.dart' as uuidx;
import './api.dart' as api;

// DiscoveryLocator offers to search the wider p2p network for query/mimetype
// when neither the local library nor the catalog (library.Known) have a
// match - the free-text analog of KnownMediaLocator, which locates a
// specific already-catalogued item.
class DiscoveryLocator extends StatefulWidget {
  final String query;
  final String mimetype;
  final Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate;
  final Future<api.LocateLookupResponse> Function(String id, {List<httpx.Option> options}) lookup;
  final Future<bool> Function(BuildContext context, {List<httpx.Option> options}) ensureP2P;
  final Future<Widget> Function(api.Locate located) onFound;
  final Widget help;

  const DiscoveryLocator({
    super.key,
    required this.query,
    required this.mimetype,
    required this.onFound,
    this.locate = api.locate.create,
    this.lookup = api.locate.get,
    this.ensureP2P = disc.ensureP2P,
    this.help = const ds.Hint(
      Text(
        "searches the peer-to-peer network and your search plugins for this title. "
        "once found, you can download it immediately.",
      ),
    ),
  });

  @override
  State<StatefulWidget> createState() => _DiscoveryLocator();
}

enum _LocateState { idle, loading, pending, found }

class _DiscoveryLocator extends State<DiscoveryLocator> {
  _LocateState _state = _LocateState.idle;
  String _locateId = '';
  Duration _interval = Duration.zero;
  Widget _cause = ds.Error.zero;
  Future<Widget>? _foundContent;

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
      _state = _LocateState.loading;
      _cause = ds.Error.zero;
    });

    final options = [authn.request(authn.AuthzCache.meta(context))];

    widget
        .ensureP2P(context, options: options)
        .then((proceed) {
          if (!proceed) {
            setState(() {
              _state = _LocateState.idle;
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
                _locateId = v.locate.id;
                setState(() {
                  _state = _LocateState.pending;
                  _interval = const Duration(seconds: 10);
                });
              });
        })
        .catchError((e) {
          setState(() {
            _state = _LocateState.idle;
            _cause = ds.Error.unknown(e, onTap: reseterr);
          });
        });
  }

  Future<bool> _checkLocated() {
    final options = [authn.request(authn.AuthzCache.meta(context))];

    return widget
        .lookup(_locateId, options: options)
        .then<bool>((v) {
          print("RECEIVED ${v}");
          if (uuidx.isMax(uuidx.fromString(v.locate.locatedTorrentId))) {
            return false;
          }

          print("FOUND ${v}");
          setState(() {
            _state = _LocateState.found;
            _interval = Duration.zero;
            _foundContent = widget.onFound(v.locate);
          });
          return true;
        })
        .catchError((e) {
          print("WAKA");
          debugPrint('$e');
          return false;
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    if (widget.query.trim().isEmpty || widget.mimetype.isEmpty) {
      return ds.Empty;
    }

    final icon = switch (_state) {
      _LocateState.idle => Icons.travel_explore_rounded,
      _LocateState.loading => Icons.travel_explore_rounded,
      _LocateState.pending => Icons.query_builder_rounded,
      _LocateState.found => Icons.check_circle_rounded,
    };

    final iconColor = switch (_state) {
      _LocateState.idle => null,
      _LocateState.loading => null,
      _LocateState.pending => null,
      _LocateState.found => defaults.success,
    };

    final idle = Icon(icon, size: 48, color: iconColor);
    final content = switch (_state) {
      _LocateState.found =>
        _foundContent == null
            ? idle
            : FutureBuilder<Widget>(
                future: _foundContent,
                builder: (context, snapshot) {
                  return ds.Loading(
                    snapshot.data ?? idle,
                    loading: !(snapshot.hasData || snapshot.hasError),
                    cause: snapshot.hasError ? ds.Errors.httpauto(snapshot.error!, onTap: _onTap) : ds.Error.zero,
                  );
                },
              ),
      _ => idle,
    };

    final trailingText = switch (_state) {
      _LocateState.idle => 'search networks for media',
      _LocateState.loading => 'search networks for media',
      _LocateState.pending => 'searching networks for media...',
      _LocateState.found => 'found',
    };

    final onTap = switch (_state) {
      _LocateState.idle => _onTap,
      _LocateState.loading => null,
      _LocateState.pending => null,
      _LocateState.found => null,
    };

    return ds.Loading(
      loading: _state == _LocateState.loading,
      cause: _cause,
      ds.Poll(
        Column(
          children: [
            ds.Card(
              margin: defaults.padding.copyWith(bottom: 0, top: 0),
              Center(
                child: content,
              ),
              onTap: onTap,
              help: widget.help,
              trailing: [
                Center(
                  child: Text(trailingText, textAlign: TextAlign.center),
                ),
              ],
            ),
          ],
        ),
        interval: _interval,
        onTick: _checkLocated,
      ),
    );
  }
}
