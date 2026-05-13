import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/theme.defaults.dart';

class HelpLabelled extends StatelessWidget {
  final Widget label;
  final Widget description;
  const HelpLabelled({super.key, required this.label, required this.description});

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    final padding = EdgeInsets.symmetric(
      horizontal: defaults.spacing,
      vertical: defaults.spacing / 2,
    );

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      spacing: defaults.spacing,
      children: [
        SizedBox(
          width: 120,
          child: Container(
            padding: padding,
            decoration: BoxDecoration(
              color: defaults.highlight,
              borderRadius: defaults.borderRadius,
            ),
            child: label,
          ),
        ),
        Expanded(
          child: Container(
            padding: padding,
            child: description,
          ),
        ),
      ],
    );
  }
}
