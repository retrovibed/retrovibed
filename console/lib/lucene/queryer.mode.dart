import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'parser.results.dart';

class QueryerMode extends StatelessWidget {
  final ParserResult mode;

  const QueryerMode({
    super.key,
    required this.mode,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final chipTheme = theme.chipTheme;
    final bgColor = chipTheme.selectedColor ?? theme.colorScheme.secondaryContainer;
    final defaults = ds.Defaults.of(context);

    return Container(
      padding: EdgeInsets.symmetric(horizontal: defaults.spacing, vertical: 4),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(8),
      ),
      child: mode,
    );
  }
}
