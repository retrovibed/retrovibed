import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as _meta;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/uuidx.dart' as uuidx;

typedef FnAuthzCurrent = Future<_meta.AuthzResponse> Function({String? host});

class AuthzCache extends StatefulWidget {
  final Widget child;
  final FnAuthzCurrent current;

  const AuthzCache(this.child, {super.key, this.current = _meta.authz.current});

  static Future<_meta.AuthzResponse> fake({String? host}) {
    return Future.value(
      _meta.AuthzResponse(
        bearer: uuidx.min(),
        token: _meta.Token(
          expires: fixnum.Int64(DateTime.now().millisecondsSinceEpoch + 3600000),
        ),
      ),
    );
  }

  static _AuthzCache? of(BuildContext context) {
    return context.findAncestorStateOfType<_AuthzCache>();
  }

  static _meta.Token authzmetadata(BuildContext context) {
    final cache = of(context) ?? _AuthzCache();
    return cache.meta.current.metadata;
  }

  static httpx.Option bearer(BuildContext context) {
    final cache = of(context) ?? _AuthzCache();
    return httpx.Request.bearer(() => cache.meta.token().then((v) => v.bearer));
  }

  static Future<String> bearerString(BuildContext context) {
    final cache = of(context) ?? _AuthzCache();
    return cache.meta.token().then((v) => v.bearer);
  }

  @override
  State<AuthzCache> createState() => _AuthzCache();
}

class _AuthzCache extends State<AuthzCache> {
  bool _loading = true;
  authz.Cached<_meta.Token> meta = authz.Cached(
    authz.Bearer(_meta.Token(), ""),
    authz.Cached.pending,
  );

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void refresh() {
    setState(() {
      meta = authz.Cached(
        authz.Bearer(_meta.Token(), ""),
        authz.refresh(
          (c) => httpx
              .withRetry(
                widget.current,
                checks: const [
                  ...httpx.RetryChecks.auto,
                  httpx.RetryChecks.unauthorized,
                ],
              )
              .then((v) {
                return authz.Bearer(v.token, v.bearer);
              })
              .catchError((e) {
                return authz.Bearer(c, "");
              })
              .whenComplete(() {
                setState(() {
                  _loading = false;
                });
              }),
          (c, ts) {
            return DateTime.fromMillisecondsSinceEpoch(
              c.expires.toInt() * 1000,
              isUtc: true,
            ).isBefore(ts);
          },
        ),
      );
    });
  }

  @override
  void initState() {
    super.initState();
    refresh();
    _meta.EndpointAuto.of(context)?.changed.addListener(refresh);
    meta.refresh(meta);
  }

  @override
  void deactivate() {
    _meta.EndpointAuto.of(context)?.changed.removeListener(refresh);
    super.deactivate();
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(
      loading: _loading,
      _loading ? SizedBox() : widget.child,
    );
  }
}
