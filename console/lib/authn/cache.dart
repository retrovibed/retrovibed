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
    return fakeWith(_meta.Token())(host: host);
  }

  static FnAuthzCurrent fakeWith(_meta.Token token) {
    return ({String? host}) {
      return Future.value(
        _meta.AuthzResponse(
          bearer: uuidx.min(),
          token: token..expires = fixnum.Int64(DateTime.now().millisecondsSinceEpoch + 3600000),
        ),
      );
    };
  }

  static _AuthzCache of(BuildContext context) {
    return context.findAncestorStateOfType<_AuthzCache>() ?? _AuthzCache();
  }

  static AuthzTokenData cached(BuildContext context) {
    return context.dependOnInheritedWidgetOfExactType<AuthzTokenData>() ?? AuthzTokenData.empty;
  }

  static _meta.Token authzmetadata(BuildContext context) => cached(context).meta.current.token;

  static authz.Cached<_meta.Token> meta(BuildContext context) => cached(context).meta;

  @override
  State<AuthzCache> createState() => _AuthzCache();

  // Hook allowing a subclass (e.g. a daemon-scoped variant) to publish its
  // token under a distinct InheritedWidget type, so plain AuthzCache.meta()
  // lookups from descendants skip past it and keep resolving the nearest
  // "real" AuthzCache ancestor instead of this scoped one.
  AuthzTokenData publish(authz.Cached<_meta.Token> meta, Widget child) => AuthzTokenData(meta: meta, child: child);
}

class AuthzTokenData extends InheritedWidget {
  final authz.Cached<_meta.Token> meta;

  const AuthzTokenData({required this.meta, required super.child});

  static final empty = AuthzTokenData(
    meta: authz.Cached(authz.Bearer(_meta.Token(), ""), authz.Cached.pending),
    child: const SizedBox(),
  );

  @override
  bool updateShouldNotify(AuthzTokenData old) => meta != old.meta;
}

class _AuthzCache extends State<AuthzCache> {
  bool _loading = true;
  authz.Cached<_meta.Token> meta = authz.Cached(
    authz.Bearer(_meta.Token(), ""),
    authz.Cached.pending,
  );
  final ValueNotifier<authz.Bearer<_meta.Token>> changed = ValueNotifier<authz.Bearer<_meta.Token>>(
    authz.Bearer(_meta.Token(), ""),
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
                final bearer = authz.Bearer(v.token, v.bearer);
                setState(() {
                  meta.current = bearer;
                  changed.value = bearer;
                  _loading = false;
                });
                return bearer;
              })
              .catchError((e) {
                setState(() {
                  _loading = false;
                });
                return authz.Bearer(c, "");
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
  void dispose() {
    changed.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return widget.publish(
      meta,
      ds.LoadingBoundary(
        loading: _loading,
        _loading ? SizedBox() : widget.child,
      ),
    );
  }
}
