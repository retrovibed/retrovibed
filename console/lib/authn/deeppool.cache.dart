import 'package:flutter/material.dart';
import 'package:retrovibed/meta.dart' as _meta;
import 'package:retrovibed/meta/api.deeppool.dart' as deeppool;
import 'package:retrovibed/billing/api.dart' as billing;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authz.dart' as authz;
import 'authenticated.dart';

typedef FnDeeppoolAuthz = Future<deeppool.AuthzResponse> Function({List<httpx.Option> options});
typedef FnAttribution = Future<billing.AttributionTokenResponse> Function({List<httpx.Option> options});

class DeeppoolAuthzCache extends StatefulWidget {
  final Widget child;
  final FnDeeppoolAuthz apideeppoolauthz;
  final FnAttribution apibillingattribution;

  const DeeppoolAuthzCache(
    this.child, {
    Key? key,
    this.apideeppoolauthz = deeppool.authz,
    this.apibillingattribution = billing.attribution,
  }) : super(key: key);

  static _AuthzCache of(BuildContext context) {
    return context.findAncestorStateOfType<_AuthzCache>() ?? _AuthzCache();
  }

  static httpx.Option bearer(BuildContext context) {
    final cache = of(context);
    return httpx.Request.bearer(() => cache.meta.auto().then((v) => v.bearer));
  }

  static Future<String> attributionToken(BuildContext context) {
    final cache = of(context);
    return cache._attributionToken();
  }

  @override
  State<DeeppoolAuthzCache> createState() => _AuthzCache();
}

class _AuthzCache extends State<DeeppoolAuthzCache> {
  authz.Cached<_meta.Token> meta = authz.Cached(
    authz.Bearer(_meta.Token(), ""),
    authz.Cached.pending,
  );
  String? _attribution;

  Future<String> _attributionToken() {
    return httpx
        .withRetry(
          () => widget.apibillingattribution(options: [Authenticated.bearer(context)]),
        )
        .then((v) => _attribution ??= v.token);
  }

  Future<authz.Bearer<_meta.Token>> Function(authz.Cached<_meta.Token>) _refresh() {
    return authz.refresh(
      (c) => widget
          .apideeppoolauthz(options: [Authenticated.bearer(context)])
          .then((v) {
            return authz.Bearer(v.token, v.bearer);
          })
          .then((v) {
            setState(() {
              meta = authz.Cached(
                v,
                _refresh(),
              );
            });
            return v;
          })
          .catchError((e) {
            print("failed to refresh token cache ${e}");
            return authz.Bearer(c, "");
          }),
      (c, ts) {
        return DateTime.fromMillisecondsSinceEpoch(
          c.expires.toInt() * 1000,
          isUtc: true,
        ).isBefore(ts);
      },
    );
  }

  @override
  void initState() {
    super.initState();
    meta = authz.Cached(
      authz.Bearer(_meta.Token(), ""),
      _refresh(),
    );
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return widget.child;
  }
}
