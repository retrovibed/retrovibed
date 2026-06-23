import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './diagnostics.details.dart';

class DiagnosticsCard extends StatelessWidget {
  final EdgeInsets? margin;
  final Function(Widget w) onPressed;
  const DiagnosticsCard({super.key, required this.onPressed, this.margin});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    final tap = () => onPressed(
      ds.Container(
        alignment: Alignment.topCenter,
        padding: defaults.padding,
        margin: EdgeInsets.zero,
        const DiagnosticsDetails(),
      ),
    );
    return ds.Card(
      alignment: Alignment.center,
      margin: margin ?? defaults.margin,
      onTap: tap,
      help: ds.Hint(const Text("dht, discovery, and various other subsystem diagnostics")),
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing / 2,
        children: [
          Text("Diagnostics", style: theme.textTheme.titleMedium),
          Text(
            "Inspect system diagnostics",
            style: theme.textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}
