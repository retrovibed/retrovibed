import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'ftux/community.picker.dart';
import 'ftux/wireguard.setup.dart';

export 'ftux/api.dart';

/// Shows the first time setup screens (curated community-subscription picker,
/// then wireguard setup) once, the first time this widget is mounted (i.e.
/// right after login, when wrapped around the authenticated app shell).
/// Persists "seen" via [ds.Disclaimer]'s disk cache.
class AutoHelp extends StatefulWidget {
  final Widget child;
  const AutoHelp(this.child, {super.key});

  @override
  State<AutoHelp> createState() => _AutoHelpState();
}

class _AutoHelpState extends State<AutoHelp> {
  static const String _cacheid = 'ftux';
  List<Widget> _pending = const [];

  @override
  void initState() {
    super.initState();
    _pending = [
      CommunityPicker(onDone: _next),
      WireguardSetup(onDone: _next),
    ];
  }

  void _next() {
    if (_pending.length <= 1) ds.Disclaimer.acknowledge(_cacheid);
    setState(() => _pending.removeAt(0));
  }

  @override
  Widget build(BuildContext context) {
    return ds.Disclaimer(
      widget.child,
      cacheid: _cacheid,
      overlay: ds.Masked(
        Center(
          child: SingleChildScrollView(
            child: _pending.isEmpty ? ds.NullWidget : _pending.first,
          ),
        ),
      ),
    );
  }
}
