import 'package:flutter/material.dart';
import 'package:retrovibed/authz.dart' as authz;
import 'package:retrovibed/meta.dart' as meta;
import 'cache.dart' as authn;
import 'endpoint.dart';

// Scopes a device and its matching auth token together around a subtree, so
// a caller can target a specific (possibly remote) device without disturbing
// the app-root EndpointAuto/AuthzCache. Owns its ValueNotifier<Daemon>
// internally (mirroring how meta.EndpointAuto owns `changed`) rather than
// taking one from the caller, so wrapping a subtree is a single,
// stateless-friendly widget.
class AuthedEndpoint extends StatefulWidget {
  final Widget child;
  final meta.Daemon? initial;
  final authn.FnAuthzCurrent current;

  const AuthedEndpoint(
    this.child, {
    super.key,
    this.initial,
    this.current = meta.authz.current,
  });

  static ValueNotifier<meta.Daemon> daemon(BuildContext context) => Endpoint.of(context)!.widget.daemon;

  static authz.Cached<meta.Token> token(BuildContext context) => authn.AuthzCache.meta(context);

  @override
  State<AuthedEndpoint> createState() => _AuthedEndpoint();
}

class _AuthedEndpoint extends State<AuthedEndpoint> {
  late final ValueNotifier<meta.Daemon> _daemon = ValueNotifier(
    widget.initial ?? meta.EndpointAuto.of(context)?.changed.value ?? meta.Daemon(),
  );

  @override
  void initState() {
    super.initState();
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
      authn.AuthzCache(
        key: ValueKey(_daemon.value.hostname),
        current: ({String? host}) => widget.current(host: host ?? _daemon.value.hostname),
        widget.child,
      ),
    );
  }
}
