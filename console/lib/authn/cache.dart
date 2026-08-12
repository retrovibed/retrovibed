import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as _meta;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/design.kit/stateful.dart';

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

  static _AuthzTokenData cached(BuildContext context) {
    return context.dependOnInheritedWidgetOfExactType<_AuthzTokenData>() ?? _AuthzTokenData.empty;
  }

  static _meta.Token authzmetadata(BuildContext context) => cached(context).meta.current.token;

  static authz.Cached<_meta.Token> meta(BuildContext context) => cached(context).meta;

  @override
  State<AuthzCache> createState() => _AuthzCache();
}

class _AuthzTokenData extends InheritedWidget {
  final authz.Cached<_meta.Token> meta;

  const _AuthzTokenData({required this.meta, required super.child});

  static final empty = _AuthzTokenData(
    meta: authz.Cached(authz.Bearer(_meta.Token(), ""), authz.Cached.pending),
    child: const SizedBox(),
  );

  @override
  bool updateShouldNotify(_AuthzTokenData old) => meta != old.meta;
}

class _AuthzCache extends State<AuthzCache> with LoadingState {
  bool _loading = true;
  authz.Cached<_meta.Token> meta = authz.Cached(
    authz.Bearer(_meta.Token(), ""),
    authz.Cached.pending,
  );

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
  Widget build(BuildContext context) {
    return _AuthzTokenData(
      meta: meta,
      child: ds.LoadingBoundary(
        loading: _loading,
        _loading ? SizedBox() : widget.child,
      ),
    );
  }
}
