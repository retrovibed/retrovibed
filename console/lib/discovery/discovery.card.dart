import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/design.kit/bytesx.dart';
import 'package:retrovibed/ddisc.dart' as ddisc;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'discovery.details.dart';

class DiscoveredCard extends StatefulWidget {
  final ddisc.Discovery current;
  final Future<ddisc.DiscoveryDownloadResponse> Function(String id, {List<httpx.Option> options}) download;
  final void Function(ddisc.Discovery current) onDownloaded;
  final Widget help;

  const DiscoveredCard(
    this.current, {
    super.key,
    this.download = ddisc.api.download,
    this.onDownloaded = ds.fnNoop,
    this.help = ds.HelpScope.None,
  });

  @override
  State<DiscoveredCard> createState() => _DiscoveredCardState();
}

class _DiscoveredCardState extends State<DiscoveredCard> {
  bool _loading = false;
  bool _queued = false;
  bool _resolved = false;
  Widget _cause = ds.Error.zero;

  late lib.Known _known = lib.Known(
    id: "",
    description: widget.current.title,
    summary: widget.current.description,
    rating: 0.0,
    image: "",
  );

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (_resolved) return;
    if (uuidx.isMinMax(uuidx.fromString(widget.current.knownMediaId))) return;
    _resolved = true;

    final authz = authn.AuthzCache.meta(context);
    lib.known
        .cached(
          widget.current.knownMediaId,
          () => lib.known.get(widget.current.knownMediaId, options: [authn.request(authz)]),
        )
        .then((w) => setState(() => _known = w.known..description = widget.current.title));
  }

  void _onTap() {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    final options = [authn.request(authn.AuthzCache.meta(context))];

    httpx
        .withRetry(() => widget.download(widget.current.id, options: options))
        .then((v) {
          widget.onDownloaded(widget.current);
          setState(() {
            _queued = true;
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _loading = false;
            _cause = ds.Errors.httpauto(e, onTap: _reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((e) {
          setState(() {
            _loading = false;
            _cause = ds.Error.unknown(e, onTap: _reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final size = bytesx(widget.current.bytes.toInt()).toIEC600272Format();

    return ds.Loading(
      loading: _loading,
      cause: _cause,
      lib.KnownMediaCard(
        _known,
        overlay: DiscoveryDetails(widget.current, _known),
        icon: _queued ? Icons.query_builder_rounded : Icons.download_rounded,
        help: widget.help,
        onTap: _queued || _loading ? null : _onTap,
        trailing: Text(size),
      ),
    );
  }
}
