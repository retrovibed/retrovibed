import 'package:flutter/material.dart';
import 'package:retrovibed/meta.dart' as _meta;
import 'package:retrovibed/meta/api.deeppool.dart' as deeppool;
import 'package:retrovibed/billing/api.dart' as billing;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authz.dart' as authz;
import 'authenticated.dart';
import 'cache.dart';

typedef FnDeeppoolAuthz = Future<deeppool.AuthzResponse> Function({List<httpx.Option> options});
typedef FnAttribution = Future<billing.AttributionTokenResponse> Function({List<httpx.Option> options});

// Decides whether DeeppoolAuthzCache should even be mounted: a local_only
// (guest) identity should never make the deeppool authz/billing calls that
// widget performs, so it's simplest to never create its state at all rather
// than guard each call inside it.
//
// Listens to AuthzCache's `changed` notifier rather than depending on
// AuthzTokenData (InheritedWidget) directly - refresh() only ever reassigns
// meta.current in place, not the AuthzTokenData.meta reference itself, so
// InheritedWidget dependents aren't notified past the first resolved token.
// `changed.value = bearer` fires correctly on every refresh, so it's the
// reliable signal to rebuild on here.
class DeeppoolAuthzCacheGuard extends StatelessWidget {
  final Widget child;
  const DeeppoolAuthzCacheGuard(this.child, {super.key});

  @override
  Widget build(BuildContext context) {
    final cache = AuthzCache.of(context);
    return ValueListenableBuilder<authz.Bearer<_meta.Token>>(
      valueListenable: cache.changed,
      builder: (context, bearer, _) {
        if (bearer.token.localOnly) {
          return child;
        }
        return DeeppoolAuthzCache(child);
      },
    );
  }
}

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
