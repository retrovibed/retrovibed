import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/ddisc.dart' as ddisc;
import 'package:retrovibed/discovery.dart' as disc;
import 'known.media.card.dart';
import './api.dart' as api;

class KnownMediaLocator extends StatefulWidget {
  final api.Known current;
  final Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate;
  final Future<ddisc.DiscoveryDownloadResponse> Function(String id, {List<httpx.Option> options}) download;
  final Future<bool> Function(BuildContext context, {List<httpx.Option> options}) ensureP2P;
  final IconData icon;
  final Widget help;
  final Widget? trailing;

  const KnownMediaLocator(
    this.current, {
    super.key,
    this.locate = api.locate.create,
    this.download = ddisc.api.download,
    this.ensureP2P = disc.ensureP2P,
    this.icon = Icons.download_rounded,
    this.help = ds.HelpScope.None,
    this.trailing,
  });

  @override
  State<StatefulWidget> createState() => _KnownMediaLocator();
}

class _KnownMediaLocator extends State<KnownMediaLocator> {
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

    widget
        .ensureP2P(context, options: options)
        .then((proceed) {
          if (!proceed) {
            setState(() {
              _loading = false;
            });
            return null;
          }

          switch (widget.current.source) {
            case ddisc.sources.discovered:
              return httpx.withRetry(
                () => widget
                    .download(
                      widget.current.id,
                      options: options,
                    )
                    .then((v) {
                      print("downloading ${v}");
                      setState(() {
                        _queued = true;
                        _loading = false;
                      });
                    }),
              );
            default:
              return httpx.withRetry(
                () => widget
                    .locate(
                      api.Locate.create()..knownMediaId = widget.current.id,
                      options: options,
                    )
                    .then((v) {
                      print("located ${v}");
                      setState(() {
                        _queued = true;
                        _loading = false;
                      });
                    }),
              );
          }
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
    return ds.Loading(
      loading: _loading,
      cause: _cause,
      KnownMediaCard(
        widget.current,
        icon: _queued ? Icons.query_builder_rounded : widget.icon,
        help: widget.help,
        onTap: _loading || _queued ? null : _onTap,
        trailing: widget.trailing,
      ),
    );
  }
}
