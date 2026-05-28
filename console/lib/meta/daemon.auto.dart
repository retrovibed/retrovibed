import 'dart:io';
import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/retrovibed.dart' as retro;
import './api.dart' as api;
import './daemon.mdns.dart' as mdns;

class DaemonHttpOverrides extends HttpOverrides {
  DaemonHttpOverrides() {}
  @override
  HttpClient createHttpClient(SecurityContext? context) {
    return super.createHttpClient(context)
      ..badCertificateCallback = (X509Certificate cert, String host, int port) {
        final validated = host == "localhost" || host == Platform.localHostname || retro.validatecert(host, cert.der);
        return validated;
      };
  }
}

class EndpointAuto extends StatefulWidget {
  final Widget child;
  final void Function(api.Daemon v)? onTap;
  final Future<api.DaemonLookupResponse> Function() latest;
  final Future<api.DaemonCreateResponse> Function(api.DaemonCreateRequest) create;
  final Future<api.Daemon> Function(api.Daemon) connectable;
  final Future<api.Session> Function(api.Identity, {String? host}) register;
  final Duration Function(int attempt)? backoff;

  const EndpointAuto(
    this.child, {
    super.key,
    this.latest = api.daemons.latest,
    this.onTap,
    this.create = api.daemons.create,
    this.connectable = api.daemons.connectable,
    this.register = api.register,
    this.backoff,
  });

  static _EndpointAuto? of(BuildContext context) {
    return context.findAncestorStateOfType<_EndpointAuto>();
  }

  @override
  State<StatefulWidget> createState() => _EndpointAuto();
}

class _EndpointAuto extends State<EndpointAuto> with WidgetsBindingObserver {
  final ValueNotifier<api.Daemon> changed = ValueNotifier<api.Daemon>(
    api.Daemon(),
  );
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  api.Daemon? _res;
  Widget Function(void Function(api.Daemon) connect, {void Function()? retry}) _preamble =
      (connect, {retry}) => mdns.NoLocalService(connect: connect, retry: retry);

  Future<void> setdaemon(api.Daemon? d) {
    if (d == null) return Future.value(null);
    return refresh(Future.value(d));
  }

  reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<api.Daemon> latest() {
    return httpx.withRetry(
      () => widget.latest().then((r) => r.daemon).catchError((e) {
        // no service known
        return httpx
            .withRetry(
              () => widget.create(
                api.DaemonCreateRequest(
                  daemon: api.Daemon(hostname: httpx.localhost()),
                ),
              ),
              backoff: widget.backoff,
            )
            .then((v) => v.daemon);
      }, test: httpx.ErrorsTest.err404),
      backoff: widget.backoff,
      checks: const [
        httpx.RetryChecks.missingtoken,
        httpx.RetryChecks.unauthorized,
        ...httpx.RetryChecks.auto,
      ],
    );
  }

  Future<void> refreshNoErrHandling(Future<api.Daemon> pending) {
    return pending
        .then((v) {
          return widget.connectable(v).catchError((e, _) {
            var c = Completer<api.Daemon>();
            setState(() {
              _loading = false;
              _cause = ds.Error.unauthorized(
                e,
                onTap: reseterr,
                color: Color.fromRGBO(0, 0, 0, 0.80),
                message: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      "you've attempted to access a system you havent been granted access to yet.",
                    ),
                    Text("would you like to request access?"),
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        TextButton(onPressed: reseterr, child: Text('No')),
                        ds.LoadingButton(
                          Text("Yes"),
                          onPressed: () {
                            return httpx
                                .withRetry(
                                  () => widget.register(
                                    api.Identity.create()..display = retro.username(),
                                    host: v.hostname,
                                  ),
                                )
                                .then((r) {
                                  c.complete(v);
                                })
                                .catchError((e) {
                                  c.completeError(e);
                                });
                          },
                        ),
                      ],
                    ),
                  ],
                ),
              );
            });
            return c.future;
          }, test: httpx.ErrorsTest.unauthorized);
        })
        .then((v) {
          setState(() {
            httpx.set(v.hostname);
            _res = v;
            changed.value = v;
          });
        });
  }

  Future<void> refresh(Future<api.Daemon> pending) {
    setState(() {
      _loading = true;
    });
    return refreshNoErrHandling(pending)
        .catchError((e) {
          // no service known
          setState(() {
            _preamble =
                (connect, {retry}) => mdns.InitialSetup(
                  connect: (d) => refresh(Future.value(d)),
                  retry: retry,
                );
          });
        }, test: httpx.ErrorsTest.err404)
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unauthorized(
              e,
              onTap: reseterr,
              message: Text(
                "you've attempted to access a service you havent been granted access to yet.",
              ),
            );
          });
        }, test: httpx.ErrorsTest.forbidden)
        .catchError((e) {
          // profile is pending approval.
          setState(() {
            _cause = ds.Error.unauthorized(
              e,
              onTap: reseterr,
              color: Color.fromRGBO(0, 0, 0, 0.80),
              message: Text(
                "you've not yet been approved to access this library",
              ),
            );
          });
        }, test: httpx.ErrorsTest.conflict)
        .catchError((e) {
          // fallback to manual setup.
        }, test: ds.ErrorTests.offline)
        .catchError((e) {
          // fallback to manual setup.
        }, test: ds.ErrorTests.dnsresolution)
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e, onTap: reseterr);
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
    refresh(latest());
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final Widget fallback =
        _loading
            ? ds.Empty
            : mdns.MDNSDiscovery(
              daemon: (d) {
                setState(() {
                  _res = d;
                });
              },
              preamble: _preamble,
            );

    final failed = !(_res == null && _cause == ds.Error.zero);
    return ds.Loading(
      loading: _loading,
      cause: failed ? _cause : fallback,
      widget.child,
    );
  }
}
