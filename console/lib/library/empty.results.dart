import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

// EmptyResults is shown in place of the library grid/list when a search
// turns up nothing locally, offering to switch into discovery mode to
// search the wider network via the same suggestions.
class EmptyResults extends StatelessWidget {
  final VoidCallback onDiscover;

  const EmptyResults({super.key, required this.onDiscover});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return Center(
      child: Padding(
        padding: defaults.padding,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              "no results in your library, would you like to attempt to discover results?",
              textAlign: TextAlign.center,
              style: TextStyle(
                color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 8),
            ElevatedButton.icon(
              onPressed: onDiscover,
              icon: const Icon(Icons.travel_explore_rounded),
              label: const Text("discover"),
            ),
          ],
        ),
      ),
    );
  }
}
