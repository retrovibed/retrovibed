import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/discovery/search.button.dart';

// EmptyResults is shown in place of the library grid/list when a search
// turns up nothing locally, offering to locate suggestions from the wider
// network via SearchButton.
class EmptyResults extends StatelessWidget {
  final media.MediaSearchState search;

  const EmptyResults({super.key, required this.search});

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
            SearchButton(search: search),
          ],
        ),
      ),
    );
  }
}
