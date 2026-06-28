import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'parser.results.dart';

class QueryerMode extends StatelessWidget {
  final ParserResult mode;
  final FocusNode focus;

  const QueryerMode({
    super.key,
    required this.mode,
    required this.focus,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final chipTheme = theme.chipTheme;
    final defaults = ds.Defaults.of(context);

    return Focus(
      focusNode: focus,
      child: Builder(
        builder: (context) {
          final focused = Focus.of(context).hasFocus;
          final bgColor = focused
              ? (chipTheme.selectedColor ?? theme.colorScheme.secondaryContainer)
              : (chipTheme.backgroundColor ?? theme.colorScheme.surfaceContainerLow);

          return Container(
            padding: EdgeInsets.symmetric(horizontal: defaults.spacing, vertical: 4),
            decoration: BoxDecoration(
              color: bgColor,
              borderRadius: BorderRadius.circular(8),
            ),
            child: mode,
          );
        },
      ),
    );
  }
}
