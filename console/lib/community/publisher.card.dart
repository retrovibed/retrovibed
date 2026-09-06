import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './publisher.list.dart';

class Card extends StatelessWidget {
  final EdgeInsets? margin;
  final Function(Widget w) onPressed;
  const Card({super.key, required this.onPressed, this.margin});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    final tap = () => onPressed(
      ds.Container(
        alignment: Alignment.topCenter,
        padding: defaults.padding,
        margin: EdgeInsets.zero,
        ListDisplay(),
      ),
    );
    return ds.Card(
      alignment: Alignment.center,
      margin: margin ?? defaults.margin,
      onTap: tap,
      help: ds.Hint(const Text("install and configure publishing plugins")),
      Column(
        spacing: defaults.spacing / 2,
        children: [
          Text("Publish Plugins", textAlign: TextAlign.center, style: theme.textTheme.titleMedium),
          Text(
            "share community content to social platforms",
            textAlign: TextAlign.center,
            style: theme.textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}
