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

          return ds.Help(
            ds.Container(
              padding: EdgeInsets.symmetric(
                horizontal: defaults.spacing,
                vertical: 8.0,
              ),
              decoration: BoxDecoration(
                color: bgColor,
                borderRadius: BorderRadius.circular(8),
              ),
              Center(
                widthFactor: 1.0, // Prevents width from expanding infinitely horizontally
                heightFactor: 1.0, // Dictates that the Center takes up exactly the child's height
                child: mode,
              ),
            ),
            mode.help,
          );
        },
      ),
    );
  }
}
