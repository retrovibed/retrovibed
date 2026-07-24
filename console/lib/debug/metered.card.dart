import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './metered.details.dart';

class MeteredCard extends StatelessWidget {
  final EdgeInsets? margin;
  final Function(Widget w) onPressed;
  const MeteredCard({super.key, required this.onPressed, this.margin});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    final tap = () => onPressed(
      ds.Container(
        alignment: Alignment.topCenter,
        padding: defaults.padding,
        margin: EdgeInsets.zero,
        const MeteredDetails(),
      ),
    );
    return ds.Card(
      alignment: Alignment.center,
      margin: margin ?? defaults.margin,
      onTap: tap,
      help: ds.Hint(const Text("simulate a metered network connection")),
      Column(
        // crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing / 2,
        children: [
          Text("Metered Network", textAlign: TextAlign.center, style: theme.textTheme.titleMedium),
          Text(
            "Simulate and inspect metered network conditions",
            textAlign: TextAlign.center,
            style: theme.textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}
