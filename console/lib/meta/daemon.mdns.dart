import 'package:flutter/material.dart';
import 'dart:async';
import 'dart:io';
import 'package:multicast_dns/multicast_dns.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;
import './daemon.manual.dart';
import 'package:retrovibed/design.kit/stateful.dart';

class MDNSDiscovery extends StatefulWidget {
  final void Function(api.Daemon) daemon;
  final Widget Function(
    void Function(api.Daemon) connect, {
    void Function()? retry,
  })
  preamble;
  final Future<String> Function() discover;

  MDNSDiscovery({
    super.key,
    required this.daemon,
    required this.preamble,
    this.discover = autodiscover,
  });

  static _MDNSDiscovery? of(BuildContext context) {
    return context.findAncestorStateOfType<_MDNSDiscovery>();
  }

  static MDnsClient _autoclient() {
    if (Platform.isAndroid) {
      return MDnsClient(
        rawDatagramSocketFactory: (
          dynamic host,
          int port, {
          bool reuseAddress = false,
          bool reusePort = false,
          int ttl = 255,
        }) {
          return RawDatagramSocket.bind(
            host,
            port,
            reuseAddress: reuseAddress,
            reusePort: false,
            ttl: ttl,
          );
        },
      );
    }

    return MDnsClient();
  }

  static const String _serviceName = "_retrovibed._udp.local";

  static Future<String> autodiscover() {
    final MDnsClient client = _autoclient();
    final Completer<String> c = Completer();
    client
        .start()
        .then((_) {
          client
              .lookup<PtrResourceRecord>(
                ResourceRecordQuery.serverPointer(_serviceName),
              )
              .listen((ptr) {
                client
                    .lookup<SrvResourceRecord>(
                      ResourceRecordQuery.service(ptr.domainName),
                    )
                    .listen((srv) {
                      c.complete("${srv.target}:${srv.port}");
                    }, onError: c.completeError);
              }, onError: c.completeError);
        })
        .catchError((cause) {
          c.completeError(cause);
        });

    return c.future
        .timeout(
          Duration(seconds: 3),
          onTimeout: () => Future.error(TimeoutException("operation timed out.")),
        )
        .whenComplete(() => client.stop());
  }

  @override
  State<StatefulWidget> createState() => _MDNSDiscovery();
}

class _MDNSDiscovery extends State<MDNSDiscovery> with LoadingState {
  bool _loading = true;
  Widget _cause = ds.Error.zero;

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  void _runDiscover() {
    widget
        .discover()
        .then((v) {
          setState(() {
            httpx.set(v);
          });
        })
        .catchError((cause) {
          // ignore timeouts — manual setup will handle it.
        }, test: ds.ErrorTests.timeout)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        })
        .whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    _runDiscover();
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(
      cause: _cause,
      loading: _loading,
      Center(
        child: widget.preamble(
          (daemon) {
            setState(() {
              _cause = ds.Error.zero;
            });
            widget.daemon(daemon);
          },
          retry: () {
            setState(() {
              _loading = true;
              _cause = ds.Error.zero;
            });
            _runDiscover();
          },
        ),
      ),
    );
  }
}

class NoLocalService extends StatelessWidget {
  final void Function()? retry;
  final void Function(api.Daemon) connect;

  NoLocalService({super.key, this.retry, required this.connect});

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        SelectableText(
          textAlign: TextAlign.center,
          "unable to locate retrovibed on your local network, ensure retrovibed is running or provide the details to a remote instance.",
        ),
        ManualConfiguration(retry: retry, connect: connect),
      ],
    );
  }
}

class InitialSetup extends StatelessWidget {
  final void Function()? retry;
  final void Function(api.Daemon) connect;

  InitialSetup({super.key, this.retry, required this.connect});

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        SelectableText(
          textAlign: TextAlign.center,
          "provide a server address or click connect to use your current device",
        ),
        ManualConfiguration(connect: connect),
      ],
    );
  }
}
