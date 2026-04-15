import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './display.dart';

class Card extends StatelessWidget {
  final EdgeInsets? margin;
  final Function(Widget w) onPressed;
  const Card({super.key, required this.onPressed, this.margin});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    final tap = () => onPressed(
      Display(margin: EdgeInsets.zero, padding: EdgeInsets.zero),
    );
    return ds.Card(
      alignment: Alignment.center,
      margin: margin ?? defaults.margin,
      onTap: tap,
      help: ds.Hint(const Text("administer user accounts and roles")),
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing / 2,
        children: [
          Text("User Management", style: theme.textTheme.titleMedium),
          Text(
            "Manage permissions and access controls",
            style: theme.textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}
