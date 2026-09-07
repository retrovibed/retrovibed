import 'package:flutter/material.dart';
import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/meta.dart' as _meta;
import 'cache.dart' as authn;
import 'endpoint.dart';

// Scopes a device and its matching auth token together around a subtree, so
// a caller can target a specific (possibly remote) device without disturbing
// the app-root EndpointAuto/AuthzCache. Owns its ValueNotifier<Daemon>
// internally (mirroring how _meta.EndpointAuto owns `changed`) rather than
// taking one from the caller, so wrapping a subtree is a single,
// stateless-friendly widget.
class AuthedEndpoint extends StatefulWidget {
  final Widget child;
  final _meta.Daemon? initial;
  final authn.FnAuthzCurrent current;

  const AuthedEndpoint(
    this.child, {
    super.key,
    this.initial,
    this.current = _meta.authz.current,
  });

  static ValueNotifier<_meta.Daemon> daemon(BuildContext context) => Endpoint.of(context)!.widget.daemon;

  static authz.Cached<_meta.Token> token(BuildContext context) => _EndpointAuthzCache.meta(context);

  @override
  State<AuthedEndpoint> createState() => _AuthedEndpoint();
}

class _AuthedEndpoint extends State<AuthedEndpoint> {
  final ValueNotifier<_meta.Daemon> _daemon = ValueNotifier(
    _meta.Daemon(),
  );

  @override
  void initState() {
    super.initState();
    _daemon.value = widget.initial ?? _meta.EndpointAuto.of(context)?.changed.value ?? _daemon.value;
    _daemon.addListener(_onDaemonChanged);
  }

  void _onDaemonChanged() => setState(() {});

  @override
  void dispose() {
    _daemon.removeListener(_onDaemonChanged);
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Endpoint(
      daemon: _daemon,
      _EndpointAuthzCache(
        key: ValueKey(_daemon.value.hostname),
        current: ({String? host}) => widget.current(host: host ?? _daemon.value.hostname),
        widget.child,
      ),
    );
  }
}

// A device-scoped counterpart to authn.AuthzCache, kept as its own type (via
// AuthzCache.publish) so it doesn't shadow the app-root AuthzCache for
// descendants that still want the default/global token (e.g. a widget doing
// an unrelated request against httpx.host() while nested under AuthedEndpoint
// for an unrelated remote daemon). authn.AuthzCache.meta(context) walks past
// this InheritedWidget (exact-type lookup) straight to the app-root one; only
// AuthedEndpoint.token resolves this scoped instance.
class _EndpointAuthzCache extends authn.AuthzCache {
  const _EndpointAuthzCache(super.child, {super.key, required super.current});

  static authz.Cached<_meta.Token> meta(BuildContext context) {
    return (context.dependOnInheritedWidgetOfExactType<_ScopedAuthzTokenData>() ?? _ScopedAuthzTokenData.empty).meta;
  }

  @override
  _ScopedAuthzTokenData publish(authz.Cached<_meta.Token> meta, Widget child) =>
      _ScopedAuthzTokenData(meta: meta, child: child);
}

class _ScopedAuthzTokenData extends authn.AuthzTokenData {
  _ScopedAuthzTokenData({required super.meta, required super.child});

  static final empty = _ScopedAuthzTokenData(
    meta: authz.Cached(authz.Bearer(_meta.Token(), ""), authz.Cached.pending),
    child: const SizedBox(),
  );
}
