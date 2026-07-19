import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/design.kit/bytesx.dart';
import 'package:retrovibed/ddisc.dart' as ddisc;

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
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
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
      ds.Card(
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              widget.current.title,
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
              style: Theme.of(context).textTheme.titleMedium,
            ),
            if (widget.current.description.isNotEmpty)
              Text(
                widget.current.description,
                maxLines: 4,
                overflow: TextOverflow.ellipsis,
              ),
          ],
        ),
        onTap: _queued || _loading ? null : _onTap,
        help: widget.help,
        trailing: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(size),
              Icon(_queued ? Icons.query_builder_rounded : Icons.download_rounded),
            ],
          ),
        ],
      ),
    );
  }
}
