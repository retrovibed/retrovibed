import 'dart:io';
import 'package:flutter/material.dart';
import 'package:console/designkit.dart' as ds;
import 'package:console/httpx.dart' as httpx;
import './api.dart' as api;
import './daemon.mdns.dart' as mdns;

class DaemonHttpOverrides extends HttpOverrides {
  final List<String> ips;
  DaemonHttpOverrides({this.ips = const []}) {}
  @override
  HttpClient createHttpClient(SecurityContext? context) {
    return super.createHttpClient(context)
      ..badCertificateCallback = (X509Certificate cert, String host, int port) {
        return ips.any((v) => host == v) ||
            host == "localhost" ||
            host == Platform.localHostname;
      };
  }
}

class DaemonAuto extends StatefulWidget {
  final Widget child;
  final void Function(api.Daemon v)? onTap;
  final Future<api.DaemonLookupResponse> Function() latest;

  const DaemonAuto(
    this.child, {
    super.key,
    this.latest = api.daemons.latest,
    this.onTap,
  });

  @override
  State<StatefulWidget> createState() => _DaemonAuto();
}

class _DaemonAuto extends State<DaemonAuto> {
  bool _loading = true;
  ds.Error? _cause = null;
  api.Daemon? _res;

  void refresh() {
    setState(() {
      _loading = true;
    });

    widget
        .latest()
        .then((v) {
          setState(() {
            httpx.set(v.daemon.hostname);
            HttpOverrides.global = DaemonHttpOverrides(
              ips: [v.daemon.hostname.split(":").first],
            );
            _res = v.daemon;
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e);
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    HttpOverrides.global = DaemonHttpOverrides();
    refresh();
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(
      child: _res == null ? mdns.MDNSDiscovery() : widget.child,
      cause: _cause,
      loading: _loading,
    );
  }
}
