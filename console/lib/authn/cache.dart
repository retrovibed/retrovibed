import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as _meta;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/uuidx.dart' as uuidx;

typedef FnAuthzCurrent = Future<_meta.AuthzResponse> Function({String? host});

// Placeholder token used while the real /meta/authz/ fetch is still in
// flight. Defaults local_only to true so anything reading the token before
// the real one lands fails safe (assumes guest, makes no remote calls)
// instead of fails open.
_meta.Token _pendingToken() => _meta.Token(
  localOnly: true,
);

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
    meta: authz.Cached(authz.Bearer(_pendingToken(), ""), authz.Cached.pending),
    child: const SizedBox(),
  );

  @override
  bool updateShouldNotify(AuthzTokenData old) => meta != old.meta;
}

class _AuthzCache extends State<AuthzCache> with ds.LoadingState {
  authz.Cached<_meta.Token> meta = authz.Cached(
    authz.Bearer(_pendingToken(), ""),
    authz.Cached.pending,
  );
  final ValueNotifier<authz.Bearer<_meta.Token>> changed = ValueNotifier<authz.Bearer<_meta.Token>>(
    authz.Bearer(_pendingToken(), ""),
  );

  // Mutates the existing Cached in place (current + refresh fn) rather than
  // replacing the `meta` field with a new instance. `Cached.auto()`'s lock
  // serializes overlapping fetches against this single shared object, so a
  // stale fetch started before this call can never resolve later and
  // clobber the newer result — there is no second object for it to miss.
  Future<authz.Bearer<_meta.Token>> refresh() {
    setState(() {
      meta.current = authz.Bearer(_pendingToken(), "");
    });
    return meta.auto();
  }

  @override
  void initState() {
    super.initState();
    meta.refresh = authz.refresh(
      (c) => httpx
          .withRetry(
            widget.current,
            checks: const [
              ...httpx.RetryChecks.auto,
              httpx.RetryChecks.unauthorized,
            ],
            // we're talking to the local instance cheap.
            backoff: httpx.Backoff.constant(const Duration(milliseconds: 200)),
          )
          .then((v) {
            final bearer = authz.Bearer(v.token, v.bearer);
            setState(() {
              meta.current = bearer;
              changed.value = bearer;
              loading = false;
              cause = ds.Error.zero;
            });
            return bearer;
          })
          .catchError((e) {
            setState(() {
              loading = false;
              cause = ds.Errors.httpauto(e, onTap: reseterr);
            });
            // rethrow instead of caching an empty bearer: leave
            // `meta.current` untouched so the next fetch retries rather
            // than handing callers a token that's guaranteed to fail.
            throw e;
          }),
      (c, ts) {
        return DateTime.fromMillisecondsSinceEpoch(
          c.expires.toInt() * 1000,
          isUtc: true,
        ).isBefore(ts);
      },
    );

    _meta.EndpointAuto.of(context)?.changed.addListener(refresh);
    // fire-and-forget kickoff, routed through auto() so it shares the same
    // lock as every other caller instead of invoking the fetch directly:
    // success/failure are both already handled via setState inside
    // refresh()'s catchError/then, so just log and swallow the rejection
    // here to avoid an unhandled-future warning.
    meta.auto().catchError((e) {
      print("failed to refresh token cache ${e}");
      return authz.Bearer(_pendingToken(), "");
    });
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
        loading: loading,
        cause: cause,
        origin: '_AuthzCache',
        widget.child,
      ),
    );
  }
}
