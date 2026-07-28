import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/feedback/settings.dart';

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
      help: ds.Hint(const Text("report an issue or start a discussion on GitHub")),
      Column(
        spacing: defaults.spacing / 2,
        children: [
          Text("GitHub", textAlign: TextAlign.center, style: theme.textTheme.titleMedium),
          Text(
            "Report an issue or start a discussion",
            textAlign: TextAlign.center,
            style: theme.textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}
