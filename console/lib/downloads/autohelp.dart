import 'dart:async';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/retrovibed.dart' as retro;

class MeteredWarning extends StatefulWidget {
  final Widget child;
  const MeteredWarning(this.child, {super.key});

  @override
  State<MeteredWarning> createState() => _MeteredWarningState();
}

class _MeteredWarningState extends State<MeteredWarning> {
  bool _metered = retro.metered();
  StreamSubscription<List<ConnectivityResult>>? _sub;

  @override
  void initState() {
    super.initState();
    _sub = Connectivity().onConnectivityChanged.listen((results) {
      final m = results.any((r) => r == ConnectivityResult.mobile);
      if (m != _metered) setState(() => _metered = m);
    });
  }

  @override
  void dispose() {
    _sub?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (!_metered) return widget.child;
    final theme = Theme.of(context);
    return ds.HelpAuto(
      widget.child,
      cacheid: 'downloads.metered',
      title: Text("Metered Network", style: theme.textTheme.titleMedium),
      content: const Text(
        "Downloads are heavily restricted on a metered networks. "
        "Connect to an unmetered network to avoid interruptions or unexpected data charges.",
      ),
    );
  }
}

class AutoHelp extends StatelessWidget {
  final Widget child;
  const AutoHelp(this.child, {super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ds.HelpAuto(
      this.child,
      cacheid: 'downloads',
      title: Text("Downloads", style: theme.textTheme.titleMedium),
      content: const Text(
        "Track and manage your active and completed downloads. "
        "Pause, resume, or cancel individual downloads from this view. "
        "Completed downloads are available in your media library.",
      ),
    );
  }
}
