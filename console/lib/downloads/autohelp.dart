import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

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
        "Completed downloads are available in your media library.\n\n"
        "Press Alt+? at any time to activate/deactivate help overlay",
      ),
    );
  }
}
