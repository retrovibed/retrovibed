import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/theme.defaults.dart';

class Hint extends StatelessWidget {
  final Widget child;
  const Hint(this.child, {super.key});

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    final padding = EdgeInsets.symmetric(
      horizontal: defaults.spacing,
      vertical: defaults.spacing / 2,
    );

    return Container(
      padding: padding,
      child: child,
    );
  }
}
