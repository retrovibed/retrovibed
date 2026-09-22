import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/ddisc.dart' as ddisc;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'discovery.details.dart';
import 'locate.p2p.prompt.dart' as p2p;

class DiscoveredCard extends StatefulWidget {
  final ddisc.Discovery current;
  final Future<ddisc.DiscoveryDownloadResponse> Function(
    String id, {
    ddisc.Discovery? discovery,
    bool autodownload,
    List<httpx.Option> options,
  })
  download;
  final Future<void> Function(BuildContext context, {List<httpx.Option> options}) ensureP2P;
  final void Function(ddisc.Discovery current) onDownloaded;
  final Widget help;

  const DiscoveredCard(
    this.current, {
    super.key,
    this.download = ddisc.api.download,
    this.ensureP2P = p2p.ensureP2P,
    this.onDownloaded = ds.fnNoop,
    this.help = ds.HelpScope.None,
  });

  @override
  State<DiscoveredCard> createState() => _DiscoveredCardState();
}

class _DiscoveredCardState extends State<DiscoveredCard> with ds.LoadingState {
  bool _queued = false;
  bool _resolved = false;

  late lib.Known _known = lib.Known(
    id: "",
    description: widget.current.title,
    summary: widget.current.description,
    rating: 0.0,
    image: "",
  );

  @override
  void initState() {
    super.initState();
    loading = false;
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
      loading = true;
      cause = ds.Error.zero;
    });

    final options = [authn.request(authn.AuthzCache.meta(context))];

    widget
        .ensureP2P(context, options: options)
        .then(
          (_) => httpx.withRetry(
            () => widget.download(widget.current.id, discovery: widget.current, autodownload: true, options: options),
          ),
        )
        .then((v) {
          widget.onDownloaded(widget.current);
          setState(() {
            _queued = true;
            loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            loading = false;
          });
        }, test: p2p.consentDeclined)
        .catchError((e) {
          setState(() {
            loading = false;
            cause = ds.Errors.httpauto(e, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((e) {
          setState(() {
            loading = false;
            cause = ds.Error.unknown(e, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(
      loading: loading,
      cause: cause,
      lib.KnownMediaCard(
        _known,
        overlay: DiscoveryDetails(widget.current, _known),
        icon: _queued ? Icons.query_builder_rounded : Icons.download_rounded,
        help: widget.help,
        onTap: _queued || loading ? null : _onTap,
        trailing: [
          ds.Bytes(widget.current.bytes),
          Spacer(),
          ds.Timestamp.iso8601(
            _known.released,
            format: ds.Timestamp.year,
            inf: ds.Empty,
            neginf: ds.Empty,
          ),
          Spacer(),
          Text("${widget.current.health}"),
        ],
      ),
    );
  }
}
