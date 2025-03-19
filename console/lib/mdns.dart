import 'package:flutter/material.dart';
import 'dart:async';
import 'package:multicast_dns/multicast_dns.dart';
import 'package:console/designkit.dart' as ds;
import 'package:console/design.kit/forms.dart' as forms;
import './httpx.dart' as httpx;

class ManualConfiguration extends StatefulWidget {
  final void Function() retry;
  final void Function(String) connect;

  ManualConfiguration({super.key, required this.retry, required this.connect});

  @override
  State<ManualConfiguration> createState() => _ManualConfigurationView();
}

class _ManualConfigurationView extends State<ManualConfiguration> {
  String _hostname = '';

  @override
  Widget build(BuildContext context) {
    return forms.Container(
      child: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            SelectableText(
              textAlign: TextAlign.center,
              "unable to locate retrovibed on your local network, ensure a retrovibed is running or provide the details to a remote instance.",
            ),
            forms.Field(
              label: SelectableText("hostname"),
              input: TextFormField(
                decoration: new InputDecoration(
                  hintText: "example.com:9998",
                  helperText: "hostname and port for the retrovibed instance",
                ),
                keyboardType: TextInputType.number,
                onChanged:
                    (v) => setState(() {
                      _hostname = v;
                    }),
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                TextButton(child: Text("retry"), onPressed: widget.retry),
                TextButton(
                  child: Text("connect"),
                  onPressed: () {
                    widget.connect(_hostname);
                  },
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class MDNSDiscovery extends StatefulWidget {
  final Widget child;
  final UniqueKey id = UniqueKey();

  MDNSDiscovery(this.child, {super.key});

  static _MDNSDiscovery? of(BuildContext context) {
    return context.findAncestorStateOfType<_MDNSDiscovery>();
  }

  @override
  State<StatefulWidget> createState() => _MDNSDiscovery();
}

class _MDNSDiscovery extends State<MDNSDiscovery> {
  static const String ServiceName = "_retrovibed._udp.local";
  bool _loading = true;
  Widget? _cause = null;

  void discover() {
    final MDnsClient _client = MDnsClient();
    final Completer<String> _c = new Completer();
    _client
        .start()
        .then((_) {
          _client
              .lookup<PtrResourceRecord>(
                ResourceRecordQuery.serverPointer(ServiceName),
              )
              .listen((ptr) {
                _client
                    .lookup<SrvResourceRecord>(
                      ResourceRecordQuery.service(ptr.domainName),
                    )
                    .listen((srv) {
                      _c.complete("${srv.target}:${srv.port}");
                    }, onError: _c.completeError);
              }, onError: _c.completeError);
        })
        .catchError((cause) {
          _c.completeError(cause);
        });

    _c.future
        .timeout(
          Duration(seconds: 3),
          onTimeout:
              () => Future.error(TimeoutException("operation timed out.")),
        )
        .then((v) {
          setState(() {
            httpx.set(v);
          });
        })
        .catchError((cause) {
          _loading = false;
          _cause = ConstrainedBox(
            constraints: BoxConstraints(maxHeight: 256, maxWidth: 800),
            child: ManualConfiguration(
              retry: () {
                setState(() {
                  _loading = true;
                  _cause = null;
                });
                this.discover();
              },
              connect: (hostname) {
                setState(() {
                  httpx.set(hostname);
                  _cause = null;
                });
              },
            ),
          );
        }, test: ds.ErrorTests.timeout)
        .catchError((cause) {
          setState(() {
            _loading = false;
            _cause = ds.Error.unknown(cause);
          });
        })
        .whenComplete(() {
          _client.stop();
          setState(() {
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    this.discover();
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(cause: _cause, loading: _loading, child: widget.child);
  }
}
