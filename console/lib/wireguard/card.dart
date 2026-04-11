import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './list.dart';

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
      help: ds.Hint(label: const Text("WireGuard"), description: const Text("set up VPN connections")),
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing / 2,
        children: [
          Text("VPN - WireGuard", style: theme.textTheme.titleMedium),
          Text(
            "Access your library from anywhere",
            style: theme.textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}
