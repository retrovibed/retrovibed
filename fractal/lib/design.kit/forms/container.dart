import 'package:fractal/design.kit/theme.defaults.dart' as theming;
import 'package:flutter/material.dart' as m;

class Container extends m.StatelessWidget {
  final m.Widget child;

  Container({super.key, required this.child});

  @override
  m.Widget build(m.BuildContext context) {
    return m.Container(
      padding: m.Theme.of(context).extension<theming.Defaults>()!.padding,
      child: child,
    );
  }
}
