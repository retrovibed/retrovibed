import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './list.searchable.dart';

class Card extends StatelessWidget {
  final EdgeInsets? margin;
  final Function(Widget w) onPressed;
  const Card({super.key, required this.onPressed, this.margin});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    final tap =
        () => onPressed(
          ds.Container(
            alignment: Alignment.topCenter,
            margin: EdgeInsets.zero,
            ListSearchable(),
          ),
        );
    return ds.Card(
      alignment: Alignment.center,
      margin: margin ?? defaults.margin,
      onTap: tap,
      help: ds.Hint(const Text("configure RSS feed subscriptions")),
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing / 2,
        children: [
          Text("RSS Feeds", style: theme.textTheme.titleMedium),
          Text(
            "Manage your RSS feed subscriptions",
            style: theme.textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}
