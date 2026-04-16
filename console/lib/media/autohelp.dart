import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

class AutoHelp extends StatelessWidget {
  final Widget child;
  const AutoHelp(this.child, {super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ds.HelpAuto(
      child,
      cacheid: 'media',
      title: Text("Media Library", style: theme.textTheme.titleMedium),
      content: const Text(
        "Browse and play your media library. Use the search bar to filter "
        "by title, artist, or album. Click any item to begin playback, "
        "or use the player controls at the bottom to manage the queue.",
      ),
    );
  }
}
