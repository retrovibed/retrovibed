import 'package:flutter/material.dart';
import 'package:retrovibed/meta.dart' as meta;

// Carries a caller-supplied daemon target down the tree. Unlike
// meta.EndpointAuto this never discovers/creates/registers a daemon and
// never mutates the app's global active host (httpx.set) - it just exposes
// whatever ValueNotifier<Daemon> it's given to descendants.
class Endpoint extends StatefulWidget {
  final Widget child;
  final ValueNotifier<meta.Daemon> daemon;

  const Endpoint(this.child, {super.key, required this.daemon});

  static _Endpoint? of(BuildContext context) {
    return context.findAncestorStateOfType<_Endpoint>();
  }

  @override
  State<Endpoint> createState() => _Endpoint();
}

class _Endpoint extends State<Endpoint> {
  @override
  Widget build(BuildContext context) => widget.child;
}
