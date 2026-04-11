import 'package:flutter/material.dart';

class IconCheckmark extends StatelessWidget {
  final bool checked;
  final double? size;
  const IconCheckmark(this.checked, {super.key, this.size});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Opacity(
      opacity: checked ? 1.0 : 1.0,
      child: Icon(
        size: size,
        checked ? Icons.check : Icons.clear,
        color: checked ? Colors.lightGreenAccent : theme.disabledColor,
      ),
    );
  }
}
