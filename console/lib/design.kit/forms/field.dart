import 'package:flutter/material.dart';

class Field extends StatelessWidget {
  final Widget label;
  final Widget input;
  final List<Widget> trailing;
  final EdgeInsets margin;
  final EdgeInsets padding;

  const Field({
    super.key,
    required this.input,
    this.label = const SizedBox(),
    this.trailing = const [],
    this.margin = EdgeInsets.zero,
    this.padding = EdgeInsets.zero,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return SelectionArea(
      child: Container(
        margin: margin,
        padding: padding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Row(children: [Expanded(child: input), ...trailing]),
            DefaultTextStyle(
              style: theme.textTheme.bodySmall!.copyWith(
                color: theme.hintColor,
              ),
              child: label,
            ),
          ],
        ),
      ),
    );
  }
}
