import 'package:flutter/material.dart';

class Field extends StatelessWidget {
  final Widget? label;
  final Widget input;

  Field({super.key, required this.input, this.label});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(child: label ?? Container()),
        Expanded(child: input, flex: 9),
      ],
    );
  }
}
