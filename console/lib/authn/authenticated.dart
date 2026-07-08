import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'api.dart' as api;

class Authenticated extends StatefulWidget {
  final Widget child;
  final Future<api.Authed> Function() apissh;
  final Future<api.Session> Function() apisignup;
  final Future<api.Session> Function(String token) apicurrent;

  const Authenticated(
    this.child, {
    super.key,
    this.apissh = api.ssh,
    this.apisignup = api.signup,
    this.apicurrent = api.current,
  });

  static api.Session syncSession(BuildContext context) {
    return context.findAncestorStateOfType<_AuthenticatedState>()?.syncCurrent() ?? api.Session();
  }

  static Future<api.Session> session(BuildContext context) {
    return context.findAncestorStateOfType<_AuthenticatedState>()?.current() ?? Future.value(api.Session());
  }

  static httpx.Option bearer(BuildContext context) {
    return httpx.Request.bearer(
      () => session(context).then((s) {
        return s.token;
      }),
    );
  }

  @override
  State<Authenticated> createState() => _AuthenticatedState();
}

class _AuthenticatedState extends State<Authenticated> {
  Widget _cause = ds.Error.zero;
  DateTime _expires = DateTime.timestamp();
  api.Session _current = api.Session();
  bool _loading = true;

  api.Session syncCurrent() {
    return _current;
  }

  Future<api.Session> current() {
    if (_expires.isAfter(DateTime.timestamp())) {
      return Future.value(_current);
    }

    return widget
        .apissh()
        .then<api.Session>((v) {
          switch (v.profiles.length) {
            case 0:
              return widget.apisignup();
            case 1:
              return widget.apicurrent(v.profiles.first.token);
            default:
              return Future.error(
                new Exception("multiple profiles not currently supported"),
              );
          }
        })
        .then((v) {
          setState(() {
            _expires = DateTime.fromMillisecondsSinceEpoch(
              (Duration(seconds: v.expires.toInt()) - Duration(seconds: 60)).inMilliseconds,
            );
            _current = v;
          });
          return v;
        });
  }

  void _refresh() {
    current()
        .then((v) {
          setState(() {
            _loading = false;
          });
        })
        .catchError((cause) {
          setState(() {
            _loading = false;
            _cause = ds.Errors.httpauto(cause, onTap: _reseterr);
          });
        });
  }

  void _reseterr() {
    _refresh();
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    _refresh();
  }

  @override
  Widget build(BuildContext context) {
    return ds.LoadingBoundary(
      widget.child,
      loading: _loading,
      cause: _cause,
    );
  }
}
