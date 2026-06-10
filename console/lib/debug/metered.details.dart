import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/netmonx/api.dart' as api;
import './metered.toggle.dart';

class MeteredDetails extends StatefulWidget {
  final EdgeInsets margin;
  final api.FnNetworkDiagnostics apinetwork;
  const MeteredDetails({super.key, this.margin = EdgeInsets.zero, this.apinetwork = api.network.get});

  @override
  State<MeteredDetails> createState() => _MeteredDetailsState();
}

class _MeteredDetailsState extends State<MeteredDetails> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  api.NetworkMetricsResponse _data = api.NetworkMetricsResponse();

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _fetch());
  }

  void _fetch() {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    final auth = [authn.request(authn.AuthzCache.meta(context))];
    httpx
        .withRetry(() => widget.apinetwork(options: auth))
        .then((v) {
          setState(() {
            _loading = false;
            _data = v;
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
    final wireguard = _data.wireguard;
    final network = _data.network;

    return ds.Card(
      alignment: Alignment.topLeft,
      margin: widget.margin,
      help: ds.Hint(const Text("simulate a metered network connection")),
      SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          spacing: defaults.spacing / 4,
          children: [
            Row(
              children: [
                Expanded(child: Text("Metered Network", style: theme.textTheme.titleMedium)),
                const MeteredToggle(),
              ],
            ),
            ds.Loading(
              loading: _loading,
              cause: _cause,
              Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text("WireGuard", style: theme.textTheme.titleSmall),
                  forms.Field(
                    label: Text("status"),
                    input: Text(wireguard.status.isEmpty ? '—' : wireguard.status),
                  ),
                  forms.Field(
                    label: Text("peer"),
                    input: Text(
                      wireguard.peerKey.isEmpty ? '—' : wireguard.peerKey,
                      overflow: TextOverflow.ellipsis,
                      maxLines: 1,
                    ),
                  ),
                  forms.Field(label: Text("tx"), input: ds.Bytes(wireguard.txBytes)),
                  forms.Field(label: Text("rx"), input: ds.Bytes(wireguard.rxBytes)),
                  forms.Field(
                    label: Text("last handshake"),
                    input: wireguard.lastHandshakeSec > 0
                        ? ds.Duration.until(
                            DateTime.fromMillisecondsSinceEpoch(wireguard.lastHandshakeSec.toInt() * 1000),
                          )
                        : Text("never"),
                  ),
                  Text("Network", style: theme.textTheme.titleSmall),
                  forms.Field(
                    label: Text("ipv4"),
                    input: Text(network.haveV4 ? "yes" : "no"),
                  ),
                  forms.Field(
                    label: Text("ipv6"),
                    input: Text(network.haveV6 ? "yes" : "no"),
                  ),
                  forms.Field(
                    label: Text("default interface"),
                    input: Text(network.defaultInterface.isEmpty ? '—' : network.defaultInterface),
                  ),
                  ...network.interfaces.map(
                    (iface) => forms.Field(
                      label: Text(iface.name),
                      input: Text('${iface.ip}${iface.metered ? " · metered" : ""}'),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
