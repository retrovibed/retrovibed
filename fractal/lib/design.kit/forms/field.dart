import 'package:flutter/material.dart';

class Field extends StatelessWidget {
  final Widget? label;
  final Widget input;
  final BoxConstraints constraints;

  Field({
    super.key,
    required this.input,
    this.label,
    this.constraints = const BoxConstraints(maxHeight: 48.0, minHeight: 48.0),
  });

  @override
  Widget build(BuildContext context) {
    return ConstrainedBox(
      constraints: constraints,
      child: Row(
        children: [
          Expanded(child: label ?? Container()),
          Expanded(child: input, flex: 9),
        ],
      ),
    );
  }
}
