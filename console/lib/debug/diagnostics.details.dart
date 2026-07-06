import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/ddisc/api.dart' as ddisc;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/dhtx/api.dart' as dhtx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/torrentx/api.dart' as torrentx;

class DiagnosticsDetails extends StatefulWidget {
  final EdgeInsets margin;
  final torrentx.FnTorrentDiagnostics apitorrent;
  final dhtx.FnDHTDiagnostics apidht;
  final ddisc.FnDiscoveryDiagnostics apidiscovery;
  const DiagnosticsDetails({
    super.key,
    this.margin = EdgeInsets.zero,
    this.apitorrent = torrentx.diagnostics.get,
    this.apidht = dhtx.diagnostics.get,
    this.apidiscovery = ddisc.diagnostics.get,
  });

  @override
  State<DiagnosticsDetails> createState() => _DiagnosticsDetailsState();
}

class _DiagnosticsDetailsState extends State<DiagnosticsDetails> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  torrentx.TorrentMetricsResponse _torrent = torrentx.TorrentMetricsResponse();
  dhtx.DHTMetricsResponse _dht = dhtx.DHTMetricsResponse();
  ddisc.DiscoveryMetricsResponse _discovery = ddisc.DiscoveryMetricsResponse();

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(() => _fetch());
  }

  void _fetch() {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    final auth = [authn.request(authn.AuthzCache.meta(context))];
    Future.wait([
          httpx.withRetry(() => widget.apitorrent(options: auth)),
          httpx.withRetry(() => widget.apidht(options: auth)),
          httpx.withRetry(() => widget.apidiscovery(options: auth)),
        ])
        .then((v) {
          setState(() {
            _loading = false;
            _torrent = v[0] as torrentx.TorrentMetricsResponse;
            _dht = v[1] as dhtx.DHTMetricsResponse;
            _discovery = v[2] as ddisc.DiscoveryMetricsResponse;
          });
        })
        .catchError((e) {
          setState(() {
            _loading = false;
            _cause = ds.Errors.httpauto(e, onTap: _fetch);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((e) {
          setState(() {
            _loading = false;
            _cause = ds.Error.unknown(e, onTap: _fetch);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final torrent = _torrent.torrent;
    final dht = _dht.dht;
    final discovery = _discovery.discovery;

    return ds.Card(
      alignment: Alignment.topLeft,
      margin: widget.margin,
      help: ds.Hint(const Text("dht, discovery, and various other subsystem diagnostics")),
      SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          spacing: defaults.spacing / 4,
          children: [
            Row(
              children: [
                Expanded(child: Text("Diagnostics", style: theme.textTheme.titleMedium)),
                ds.LoadingIconButton(
                  icon: const Icon(Icons.bug_report),
                  tooltip: "throw a test error",
                  onPressed: () async {
                    await Future<void>(() => throw Exception("synthetic diagnostics error")).catchError((e) {
                      setState(() {
                        _cause = ds.Error.unknown(e, onTap: _fetch);
                      });
                    });
                  },
                ),
              ],
            ),
            ds.Loading(
              loading: _loading,
              cause: _cause,
              Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text("Torrent", style: theme.textTheme.titleSmall),
                  forms.Field(label: const Text("total"), input: Text("${torrent.total}")),
                  forms.Field(label: const Text("seeding"), input: Text("${torrent.seeding}")),
                  forms.Field(label: const Text("bytes"), input: ds.Bytes(torrent.bytes)),
                  forms.Field(label: const Text("downloaded"), input: ds.Bytes(torrent.downloaded)),
                  forms.Field(label: const Text("uploaded"), input: ds.Bytes(torrent.uploaded)),
                  forms.Field(label: const Text("peers"), input: Text("${torrent.peers}")),
                  Text("DHT", style: theme.textTheme.titleSmall),
                  forms.Field(label: const Text("good nodes"), input: Text("${dht.goodNodes}")),
                  forms.Field(label: const Text("nodes"), input: Text("${dht.nodes}")),
                  forms.Field(label: const Text("bad nodes"), input: Text("${dht.badNodes}")),
                  forms.Field(
                    label: const Text("outstanding transactions"),
                    input: Text("${dht.outstandingTransactions}"),
                  ),
                  forms.Field(
                    label: const Text("successful announces"),
                    input: Text("${dht.successfulOutboundAnnouncePeerQueries}"),
                  ),
                  forms.Field(label: const Text("queries attempted"), input: Text("${dht.outboundQueriesAttempted}")),
                  Text("Discovery", style: theme.textTheme.titleSmall),
                  forms.Field(label: const Text("enabled"), input: Text(discovery.enabled ? "yes" : "no")),
                  forms.Field(label: const Text("ratio"), input: Text("${discovery.ratio}")),
                  forms.Field(label: const Text("partitions"), input: Text("${discovery.partitions}")),
                  forms.Field(label: const Text("workloads"), input: Text("${discovery.workloads}")),
                  forms.Field(
                    label: const Text("local partition"),
                    input: Text(
                      discovery.localPartition.isEmpty ? '—' : discovery.localPartition,
                      overflow: TextOverflow.ellipsis,
                      maxLines: 1,
                    ),
                  ),
                  forms.Field(label: const Text("peers"), input: Text("${discovery.peers}")),
                  forms.Field(label: const Text("ddisc peers"), input: Text("${discovery.peersDdisc}")),
                  forms.Field(label: const Text("bep51 peers"), input: Text("${discovery.peersBep51}")),
                  forms.Field(label: const Text("unknown"), input: Text("${discovery.indexing}")),
                  forms.Field(label: const Text("queued"), input: Text("${discovery.unidentified}")),
                  forms.Field(label: const Text("offload"), input: Text("${discovery.offload}")),
                  forms.Field(label: const Text("indexed"), input: Text("${discovery.indexed}")),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
