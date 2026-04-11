import 'package:flutter/material.dart';

class Field extends StatelessWidget {
  final Widget label;
  final Widget input;
  final EdgeInsets margin;
  final EdgeInsets padding;

  const Field({
    super.key,
    required this.input,
    this.label = const SizedBox(),
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
            Container(width: double.infinity, child: input),
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
