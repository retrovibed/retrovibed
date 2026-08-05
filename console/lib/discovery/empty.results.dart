import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

// EmptyResults is shown in place of the discovery grid when a network
// search turns up nothing, offering a CTA to queue a peer-to-peer search.
class EmptyResults extends StatelessWidget {
  final Widget activation;

  const EmptyResults(this.activation, {super.key});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return Center(
      child: Padding(
        padding: defaults.padding,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          spacing: defaults.spacing,
          children: [
            Text(
              "no candidates found on the network for this search",
              textAlign: TextAlign.center,
              style: TextStyle(
                color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
              ),
            ),
            activation,
          ],
        ),
      ),
    );
  }
}
