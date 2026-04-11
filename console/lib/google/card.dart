import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/google/settings.dart';

class Card extends StatelessWidget {
  final EdgeInsets margin;
  final Function(Widget w) onPressed;
  const Card({
    super.key,
    required this.onPressed,
    this.margin = EdgeInsets.zero,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    final tap = () => onPressed(
      ds.Container(
        padding: defaults.padding,
        margin: EdgeInsets.zero,
        Settings(),
      ),
    );
    return ds.Card(
      alignment: Alignment.center,
      margin: margin,
      onTap: tap,
      help: ds.Hint(label: const Text("Google"), description: const Text("connect Google services")),
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing / 2,
        children: [
          Text("Google", style: theme.textTheme.titleMedium),
          Text(
            "Connect Google services",
            style: theme.textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}
