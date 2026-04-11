import 'package:flutter/material.dart' as m;
import 'package:retrovibed/design.kit/theme.defaults.dart' as theming;

class Container extends m.StatelessWidget {
  final m.Widget child;
  final m.EdgeInsets? padding;
  final m.EdgeInsets? margin;
  final m.BoxDecoration? decoration;
  final m.BoxConstraints? constraints;
  final m.Alignment? alignment;
  final m.Color? background;
  final m.Clip clipBehavior;

  Container(
    this.child, {
    super.key,
    this.alignment,
    this.constraints,
    this.decoration,
    this.padding,
    this.margin,
    this.background,
    this.clipBehavior = m.Clip.none,
  });

  @override
  m.Widget build(m.BuildContext context) {
    final theme = m.Theme.of(context);
    final defaults = theming.Defaults.of(context);
    return m.Container(
      constraints: constraints,
      alignment: alignment,
      clipBehavior: clipBehavior,
      margin: margin ?? m.EdgeInsets.zero,
      padding: padding ?? m.EdgeInsets.zero,
      decoration:
          decoration ??
          m.BoxDecoration(
            color: background ?? theme.colorScheme.surface,
            borderRadius: defaults.borderRadius,
          ),
      child: m.SelectionArea(child: child),
    );
  }
}
