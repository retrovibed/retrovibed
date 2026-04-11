import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/theme.defaults.dart';

class Hint extends StatelessWidget {
  final Widget label;
  final Widget description;
  const Hint({super.key, required this.label, required this.description});

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(
          width: 120,
          child: Container(
            padding: EdgeInsets.symmetric(
              horizontal: defaults.spacing,
              vertical: defaults.spacing / 2,
            ),
            decoration: BoxDecoration(
              color: defaults.highlight,
              borderRadius: defaults.borderRadius,
            ),
            child: label,
          ),
        ),
        SizedBox(width: defaults.spacing),
        Expanded(child: description),
      ],
    );
  }
}
